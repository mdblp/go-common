// Package portal is a client module to support server-side use of the Diabeloop
// service called user-api.
package portal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/tidepool-org/go-common/clients/disc"
)

// Client The interface to portal-api.
type Client struct {
	config     *ClientConfig
	hostGetter disc.HostGetter
	httpClient *http.Client
}

// ClientConfig Used to configure this client
type ClientConfig struct {
}

// ClientBuilder ...
type ClientBuilder struct {
	config     *ClientConfig
	hostGetter disc.HostGetter
	httpClient *http.Client
}

const (
	routeV2GetPatientConfig = "/organization/v2/patient/params"
)

// NewPortalClientBuilder create a new ClientBuilder
func NewPortalClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		config: &ClientConfig{},
	}
}

// WithHostGetter Set the host getter
func (b *ClientBuilder) WithHostGetter(val disc.HostGetter) *ClientBuilder {
	b.hostGetter = val
	return b
}

// WithHTTPClient set the HTTP client
func (b *ClientBuilder) WithHTTPClient(val *http.Client) *ClientBuilder {
	b.httpClient = val
	return b
}

// Build the portal client
func (b *ClientBuilder) Build() *Client {
	if b.hostGetter == nil {
		panic("PortalClient requires a hostGetter to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &Client{
		config:     b.config,
		hostGetter: b.hostGetter,
		httpClient: b.httpClient,
	}
}

func (client *Client) getHost() (*url.URL, error) {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		*cpy = hostArr[0]
		return cpy, nil
	}
	return nil, errors.New("no known portal-api hosts")
}

// GetPatientConfig Return the patient configuration
//
// It use the patient token, to identify the patient to use
func (client *Client) GetPatientConfig(token string) (*PatientConfig, error) {
	host, err := client.getHost()
	if host == nil || err != nil {
		return nil, err
	}

	host.Path = path.Join(host.Path, routeV2GetPatientConfig)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unknown response code[%d] from service[%s]", res.StatusCode, req.URL)
	}

	var pc PatientConfig
	if err = json.NewDecoder(res.Body).Decode(&pc); err != nil {
		return nil, fmt.Errorf("Error parsing JSON results: %v", err)
	}

	return &pc, nil
}
