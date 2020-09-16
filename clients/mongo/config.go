package mongo

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"time"

	common "github.com/tidepool-org/go-common"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config of the mongo database
type Config struct {
	Scheme                 string `json:"scheme"`
	addresses              []string
	TLS                    bool                          `json:"tls"`
	Database               string                        `json:"database"`
	Username               string                        `json:"-"`
	Password               string                        `json:"-"`
	Timeout                time.Duration                 `json:"timeout"`
	OptParams              string                        `json:"optParams"`
	WaitConnectionInterval time.Duration                 `json:"waitConnectionInterval"`
	MaxConnectionAttempts  int64                         `json:"maxConnectionAttempts"`
	Indexes                map[string][]mongo.IndexModel `json:"indexes"`
}

// FromEnv read the mongo config from the environment variables
func (config *Config) FromEnv() {
	config.Scheme, _ = os.LookupEnv("TIDEPOOL_STORE_SCHEME")
	stringAddresses, _ := os.LookupEnv("TIDEPOOL_STORE_ADDRESSES")
	addresses := []string{}
	for _, str := range strings.Split(stringAddresses, ",") {
		if str = strings.TrimSpace(str); str != "" {
			addresses = append(addresses, str)
		}
	}
	config.addresses = addresses
	config.Username, _ = os.LookupEnv("TIDEPOOL_STORE_USERNAME")
	config.Password, _ = os.LookupEnv("TIDEPOOL_STORE_PASSWORD")
	config.Database, _ = os.LookupEnv("TIDEPOOL_STORE_DATABASE")
	config.OptParams, _ = os.LookupEnv("TIDEPOOL_STORE_OPT_PARAMS")
	tls, found := os.LookupEnv("TIDEPOOL_STORE_TLS")
	config.TLS = found && tls == "true"

	defaultTimeout := common.GetEnvironmentInt64("TIDEPOOL_STORE_DEFAULT_TIMEOUT", 2)
	config.Timeout = time.Duration(defaultTimeout) * time.Second
	waitConnectionInterval := common.GetEnvironmentInt64("TIDEPOOL_STORE_WAIT_CONNECTION_INTERVAL", 5)
	config.WaitConnectionInterval = time.Duration(waitConnectionInterval) * time.Second
	config.MaxConnectionAttempts = common.GetEnvironmentInt64("TIDEPOOL_STORE_MAX_CONNECTION_ATTEMPTS", 0)

}

func (config *Config) toConnectionString() (string, error) {
	if len(config.addresses) == 0 {
		config.addresses = []string{"localhost"}
	}
	for _, address := range config.addresses {
		if address == "" {
			return "", errors.New("address is missing")
		} else if _, err := url.Parse(address); err != nil {
			return "", errors.New("address is invalid")
		}
	}
	if config.Database == "" {
		return "", errors.New("database is missing")
	}

	var url string
	if config.Scheme != "" {
		url += config.Scheme + "://"
	} else {
		url += "mongodb://"
	}

	if config.Username != "" {
		url += config.Username
		if config.Password != "" {
			url += ":"
			url += config.Password
		}
		url += "@"
	}
	url += strings.Join(config.addresses, ",")
	url += "/"
	url += config.Database
	if config.TLS {
		url += "?ssl=true"
	} else {
		url += "?ssl=false"
	}
	if config.OptParams != "" {
		url += "&" + config.OptParams
	}

	return url, nil
}
