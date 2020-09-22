package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Store and Mongo clients used to manage database connections 
type (
	// StorageIterator - Interface for the query iterator
	StorageIterator interface {
		Next(ctx context.Context) bool
		Close(ctx context.Context) error
		Decode(val interface{}) error
	}
	// Storage - Interface for our storage layer
	Storage interface {
		Close() error
		Collection(collectionName string, databaseName ...string) *mongo.Collection
		Ping() error
		PingOK() bool
		Start()
		WaitUntilStarted()
	}
	// StoreClient - Mongo Storage Client
	StoreClient struct {
		client          *mongo.Client
		Context         context.Context
		config          *Config
		logger          *log.Logger
		closingChannel  chan bool
		initializeGroup sync.WaitGroup
		pingOK          bool
		clientMux       sync.Mutex
	}
)

func NewStoreClient(config *Config, logger *log.Logger) (*StoreClient, error) {
	if config.Timeout <= 0 {
		return nil, errors.New("timeout is invalid")
	}
	mongoClient, err := newMongoClient(config)
	if err != nil {
		return nil, err
	}

	store := &StoreClient{
		client:  mongoClient,
		Context: context.Background(),
		config:  config,
		logger:  logger,
	}

	return store, nil
}
func newMongoClient(config *Config) (*mongo.Client, error) {
	connectionString, err := config.toConnectionString()
	if err != nil {
		return nil, err
	}
	clientOptions := options.Client().ApplyURI(connectionString)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return mongoClient, nil
}

// Mutex getters / setters
func (s *StoreClient) getClient() *mongo.Client {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	return s.client
}
func (s *StoreClient) setClient(cli *mongo.Client) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	s.client = cli
}
func (s *StoreClient) PingOK() bool {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	return s.pingOK
}
func (s *StoreClient) setPingOK(ping bool) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	s.pingOK = ping
}

func (s *StoreClient) Start() {
	if s.closingChannel == nil {
		s.initializeGroup.Add(1)
		go s.connectionRoutine()
	}
}

func (s *StoreClient) Close() error {
	if s.closingChannel != nil {
		s.closingChannel <- true
	}
	s.initializeGroup.Wait()
	return s.getClient().Disconnect(s.Context)
}

func (s *StoreClient) Ping() error {
	if s.getClient() == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	err := s.getClient().Ping(ctx, readpref.PrimaryPreferred())
	s.setPingOK(err == nil)
	return err
}

func (s *StoreClient) connectionRoutine() {
	err := s.Ping()
	var attempts int64
	if err != nil {
		s.logger.Printf("Unable to open inital store session : %v", err)
		s.closingChannel = make(chan bool, 1)
		for {
			timer := time.After(s.config.WaitConnectionInterval)
			select {
			case <-s.closingChannel:
				close(s.closingChannel)
				s.closingChannel = nil
				s.initializeGroup.Done()
				return
			case <-timer:
				err := s.Ping()
				if err == nil {
					s.logger.Debug("Store session opened succesfully")
					s.logger.Printf("Store pinged succesfully after %v attempts, creating indexes", attempts)
					s.createIndexesFromConfig()
					s.closingChannel <- true
				} else {
					if s.config.MaxConnectionAttempts > 0 && s.config.MaxConnectionAttempts < attempts {
						s.logger.Printf("Unable to open store session, maximum connection attempts reached (%v) : %v", s.config.MaxConnectionAttempts, err)
						s.closingChannel <- true
						panic(err)
					} else {
						s.logger.Printf("Unable to open store session : %v", err)
						attempts++
					}
				}
			}
		}
	} else {
		s.logger.Printf("Store client up and running, creating indexes")
		s.createIndexesFromConfig()
		if s.closingChannel != nil {
			close(s.closingChannel)
			s.closingChannel = nil
		}
		s.initializeGroup.Done()
		return
	}
}
func (s *StoreClient) WaitUntilStarted() {
	s.initializeGroup.Wait()
}

func (s *StoreClient) Collection(collectionName string, databaseName ...string) *mongo.Collection {
	dbName := s.config.Database
	if len(databaseName) > 0 {
		dbName = databaseName[0]
	}
	return s.getClient().Database(dbName).Collection(collectionName)
}

func (s *StoreClient) createIndexesFromConfig() {
	if s.config.Indexes != nil {
		for collection, idxs := range s.config.Indexes {
			if _, err := s.Collection(collection).Indexes().CreateMany(context.Background(), idxs); err != nil {
				s.logger.Printf(fmt.Sprintf("Unable to create indexes: %s", err))
			}
		}
	}
}
