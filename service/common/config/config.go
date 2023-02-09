package config

import (
	"github.com/uber-go/tally/v4/prometheus"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type (
	Config struct {
		// Log is the logging config
		Log Logger `yaml:"log"`
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
		// Cadence config is the config to connect to Cadence
		Cadence *CadenceConfig `yaml:"cadence"`
	}

	TemporalConfig struct {
		// HostPort to connect to, default to localhost:7233
		HostPort string `yaml:"hostPort"`
		// Namespace to connect to, default to default
		Namespace string `yaml:"namespace"`
		// Prometheus is configuring the metric exposer
		Prometheus *prometheus.Configuration `yaml:"prometheus"`
	}

	CadenceConfig struct {
		// HostPort to connect to, default to 127.0.0.1:7833
		HostPort string `yaml:"hostPort"`
		// Domain to connect to, default to default
		Domain string `yaml:"domain"`
		// DisableSearchAttributes will not use system search attributes
		// this is for Cadence service without advanced visibility because of
		// https://github.com/uber/cadence/issues/5085
		DisableSystemSearchAttributes bool `yaml:"disableSystemSearchAttributes"`
	}

	// Logger contains the config items for logger
	Logger struct {
		// Stdout is true then the output needs to goto standard out
		// By default this is false and output will go to standard error
		Stdout bool `yaml:"stdout"`
		// Level is the desired log level
		Level string `yaml:"level"`
		// OutputFile is the path to the log output file
		// Stdout must be false, otherwise Stdout will take precedence
		OutputFile string `yaml:"outputFile"`
		// LevelKey is the desired log level, defaults to "level"
		LevelKey string `yaml:"levelKey"`
		// Encoding decides the format, supports "console" and "json".
		// "json" will print the log in JSON format(better for machine), while "console" will print in plain-text format(more human friendly)
		// Default is "json"
		Encoding string `yaml:"encoding"`
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
