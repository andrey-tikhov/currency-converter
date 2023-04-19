package config

import (
	"go.uber.org/config"
)

func New() (config.Provider, error) {
	configFile := config.File("config/base.yaml")
	provider, err := config.NewYAML(configFile)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

type RussiaCBConfig struct {
	APIURL   string `yaml:"api_url,omitempty"`
	Timezone string `yaml:"timezone,omitempty"`
}

type ThailandCBConfig struct {
	APIURL   string `yaml:"api_url,omitempty"`
	Timezone string `yaml:"timezone,omitempty"`
}

type Defaults struct {
	DefaultCB string `yaml:"default_cb"`
}
