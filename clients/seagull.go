package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/tidepool-org/go-common/clients/status"
	"github.com/tidepool-org/go-common/errors"
)

type (
	Seagull interface {
		// Retrieves arbitrary collection information from metadata
		//
		// userID -- the Tidepool-assigned userId
		// hashName -- the name of what we are trying to get
		// token -- a server token or the user token
		GetPrivatePair(userID, hashName, token string) *PrivatePair
		// Retrieves arbitrary collection information from metadata
		//
		// userID -- the Tidepool-assigned userId
		// collectionName -- the collection being retrieved
		// token -- a server token or the user token
		// v - the interface to return the value in
		GetCollection(userID, collectionName, token string, v interface{}) error
	}

	seagullClient struct {
		httpClient *http.Client // store a reference to the http client so we can reuse it
		host       string
	}

	seagullClientBuilder struct {
		httpClient *http.Client
		host       string
	}

	PrivatePair struct {
		ID    string
		Value string
	}
)

func NewSeagullClientBuilder() *seagullClientBuilder {
	return &seagullClientBuilder{}
}

func (b *seagullClientBuilder) WithHttpClient(httpClient *http.Client) *seagullClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *seagullClientBuilder) WithHost(host string) *seagullClientBuilder {
	b.host = host
	return b
}

func (b *seagullClientBuilder) Build() *seagullClient {
	if b.httpClient == nil {
		panic("seagullClient requires an httpClient to be set")
	}
	if b.host == "" {
		panic("seagullClient requires a host to be set")
	}
	return &seagullClient{
		httpClient: b.httpClient,
		host:       b.host,
	}
}
func NewSeagullClientFromEnv(httpClient *http.Client) *seagullClient {
	builder := NewSeagullClientBuilder()
	host, _ := os.LookupEnv("SEAGULL_HOST")
	return builder.WithHost(host).
		WithHttpClient(httpClient).
		Build()
}

func (client *seagullClient) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	host, err := client.getHost()
	if err != nil {
		return nil
	}
	host.Path = path.Join(host.Path, userID, "private", hashName)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	log.Println(req)
	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Problem when looking up private pair for userID[%s]. %s", userID, err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Unknown response code[%v] from service[%v]", res.StatusCode, req.URL)
		return nil
	}

	var retVal PrivatePair
	if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
		log.Println("Error parsing JSON results", err)
		return nil
	}
	return &retVal
}

func (client *seagullClient) GetCollection(userID, collectionName, token string, v interface{}) error {
	host, err := client.getHost()
	if err != nil {
		return err
	}
	host.Path = path.Join(host.Path, userID, collectionName)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	log.Println(req)
	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Problem when looking up collection for userID[%s]. %s", userID, err)
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
			log.Println("Error parsing JSON results", err)
			return err
		}
		return nil
	case http.StatusNotFound:
		log.Printf("No [%s] collection found for [%s]", collectionName, userID)
		return nil
	default:
		return &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}

}

func (client *seagullClient) getHost() (*url.URL, error) {
	if client.host == "" {
		return nil, errors.New("No client host defined")
	}
	theURL, err := url.Parse(client.host)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse urlString[%s]", client.host)
	}
	return theURL, nil
}
