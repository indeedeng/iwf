package service

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type (
	Config struct {
		// Api is the API config
		Api ApiConfig `yaml:"api"`
		// Backend is the service behind, either Cadence or Temporal is required
		Backend Backend `yaml:"backend"`
	}

	ApiConfig struct {
		// Port is the port on which the API service will bind to
		Port int `yaml:"port"`
	}

	Backend struct {
		// Temporal config is the config to connect to Temporal
		Temporal *TemporalConfig `yaml:"temporal"`
	}

	TemporalConfig struct {
		// HostPort to connect to, default to localhost:7233
		HostPort string `yaml:"hostPort"`
		// Namespace to connect to, default to default
		Namespace string `yaml:"namespace"`
	}
)

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	log.Printf("Loading configFile=%v\n", configPath)

	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
