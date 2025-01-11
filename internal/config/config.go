package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// loadConfig reads and parses the YAML configuration file
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	// Validate required fields
	if cfg.Azure.TenantID == "" || cfg.Azure.ClientID == "" || cfg.Azure.ClientSecret == "" || cfg.Database.Path == "" {
		return nil, fmt.Errorf("missing required configuration fields")
	}
	return &cfg, nil
}
