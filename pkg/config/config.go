package config

import (
	"api-golang-kubernetes/pkg/log"
	"encoding/json"
	"io/ioutil"
)

// Config is the only one instance holding configuration
// of this service.
var config *AppConfig

// AppConfig is a structure into which config file
// (e.g., config/config.json) is loaded.
type AppConfig struct {
	Logging struct {
		Enable bool   `json:"Enable"`
		Level  string `json:"Level"`
	} `json:"Logging"`

	GracefulTermTimeMillis int64
}

// GetInstance returns the pointer to the singleton instance of Config
func GetInstance() *AppConfig {
	if config == nil {
		config = &AppConfig{}
	}
	return config
}

// Load reads config file (e.g., configs/config.json) and
// unmarshalls JSON string in it into Config structure
func (AppConfig) Load(fname string) bool {
	log.D("Load config from the file \"" + fname + "\".")

	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.E("%v", err)
		return false
	}

	errCode := json.Unmarshal(b, &config)
	log.D("config: %v , err: %v", config, errCode)

	return true
}
