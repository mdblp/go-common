package mongo

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mdblp/go-common/errors"
	"github.com/mdblp/go-common/jepson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config of the mongo database
type Config struct {
	ConnectionString string           `json:"connectionString"`
	Timeout          *jepson.Duration `json:"timeout"`
	Scheme           string           `json:"scheme"`
	User             string           `json:"user"`
	Password         string           `json:"password"`
	Database         string           `json:"database"`
	Ssl              bool             `json:"ssl"`
	Hosts            string           `json:"hosts"`
	OptParams        string           `json:"optParams"`
}

// Store contains the connection information
type Store struct {
	logger    *log.Logger
	client    *mongo.Client
	pingOK    bool
	clientMux sync.Mutex
	pingOkMux sync.Mutex
}

const defaultTimeout = time.Second * 2

// FromEnv read the mongo config from the environment variables
func (config *Config) FromEnv() {
	config.Scheme, _ = os.LookupEnv("TIDEPOOL_STORE_SCHEME")
	config.Hosts, _ = os.LookupEnv("TIDEPOOL_STORE_ADDRESSES")
	config.User, _ = os.LookupEnv("TIDEPOOL_STORE_USERNAME")
	config.Password, _ = os.LookupEnv("TIDEPOOL_STORE_PASSWORD")
	config.Database, _ = os.LookupEnv("TIDEPOOL_STORE_DATABASE")
	config.OptParams, _ = os.LookupEnv("TIDEPOOL_STORE_OPT_PARAMS")
	ssl, found := os.LookupEnv("TIDEPOOL_STORE_TLS")
	config.Ssl = found && ssl == "true"
}

func (config *Config) toConnectionString() (string, error) {
	if config.ConnectionString != "" {
		return config.ConnectionString, nil
	}
	if config.Database == "" {
		return "", errors.New("Must specify a database in Mongo config")
	}

	var cs string
	if config.Scheme != "" {
		cs = config.Scheme + "://"
	} else {
		cs = "mongodb://"
	}

	if config.User != "" {
		cs += config.User
		if config.Password != "" {
			cs += ":"
			cs += config.Password
		}
		cs += "@"
	}

	if config.Hosts != "" {
		cs += config.Hosts
		cs += "/"
	} else {
		cs += "localhost/"
	}

	if config.Database != "" {
		cs += config.Database
	}

	if config.Ssl {
		cs += "?ssl=true"
	} else {
		cs += "?ssl=false"
	}

	if config.OptParams != "" {
		cs += "&"
		cs += config.OptParams
	}
	return cs, nil
}

// Connect perform a mongo connexion.
// The connexion may not be directly available, but we will retry
func Connect(config *Config, logger *log.Logger) (*Store, error) {
	connectionString, err := config.toConnectionString()
	if err != nil {
		return nil, err
	}

	store := &Store{
		logger: logger,
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	// Do the connection async, since the services must be able to start
	// without the database
	go store.connectionRoutine(client)

	return store, nil
}

func (s *Store) getClient() *mongo.Client {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	return s.client
}
func (s *Store) setClient(cli *mongo.Client) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	s.client = cli
}

func (s *Store) PingOK() bool {
	s.pingOkMux.Lock()
	defer s.pingOkMux.Unlock()
	return s.pingOK
}
func (s *Store) setPingOK(ping bool) {
	s.pingOkMux.Lock()
	defer s.pingOkMux.Unlock()
	s.pingOK = ping
}

func (s *Store) connectionRoutine(client *mongo.Client) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.logger.Printf("Connecting to mongo...")
	err = client.Connect(ctx)
	if err != nil {
		s.logger.Printf("Connection to mongo failed: %v", err)
		time.Sleep(defaultTimeout)
		go s.connectionRoutine(client)
	} else {
		s.logger.Printf("Connected to mongo")
		s.setPingOK(true)
		s.setClient(client)
	}
}

// GetCollection return a collection on a database
func (s *Store) GetCollection(database, collection string) *mongo.Collection {
	return s.getClient().Database(database).Collection(collection)
}

// Ping the database
func (s *Store) Ping() error {
	if s.getClient() == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	return s.getClient().Ping(ctx, readpref.PrimaryPreferred())
}

// ContinuousPing the database to monitor the database connection
//
// Update the Store.PingOK status
func (s *Store) ContinuousPing(timeout time.Duration) {
	if s.getClient() == nil {
		s.logger.Printf("Stopping continuous ping")
		return
	}

	time.Sleep(timeout)

	err := s.Ping()
	if err != nil && s.PingOK() {
		s.logger.Printf("Ping error: %v", err)
	} else if err == nil && !s.PingOK() {
		s.logger.Printf("Ping: mongo restored")
	}
	s.setPingOK(err == nil)

	go s.ContinuousPing(timeout)
}

// Disconnect from the database
func (s *Store) Disconnect() error {
	if s.getClient() == nil {
		return nil
	}

	s.logger.Printf("Disconnecting mongo database...")
	client := s.getClient()
	s.setClient(nil)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := client.Disconnect(ctx)

	return err
}
