package mongo

import (
	"log"
	"os"
	"strconv"
	"strings"
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

func overrideEnvs(overrides map[string]string) map[string]string {
	current := make(map[string]string, len(overrides))
	for env := range overrides {
		orig := os.Getenv(env)
		os.Setenv(env, overrides[env])
		current[env] = orig
	}
	return current
}
func restoreEnvs(origin map[string]string) {
	for env := range origin {
		os.Setenv(env, origin[env])
	}
}

func helperTestFromEnv(t *testing.T, testEnvs map[string]string) {
	cfg := Config{}
	origEnvs := overrideEnvs(testEnvs)
	cfg.FromEnv()
	restoreEnvs(origEnvs)
	if cfg.Scheme != testEnvs["TIDEPOOL_STORE_SCHEME"] {
		t.Errorf("cfg.Scheme not matching env.\nExpected %v got %v\n", testEnvs["TIDEPOOL_STORE_SCHEME"], cfg.Scheme)
	}
	adresses := strings.Split(testEnvs["TIDEPOOL_STORE_ADDRESSES"], ",")
	if len(adresses) != len(cfg.addresses) || cfg.addresses[0] != adresses[0] {
		t.Errorf("cfg.addresses not matching env.\nExpected %v got %v\n", adresses, cfg.addresses)
	}
	if cfg.Username != testEnvs["TIDEPOOL_STORE_USERNAME"] {
		t.Errorf("cfg.Username not matching env.\nExpected %v got %v\n", testEnvs["TIDEPOOL_STORE_USERNAME"], cfg.Username)
	}
	if cfg.Password != testEnvs["TIDEPOOL_STORE_PASSWORD"] {
		t.Errorf("cfg.Password not matching env.\nExpected %v got %v\n", testEnvs["TIDEPOOL_STORE_PASSWORD"], cfg.Password)
	}
	if cfg.Database != testEnvs["TIDEPOOL_STORE_DATABASE"] {
		t.Errorf("cfg.Database not matching env.\nExpected %v got %v\n", testEnvs["TIDEPOOL_STORE_DATABASE"], cfg.Database)
	}
	if cfg.OptParams != testEnvs["TIDEPOOL_STORE_OPT_PARAMS"] {
		t.Errorf("cfg.OptParams not matching env.\nExpected %v got %v\n", testEnvs["TIDEPOOL_STORE_OPT_PARAMS"], cfg.OptParams)
	}
	if cfg.TLS != (testEnvs["TIDEPOOL_STORE_TLS"] == "true") {
		t.Errorf("cfg.TLS not matching env.\nExpected %v got %v\n", testEnvs["TIDEPOOL_STORE_TLS"] == "true", cfg.TLS)
	}
	expectedTimeout, err := strconv.Atoi(testEnvs["TIDEPOOL_STORE_DEFAULT_TIMEOUT"])
	if err != nil {
		expectedTimeout = 2
	}
	if cfg.Timeout != time.Duration(expectedTimeout) * time.Second {
		t.Errorf("cfg.Timeout not matching env.\nExpected %v seconds got %v\n", expectedTimeout, cfg.Timeout)
	}
	expectedWaitConnectionInterval, err := strconv.Atoi(testEnvs["TIDEPOOL_STORE_WAIT_CONNECTION_INTERVAL"])
	if err != nil {
		expectedWaitConnectionInterval = 5
	}
	if cfg.WaitConnectionInterval != time.Duration(expectedWaitConnectionInterval) * time.Second {
		t.Errorf("cfg.WaitConnectionInterval not matching env.\nExpected %v seconds got %v\n", expectedWaitConnectionInterval, cfg.WaitConnectionInterval)
	}
	expectedMaxConnectionAttempts, err := strconv.Atoi(testEnvs["TIDEPOOL_STORE_MAX_CONNECTION_ATTEMPTS"])
	if err != nil {
		expectedMaxConnectionAttempts = 0
	}
	if cfg.MaxConnectionAttempts != int64(expectedMaxConnectionAttempts) {
		t.Errorf("cfg.MaxConnectionAttempts not matching env.\nExpected %v got %v\n", expectedMaxConnectionAttempts, cfg.MaxConnectionAttempts)
	}
	
	if testEnvs["TIDEPOOL_STORE_READ_MODE"] != "" {
		if cfg.ReadPreferences.Mode().String() != testEnvs["TIDEPOOL_STORE_READ_MODE"] {
			t.Errorf("cfg.ReadPreferences mode not matching env.\nExpected %v got %v\n", testEnvs["TIDEPOOL_STORE_READ_MODE"], cfg.ReadPreferences.Mode().String())
		}
		expectedStaleness, err := strconv.Atoi(testEnvs["TIDEPOOL_STORE_MAX_STALENESS"])
		expectedSet := true
		if err != nil {
			expectedSet = false
		}
		staleness, set := cfg.ReadPreferences.MaxStaleness()
		if set != expectedSet {
			t.Errorf("cfg.ReadPreferences MaxStaleness should be set: %v.", expectedSet)
		}
		if  set && staleness != time.Duration(expectedStaleness) * time.Second {
			t.Errorf("cfg.ReadPreferences MaxStaleness not matching env.\nExpected %v seconds got %v\n", expectedStaleness, staleness)
		}
	} else {
		if cfg.ReadPreferences != nil {
			t.Errorf("cfg.ReadPreferences should be nil.\nFound %v", cfg.ReadPreferences)
		}
	}
	
}

func TestFromEnv(t *testing.T) {
	testEnv := map[string]string{
		"TIDEPOOL_STORE_SCHEME": "http",
		"TIDEPOOL_STORE_ADDRESSES": "mongo.dblp.fr",
		"TIDEPOOL_STORE_USERNAME": "diabeloop",
		"TIDEPOOL_STORE_PASSWORD": "superSafePassword",
		"TIDEPOOL_STORE_DATABASE": "dbl-data",
		"TIDEPOOL_STORE_OPT_PARAMS": "&opt1=1&opt2=2",
		"TIDEPOOL_STORE_TLS": "true",
		"TIDEPOOL_STORE_DEFAULT_TIMEOUT": "10",
		"TIDEPOOL_STORE_WAIT_CONNECTION_INTERVAL": "15",
		"TIDEPOOL_STORE_MAX_CONNECTION_ATTEMPTS": "20",
		"TIDEPOOL_STORE_READ_MODE": "secondaryPreferred",
		"TIDEPOOL_STORE_MAX_STALENESS": "120",
	}
	helperTestFromEnv(t, testEnv)
	testEnv = map[string]string{
		"TIDEPOOL_STORE_SCHEME": "http",
		"TIDEPOOL_STORE_ADDRESSES": "mongo.dblp.fr",
		"TIDEPOOL_STORE_USERNAME": "diabeloop",
		"TIDEPOOL_STORE_PASSWORD": "superSafePassword",
		"TIDEPOOL_STORE_DATABASE": "dbl-data",
		"TIDEPOOL_STORE_OPT_PARAMS": "&opt1=1&opt2=2",
	}
	helperTestFromEnv(t, testEnv)
	testEnv = map[string]string{
		"TIDEPOOL_STORE_SCHEME": "http",
		"TIDEPOOL_STORE_ADDRESSES": "mongo.dblp1.fr,mongo.dblp2.fr",
		"TIDEPOOL_STORE_USERNAME": "diabeloop",
		"TIDEPOOL_STORE_PASSWORD": "superSafePassword",
		"TIDEPOOL_STORE_DATABASE": "dbl-data",
		"TIDEPOOL_STORE_OPT_PARAMS": "&opt1=1&opt2=2",
	}
	helperTestFromEnv(t, testEnv)


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
