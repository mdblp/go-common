package mongo

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestNoDatabase(t *testing.T) {
	x := Config{}
	_, err := x.toConnectionString()

	if err == nil {
		t.Error("database is required")
	}
}

func TestDatabase(t *testing.T) {
	x := Config{Database: "admin"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Errorf("should not error %v", err)
	}
	if s != "mongodb://localhost/admin?ssl=false" {
		t.Errorf("found %v", s)
	}
}

func TestScheme(t *testing.T) {
	x := Config{Database: "admin", Scheme: "mongodb+srv"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb+srv://localhost/admin?ssl=false" {
		t.Errorf("found %v", s)
	}
}

func TestUser(t *testing.T) {
	x := Config{Database: "admin", Username: "derrick"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick@localhost/admin?ssl=false" {
		t.Errorf("found %v", s)
	}
}

func TestPassword(t *testing.T) {
	x := Config{Database: "admin", Username: "derrick", Password: "password"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick:password@localhost/admin?ssl=false" {
		t.Errorf("found %v", s)
	}
}

func TestSsl(t *testing.T) {
	x := Config{Database: "admin", Username: "derrick", Password: "password", TLS: true}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick:password@localhost/admin?ssl=true" {
		t.Errorf("found %v", s)
	}
}

func TestHosts(t *testing.T) {
	x := Config{Database: "admin", Username: "derrick", Password: "password", TLS: true, addresses: []string{"mongodb1", "mongodb2"}}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick:password@mongodb1,mongodb2/admin?ssl=true" {
		t.Errorf("found %v", s)
	}
}

func TestOptParams(t *testing.T) {
	x := Config{Database: "admin", Username: "derrick", Password: "password", TLS: true, addresses: []string{"mongodb1", "mongodb2"}, OptParams: "x=y"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick:password@mongodb1,mongodb2/admin?ssl=true&x=y" {
		t.Errorf("found %v", s)
	}
}

func TestConnection(t *testing.T) {
	config := Config{
		// ConnectionString:       "mongodb://localhost/",
		Database:               "admin",
		Timeout:                2 * time.Second,
		WaitConnectionInterval: 5 * time.Second,
		MaxConnectionAttempts:  0,
	}
	if _, exist := os.LookupEnv("TIDEPOOL_STORE_ADDRESSES"); exist {
		// if mongo connexion information is provided via env var
		config.FromEnv()
	}
	logger := log.New(os.Stdout, "mongo-test ", log.LstdFlags|log.LUTC|log.Lshortfile)

	store, err := NewStoreClient(&config, logger)
	if err != nil {
		t.Errorf("connection failed: %v", err)
	}
	defer store.Close()
	store.Start()
	store.WaitUntilStarted()
	// Expect the connection to be established
	if !store.PingOK() {
		t.Errorf("connection routine failed")
	}

	err = store.Ping()
	if err != nil || !store.PingOK() {
		t.Errorf("ping failed: %v", err)
	}

	err = store.Close()
	if err != nil {
		t.Errorf("disconnect failed: %v", err)
	}
}
func TestReConnectionOnStartup(t *testing.T) {
	config := Config{
		addresses:              []string{"hheheihezhfoehz"},
		Database:               "admin",
		Timeout:                1 * time.Second,
		WaitConnectionInterval: 1 * time.Second,
		MaxConnectionAttempts:  10,
	}
	address := "localhost:27017"
	if env_adress, exist := os.LookupEnv("TIDEPOOL_STORE_ADDRESSES"); exist {
		// if mongo connexion information is provided via env var
		address = env_adress
	}
	logger := log.New(os.Stdout, "mongo-test ", log.LstdFlags|log.LUTC|log.Lshortfile)

	store, err := NewStoreClient(&config, logger)
	if err != nil {
		t.Errorf("connection failed: %v", err)
	}
	defer store.Close()
	store.Start()
	time.Sleep(2 * config.WaitConnectionInterval)
	// Expect the connection to not be established yet
	if store.PingOK() {
		t.Errorf("connection should have fail")
	}
	// Expect the connection to be established once server is up
	store.config.addresses = []string{address}
	client, err := newMongoClient(store.config)
	if err != nil {
		t.Errorf("Error creating mongo.client : %v", err)
	}
	store.setClient(client)
	store.WaitUntilStarted()
	if !store.PingOK() {
		t.Errorf("connection routine failed")
	}
	err = store.Ping()
	if err != nil || !store.PingOK() {
		t.Errorf("ping failed: %v", err)
	}

	err = store.Close()
	if err != nil {
		t.Errorf("disconnect failed: %v", err)
	}
}
