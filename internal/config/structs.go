package config

import "time"

type Config struct {
	Azure struct {
		TenantID     string `yaml:"tenant_id"`
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"azure"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	ServiceDefinitions struct {
		Path string `yaml:"path"`
	} `yaml:"serviceDefinitions"`
	Scan struct {
		Interval time.Duration `yaml:"interval"`
	} `yaml:"scan"`
}
