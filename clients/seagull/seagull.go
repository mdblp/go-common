package seagull

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/mdblp/go-common/clients/status"
	"github.com/mdblp/go-common/errors"
)

type (
	// API seagull client interface
	API interface {
		// Retrieves arbitrary collection information from metadata
		//
		// userID -- the Tidepool-assigned userId
		// collectionName -- the collection being retrieved
		// token -- a server token or the user token
		// v - the interface to return the value in
		GetCollection(userID, collectionName, token string, v interface{}) error
		// Set arbitrary collection information from metadata
		SetCollection(userID, collectionName, token string, payload interface{}) error
	}

	// Client for seagull
	Client struct {
		httpClient *http.Client // store a reference to the http client so we can reuse it
		host       string
	}
)

// NewClient create a new seagull client
func NewClient(httpClient *http.Client, host string) (*Client, error) {
	_, err := url.Parse(host)
	if err != nil {
		return nil, errors.New("Invalid seagull host url")
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

// NewClientFromEnv create a new seagull client using SEAGULL_HOST environnement variable
func NewClientFromEnv(httpClient *http.Client) (*Client, error) {
	host, haveHost := os.LookupEnv("SEAGULL_HOST")
	if !haveHost || len(host) == 0 {
		return nil, errors.New("Missing SEAGULL_HOST environnement variable")
	}
	return NewClient(httpClient, host)
}

// GetCollection return a seagull collection content
func (client *Client) GetCollection(userID, collectionName, token string, v interface{}) error {
	host, err := url.Parse(client.host)
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
		return &status.StatusError{Status: status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

func (client *Client) SetCollection(userID, collectionName, token string, payload interface{}) error {
	host, err := url.Parse(client.host)
	if err != nil {
		return err
	}
	host.Path = path.Join(host.Path, userID, collectionName)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", host.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("x-tidepool-session-token", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Problem when looking up collection for userID[%s]. %s", userID, err)
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		log.Printf("No [%s] collection found for [%s]", collectionName, userID)
		return nil
	default:
		return &status.StatusError{Status: status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}
