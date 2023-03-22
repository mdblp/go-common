// Package portal is a client module to support server-side use of the Diabeloop
// service called user-api.
package portal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
)

// API is the interface to portal-api.
type API interface {
	GetPatientConfig(token string) (*PatientConfig, error)
}

// Client used to store infos for this API
type Client struct {
	host       string
	httpClient *http.Client
}

const (
	routeV2GetPatientConfig = "/organization/v2/patient/params"
)

// NewClient create a new portal-api client
func NewClient(httpClient *http.Client, host string) (*Client, error) {
	_, err := url.Parse(host)
	if err != nil {
		return nil, errors.New("Invalid host url")
	}

	client := httpClient
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		host:       host,
		httpClient: client,
	}, nil
}

// NewClientFromEnv create a new portal-api client using PORTAL_HOST environnement variable
func NewClientFromEnv(httpClient *http.Client) (*Client, error) {
	host, haveHost := os.LookupEnv("PORTAL_HOST")
	if !haveHost || len(host) == 0 {
		return nil, errors.New("Missing PORTAL_HOST environnement variable")
	}
	return NewClient(httpClient, host)
}

// GetPatientConfig Return the patient configuration
//
// The token parameter is used to identify the patient.
func (client *Client) GetPatientConfig(token string) (*PatientConfig, error) {
	host, err := url.Parse(client.host)
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
