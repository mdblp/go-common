package opa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	URL "net/url"
	"os"
	"path"
	"strings"

	"github.com/mdblp/go-common/clients/status"
)

// API is the interface to opa.
type API interface {
	GetOpaAuth(req *http.Request) (*Authorization, error)
}

// Client used to store infos for this API
type Client struct {
	host              string
	requestingService string
	httpClient        *http.Client
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

// NewClient create a new OPA client with the specified host & service
func NewClient(httpClient *http.Client, host string, service string) (*Client, error) {
	if len(host) == 0 {
		return nil, errors.New("host is empty")
	}
	_, err := url.Parse(host)
	if err != nil {
		return nil, errors.New("Invalid host url")
	}
	if len(service) == 0 {
		return nil, errors.New("Empty service name")
	}

	client := httpClient
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		host:              host,
		requestingService: service,
		httpClient:        client,
	}, nil
}

// NewClientFromEnv create a new opa client using environnement variables
//
// OPA_HOST for the host
//
// SERVICE_NAME For the current (requests) service name
func NewClientFromEnv(httpClient *http.Client) (*Client, error) {
	host, haveHost := os.LookupEnv("OPA_HOST")
	if !haveHost {
		return nil, errors.New("Missing OPA_HOST environnement variable")
	}
	service, haveService := os.LookupEnv("SERVICE_NAME")
	if !haveService {
		return nil, errors.New("Missing SERVICE_NAME for OPA environnement variable")
	}

	return NewClient(httpClient, host, service)
}

func (client *Client) formatRequest(req *http.Request) (*HTTPInput, error) {
	var err error
	var opaReq HTTPInput
	var decodedString string

	url := *req.URL
	headers := make(map[string]string)
	for k := range req.Header {
		headers[strings.ToLower(k)] = req.Header.Get(k)
	}
	if decodedString, err = URL.QueryUnescape(url.RawQuery); err != nil {
		return nil, fmt.Errorf("Unable to parse query String [%s]", err)
	}

	url.RawQuery = decodedString
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
func (client *Client) GetOpaAuth(req *http.Request) (*Authorization, error) {
	var jsonRequest []byte
	var err error
	host, err := url.Parse(client.host)
	if host == nil || err != nil {
		return nil, err
	}

	host.Path = path.Join(host.Path, routeAuth)
	myRequest, _ := client.formatRequest(req)
	if jsonRequest, err = json.Marshal(*myRequest); err != nil {
		return nil, &status.StatusError{
			Status: status.NewStatusf(http.StatusInternalServerError, "Error formatting request [%s]", err.Error()),
		}
	}
	req, err = http.NewRequest("POST", host.String(), bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, err
	}

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
