package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Config is the main configuration
type Config struct {
	HTTP HTTP
	RPC  RPC
}

// HTTP configuration
type HTTP struct {
	Port int
}

// RPC configuration
type RPC struct {
	Port int
}

// FromFile reads the configuration from a file
func FromFile(path string, to interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, to)
}
