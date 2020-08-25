package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/mdblp/go-common/clients/status"
	"github.com/mdblp/go-common/errors"
)

type (
	//Inteface so that we can mock gatekeeperClient for tests
	Gatekeeper interface {
		//userID  -- the Tidepool-assigned userID
		//groupID  -- the Tidepool-assigned groupID
		//
		// returns the Permissions
		UserInGroup(userID, groupID string) (Permissions, error)

		//groupID  -- the Tidepool-assigned groupID
		//
		// returns the map of user id to Permissions
		UsersInGroup(groupID string) (UsersPermissions, error)

		// returns the map of user id to Permissions
		GroupsForUser(userID string) (UsersPermissions, error)

		//userID  -- the Tidepool-assigned userID
		//groupID  -- the Tidepool-assigned groupID
		//permissions -- the permisson we want to give the user for the group
		SetPermissions(userID, groupID string, permissions Permissions) (Permissions, error)
	}

	gatekeeperClient struct {
		httpClient    *http.Client // store a reference to the http client so we can reuse it
		host          string
		tokenProvider TokenProvider // An object that provides tokens for communicating with gatekeeper
	}

	gatekeeperClientBuilder struct {
		httpClient    *http.Client // store a reference to the http client so we can reuse it
		host          string
		tokenProvider TokenProvider // An object that provides tokens for communicating with gatekeeper
	}

	Permission       map[string]interface{}
	Permissions      map[string]Permission
	UsersPermissions map[string]Permissions
)

// defaultHost for Gatekeeper client
const defaultHost = "http://localhost:9123"

var (
	Allowed Permission = Permission{}
)

func NewGatekeeperClientBuilder() *gatekeeperClientBuilder {
	return &gatekeeperClientBuilder{}
}

func (b *gatekeeperClientBuilder) WithHttpClient(httpClient *http.Client) *gatekeeperClientBuilder {
	b.httpClient = httpClient
	return b
}

// WithHost set the gatekeeper URL
func (b *gatekeeperClientBuilder) WithHost(host string) *gatekeeperClientBuilder {
	b.host = host
	return b
}

func (b *gatekeeperClientBuilder) WithTokenProvider(tokenProvider TokenProvider) *gatekeeperClientBuilder {
	b.tokenProvider = tokenProvider
	return b
}

func (b *gatekeeperClientBuilder) Build() *gatekeeperClient {
	if b.host == "" {
		log.Printf("No Gatekeeper host set using default %s", defaultHost)
		b.host = defaultHost
	}
	if b.tokenProvider == nil {
		panic("gatekeeperClient requires a tokenProvider to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &gatekeeperClient{
		httpClient:    b.httpClient,
		host:          b.host,
		tokenProvider: b.tokenProvider,
	}
}

// UserInGroup Check whether one subject is sharing data with one other user
// original route /access/{userid}/{granteeid}
// userID ID of the user to check for having permissions to view subject's data
// groupID ID of the user subject
func (client *gatekeeperClient) UserInGroup(userID, groupID string) (Permissions, error) {
	host, err := client.getHost()
	if host == nil {
		return nil, errors.New("No known gatekeeper hosts")
	}
	host.Path = path.Join(host.Path, "access", groupID, userID)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retVal := make(Permissions)
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			log.Println(err)
			return nil, &status.StatusError{status.NewStatus(500, "UserInGroup Unable to parse response.")}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

// UsersInGroup List of users one subject is sharing data with
// original route /access/{userid}
// groupID ID of the user subject
func (client *gatekeeperClient) UsersInGroup(groupID string) (UsersPermissions, error) {
	host, err := client.getHost()
	if err != nil {
		return nil, fmt.Errorf("Gatekeeper url is invalid: %v", err)
	}
	host.Path = path.Join(host.Path, "access", groupID)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retVal := make(UsersPermissions)
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			log.Println(err)
			return nil, &status.StatusError{status.NewStatus(500, "UserInGroup Unable to parse response.")}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

// GroupsForUser returns the list of users sharing data with one subject
// original route /access/groups/{userid}
// userID ID of the user subject
func (client *gatekeeperClient) GroupsForUser(userID string) (UsersPermissions, error) {
	host, err := client.getHost()
	if err != nil {
		return nil, fmt.Errorf("Gatekeeper url is invalid: %v", err)
	}
	host.Path = path.Join(host.Path, "access", "groups", userID)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retVal := make(UsersPermissions)
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			log.Println(err)
			return nil, &status.StatusError{status.NewStatus(500, "GroupsForUser Unable to parse response.")}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

func (client *gatekeeperClient) SetPermissions(userID, groupID string, permissions Permissions) (Permissions, error) {
	host, err := client.getHost()
	if err != nil {
		return nil, fmt.Errorf("Gatekeeper url is invalid: %v", err)
	}
	host.Path = path.Join(host.Path, "access", groupID, userID)

	if jsonPerms, err := json.Marshal(permissions); err != nil {
		log.Println(err)
		return nil, &status.StatusError{status.NewStatusf(http.StatusInternalServerError, "Error marshaling the permissons [%s]", err)}
	} else {
		req, _ := http.NewRequest("POST", host.String(), bytes.NewBuffer(jsonPerms))
		req.Header.Set("content-type", "application/json")
		req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

		res, err := client.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == 200 {
			retVal := make(Permissions)
			if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
				log.Printf("SetPermissions: Unable to parse response: [%s]", err.Error())
				return nil, &status.StatusError{status.NewStatus(500, "SetPermissions: Unable to parse response:")}
			}
			return retVal, nil
		} else if res.StatusCode == 404 {
			return nil, nil
		} else {
			return nil, &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
		}
	}
}

func (client *gatekeeperClient) getHost() (*url.URL, error) {
	return url.Parse(client.host)
}
