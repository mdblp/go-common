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
		t.Error("should not error")
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
	x := Config{Database: "admin", User: "derrick"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick@localhost/admin?ssl=false" {
		t.Errorf("found %v", s)
	}
}

func TestPassword(t *testing.T) {
	x := Config{Database: "admin", User: "derrick", Password: "password"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick:password@localhost/admin?ssl=false" {
		t.Errorf("found %v", s)
	}
}

func TestSsl(t *testing.T) {
	x := Config{Database: "admin", User: "derrick", Password: "password", Ssl: true}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick:password@localhost/admin?ssl=true" {
		t.Errorf("found %v", s)
	}
}

func TestHosts(t *testing.T) {
	x := Config{Database: "admin", User: "derrick", Password: "password", Ssl: true, Hosts: "mongodb1,mongodb2"}
	s, err := x.toConnectionString()

	if err != nil {
		t.Error("should not error")
	}
	if s != "mongodb://derrick:password@mongodb1,mongodb2/admin?ssl=true" {
		t.Errorf("found %v", s)
	}
}

func TestOptParams(t *testing.T) {
	x := Config{Database: "admin", User: "derrick", Password: "password", Ssl: true, Hosts: "mongodb1,mongodb2", OptParams: "x=y"}
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
		ConnectionString: "mongodb://localhost/",
	}
	logger := log.New(os.Stdout, "mongo-test ", log.LstdFlags|log.LUTC|log.Lshortfile)

	store, err := Connect(&config, logger)
	if err != nil {
		t.Errorf("connection failed: %v", err)
	}

	defer store.Disconnect()

	// Expect the connection to be established
	time.Sleep(100 * time.Millisecond)
	if !store.PingOK {
		t.Errorf("connection routine failed")
	}

	err = store.Ping()
	if err != nil || !store.PingOK {
		t.Errorf("ping failed: %v", err)
	}

	err = store.Disconnect()
	if err != nil {
		t.Errorf("disconnect failed: %v", err)
	}
}
