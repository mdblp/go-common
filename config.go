package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func LoadConfig(filenames []string, obj interface{}) error {
	for _, filename := range filenames {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}

		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bytes, obj); err != nil {
			return err
		}
	}
	return nil
}

func LoadEnvironmentConfig(envVars []string, obj interface{}) error {
	for _, envVar := range envVars {
		envValue := os.Getenv(envVar)
		if envValue == "" {
			return fmt.Errorf("%s not found", envVar)
		}

		if err := json.Unmarshal([]byte(envValue), obj); err != nil {
			return fmt.Errorf("%s errored: %s", envVar, err.Error())
		}
	}
	return nil
}

// GetEnvironmentInt64 return the int value from the env, used the default provided if not found
func GetEnvironmentInt64(envVar string, defaultValue int64) int64 {
	stringValue, found := os.LookupEnv(envVar)
	var intValue int64
	var err error
	if found {
		intValue, err = strconv.ParseInt(stringValue, 10, 0)
	}
	if !found || err != nil {
		intValue = defaultValue
	}
	return intValue
}
