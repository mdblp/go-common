package opa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/tidepool-org/go-common/clients/status"
)

// Client is the interface to opa.
type Client interface {
	GetOpaAuth(req *http.Request) (*Authorization, error)
}

// ClientStruct used to store infos for this API
type ClientStruct struct {
	host              *string
	httpClient        *http.Client
	requestingService string
}

// ClientBuilder same as Client but with a different API
type ClientBuilder struct {
	host              *string
	httpClient        *http.Client
	requestingService string
}

// Authorization struct for authz
// "result": {
//     "authorized": false,
//     "data": {
//         "userIds": [
//             "00004"
//         ]
//     },
//     "route": "tidewhisperer-get"
// }
type Authorization struct {
	Result *opaResult `json:"result"`
}
type opaResult struct {
	Authorized bool                   `json:"authorized"`
	Data       map[string]interface{} `json:"data"`
	Route      string                 `json:"route"`
}

// HTTPInput struct sent to OPA
type HTTPInput struct {
	Input struct {
		Request struct {
			Headers  map[string]string `json:"headers"`
			Host     string            `json:"host"`
			Method   string            `json:"method"`
			Path     string            `json:"path"`
			Query    string            `json:"query"`
			Fragment string            `json:"fragment"`
			Protocol string            `json:"protocol"`
			Service  string            `json:"service"`
		} `json:"request"`
	} `json:"input"`
}

const (
	routeAuth = "/v1/data/backloops/access"
)

// NewClientFromEnv read the config from the environment variables
func NewClientFromEnv(httpClient *http.Client) *ClientStruct {
	builder := NewOpaClientBuilder()
	host, _ := os.LookupEnv("OPA_HOST")
	requestingService, _ := os.LookupEnv("SERVICE_NAME")
	return builder.WithHost(host).
		WithHTTPClient(httpClient).
		WithRequestingService(requestingService).
		Build()
}

// NewOpaClientBuilder create a new ClientBuilder
func NewOpaClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

// WithHost set the host getter
func (b *ClientBuilder) WithHost(host string) *ClientBuilder {
	b.host = &host
	return b
}

// WithHTTPClient set the HTTP client
func (b *ClientBuilder) WithHTTPClient(val *http.Client) *ClientBuilder {
	b.httpClient = val
	return b
}

// WithRequestingService set the HTTP client
func (b *ClientBuilder) WithRequestingService(val string) *ClientBuilder {
	b.requestingService = val
	return b
}

// Build the portal client
func (b *ClientBuilder) Build() *ClientStruct {
	if b.host == nil {
		panic("OpaClient requires a hostGetter to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &ClientStruct{
		host:              b.host,
		httpClient:        b.httpClient,
		requestingService: b.requestingService,
	}
}

func (client *ClientStruct) getHost() (*url.URL, error) {
	if client.host == nil {
		return nil, errors.New("No client host defined")
	}
	theURL, err := url.Parse(*client.host)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse urlString[%s]", *client.host)
	}
	return theURL, nil
}

func (client *ClientStruct) formatRequest(req *http.Request) (*HTTPInput, error) {
	var opaReq HTTPInput
	url := *req.URL
	headers := make(map[string]string)
	for k := range req.Header {
		headers[strings.ToLower(k)] = req.Header.Get(k)
	}
	opaReq.Input.Request.Headers = headers
	opaReq.Input.Request.Method = req.Method
	opaReq.Input.Request.Protocol = req.Proto
	opaReq.Input.Request.Host = req.Host
	opaReq.Input.Request.Path = url.Path
	opaReq.Input.Request.Query = url.RawQuery
	opaReq.Input.Request.Fragment = url.RawFragment
	opaReq.Input.Request.Service = client.requestingService
	return &opaReq, nil
}

// GetOpaAuth Return the patient configuration
//
// The token parameter is used to identify the patient.
func (client *ClientStruct) GetOpaAuth(req *http.Request) (*Authorization, error) {
	host, err := client.getHost()
	if host == nil || err != nil {
		return nil, err
	}

	host.Path = path.Join(host.Path, routeAuth)
	myRequest, _ := client.formatRequest(req)
	if jsonRequest, err := json.Marshal(*myRequest); err != nil {
		return nil, &status.StatusError{status.NewStatusf(http.StatusInternalServerError, "Error formatting request [%s]", err.Error())}
	} else {
		req, _ := http.NewRequest("POST", host.String(), bytes.NewBuffer(jsonRequest))

		res, err := client.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("Unknown response code[%d] from service[%s]", res.StatusCode, req.URL)
		}

		var auth Authorization
		if err = json.NewDecoder(res.Body).Decode(&auth); err != nil {
			return nil, fmt.Errorf("Error parsing JSON results: %v", err)
		}

		return &auth, nil
	}
}
