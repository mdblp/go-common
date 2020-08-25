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

	"github.com/mdblp/go-common/clients/disc"
)

// Client is the interface to portal-api.
type Client interface {
	GetPatientConfig(token string) (*PatientConfig, error)
}

// ClientStruct used to store infos for this API
type ClientStruct struct {
	hostGetter disc.HostGetter
	httpClient *http.Client
}

// ClientBuilder same as Client but with a different API
type ClientBuilder struct {
	hostGetter disc.HostGetter
	httpClient *http.Client
}

const (
	routeV2GetPatientConfig = "/organization/v2/patient/params"
)

// NewPortalClientBuilder create a new ClientBuilder
func NewPortalClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

// WithHostGetter set the host getter
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
func (b *ClientBuilder) Build() *ClientStruct {
	if b.hostGetter == nil {
		panic("PortalClient requires a hostGetter to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &ClientStruct{
		hostGetter: b.hostGetter,
		httpClient: b.httpClient,
	}
}

func (client *ClientStruct) getHost() (*url.URL, error) {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		// TODO allow to use more than one hostname? :
		*cpy = hostArr[0]
		return cpy, nil
	}
	return nil, errors.New("no known portal-api hosts")
}

// GetPatientConfig Return the patient configuration
//
// The token parameter is used to identify the patient.
func (client *ClientStruct) GetPatientConfig(token string) (*PatientConfig, error) {
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
