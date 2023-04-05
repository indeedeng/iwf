package config

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/uber-go/tally/v4/prometheus"
	temporalWorker "go.temporal.io/sdk/worker"
	cadenceWorker "go.uber.org/cadence/worker"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type (
	Config struct {
		// Log is the logging config
		Log Logger `yaml:"log"`
		// Api is the API config
		Api ApiConfig `yaml:"api"`
		// Interpreter is the service behind, either Cadence or Temporal is required
		Interpreter Interpreter `yaml:"interpreter"`
	}

	ApiConfig struct {
		// Port is the port on which the API service will bind to
		Port int `yaml:"port"`
	}

	Interpreter struct {
		// Temporal config is the config to connect to Temporal
		Temporal *TemporalConfig `yaml:"temporal"`
		// Cadence config is the config to connect to Cadence
		Cadence                   *CadenceConfig `yaml:"cadence"`
		DefaultWorkflowConfig     iwfidl.WorkflowConfig
		InterpreterActivityConfig InterpreterActivityConfig
		VerboseDebug              bool
	}

	TemporalConfig struct {
		// HostPort to connect to, default to localhost:7233
		HostPort string `yaml:"hostPort"`
		// Namespace to connect to, default to default
		Namespace string `yaml:"namespace"`
		// Prometheus is configuring the metric exposer
		Prometheus    *prometheus.Configuration `yaml:"prometheus"`
		WorkerOptions *temporalWorker.Options
	}

	CadenceConfig struct {
		// HostPort to connect to, default to 127.0.0.1:7833
		HostPort string `yaml:"hostPort"`
		// Domain to connect to, default to default
		Domain        string `yaml:"domain"`
		WorkerOptions *cadenceWorker.Options
	}

	InterpreterActivityConfig struct {
		// ApiServiceAddress is the address that core engine workflow talks to API service
		// It's used in DumpWorkflowInternal activity for continueAsNew
		// default is http://localhost:ApiConfig.Port
		ApiServiceAddress                  string `json:"serviceAddress"`
		DumpWorkflowInternalActivityConfig *DumpWorkflowInternalActivityConfig
	}

	DumpWorkflowInternalActivityConfig struct {
		StartToCloseTimeout time.Duration
		RetryPolicy         *iwfidl.RetryPolicy
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

func GetApiServiceAddressWithDefault(config Config) string {
	if config.Interpreter.InterpreterActivityConfig.ApiServiceAddress != "" {
		return config.Interpreter.InterpreterActivityConfig.ApiServiceAddress
	}
	return fmt.Sprintf("http://localhost:%v", config.Api.Port)
}
