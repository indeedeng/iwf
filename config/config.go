package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/uber-go/tally/v4/prometheus"
	temporalWorker "go.temporal.io/sdk/worker"
	cadenceWorker "go.uber.org/cadence/worker"
	"gopkg.in/yaml.v3"
)

const (
	StorageStatusActive   = "active"
	StorageStatusInactive = "inactive"
)

type (
	Config struct {
		// Log is the logging config
		Log Logger `yaml:"log"`
		// Api is the API config
		Api ApiConfig `yaml:"api"`
		// Interpreter is the service behind, either Cadence or Temporal is required
		Interpreter Interpreter `yaml:"interpreter"`
		// ExternalStorage is the external storage config
		ExternalStorage ExternalStorageConfig `yaml:"externalStorage"`
	}

	ExternalStorageConfig struct {
		Enabled bool `yaml:"enabled"`
		// ThresholdInBytes is the size threshold of encodedObject
		// that will be stored by external storage(picking the current active one)
		ThresholdInBytes int `yaml:"thresholdInBytes"`
		// SupportedStorages is the list of supported storage
		// Only one can be active, meaning the one that will be used for writing.
		// The non-active ones are for read only.
		SupportedStorages []SupportedStorage `yaml:"supportedStorages"`
	}

	StorageStatus string

	SupportedStorage struct {
		// Status means whether this storage is active for writing.
		// Only one of the supported storages can be active
		Status StorageStatus
		// StorageId is the id of the external storage, it's used to identify the external storage in the EncodedObject that is stored in the workflow history
		StorageId string `yaml:"storageId"`
		// StorageType is the type of the external storage, currently only s3 is supported
		StorageType string `yaml:"storageType"`
		// S3Endpoint is the endpoint of s3 service
		S3Endpoint string `yaml:"s3Endpoint"`
		// S3Bucket is the bucket name of the S3 storage
		S3Bucket string `yaml:"s3Bucket"`
		// S3Region is the region of the S3 storage
		S3Region string `yaml:"s3Region"`
		// S3AccessKey is the access key of the S3 storage
		S3AccessKey string `yaml:"s3AccessKey"`
		// S3SecretKey is the secret key of the S3 storage
		S3SecretKey string `yaml:"s3SecretKey"`
	}

	ApiConfig struct {
		// Port is the port on which the API service will bind to
		Port           int   `yaml:"port"`
		MaxWaitSeconds int64 `yaml:"maxWaitSeconds"`
		// omitRpcInputOutputInHistory is the flag to omit rpc input/output in history
		// the input/output is only for debugging purpose but could be too expensive to store
		OmitRpcInputOutputInHistory *bool `yaml:"omitRpcInputOutputInHistory"`
		// WaitForStateCompletionMigration is used to control workflowId of the WaitForStateCompletion system/internal workflows
		WaitForStateCompletionMigration WaitForStateCompletionMigration `yaml:"waitForStateCompletionMigration"`
		QueryWorkflowFailedRetryPolicy  QueryWorkflowFailedRetryPolicy  `yaml:"queryWorkflowFailedRetryPolicy"`
	}

	QueryWorkflowFailedRetryPolicy struct {
		// defaults to 1
		InitialIntervalSeconds int `yaml:"initialIntervalSeconds"`
		// defaults to 5
		MaximumAttempts int `yaml:"maximumAttempts"`
	}

	WaitForStateCompletionMigration struct {
		// expected values: old/both/new; defaults to 'old'
		SignalWithStartOn string `yaml:"signalWithStartOn"`
		// expected values: old/new; defaults to 'old'
		WaitForOn string `yaml:"waitForOn"`
	}

	Interpreter struct {
		// Temporal config is the config to connect to Temporal
		Temporal *TemporalConfig `yaml:"temporal"`
		// Cadence config is the config to connect to Cadence
		Cadence                   *CadenceConfig            `yaml:"cadence"`
		DefaultWorkflowConfig     *iwfidl.WorkflowConfig    `json:"defaultWorkflowConfig"`
		InterpreterActivityConfig InterpreterActivityConfig `yaml:"interpreterActivityConfig"`
		VerboseDebug              bool
		FailAtMemoIncompatibility bool
	}

	TemporalConfig struct {
		// HostPort to connect to, default to localhost:7233
		HostPort string `yaml:"hostPort"`
		// API key to connect to Temporal Cloud, default to empty
		CloudAPIKey string `yaml:"cloudAPIKey"`
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
		ApiServiceAddress                  string                              `json:"serviceAddress"`
		DumpWorkflowInternalActivityConfig *DumpWorkflowInternalActivityConfig `json:"dumpWorkflowInternalActivityConfig"`
		DefaultHeaders                     map[string]string                   `json:"defaultHeaders"`
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

var DefaultWorkflowConfig = &iwfidl.WorkflowConfig{
	ContinueAsNewThreshold: iwfidl.PtrInt32(100),
}

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

func (c Config) GetApiServiceAddressWithDefault() string {
	if c.Interpreter.InterpreterActivityConfig.ApiServiceAddress != "" {
		return c.Interpreter.InterpreterActivityConfig.ApiServiceAddress
	}
	return fmt.Sprintf("http://localhost:%v", c.Api.Port)
}

func (c Config) GetSignalWithStartOnWithDefault() string {
	if c.Api.WaitForStateCompletionMigration.SignalWithStartOn != "" {
		return c.Api.WaitForStateCompletionMigration.SignalWithStartOn
	}
	return "old"
}

func (c Config) GetWaitForOnWithDefault() string {
	if c.Api.WaitForStateCompletionMigration.WaitForOn != "" {
		return c.Api.WaitForStateCompletionMigration.WaitForOn
	}
	return "old"
}

func QueryWorkflowFailedRetryPolicyWithDefaults(retryPolicy *QueryWorkflowFailedRetryPolicy) QueryWorkflowFailedRetryPolicy {
	var rp QueryWorkflowFailedRetryPolicy

	if retryPolicy != nil && retryPolicy.InitialIntervalSeconds != 0 {
		rp.InitialIntervalSeconds = retryPolicy.InitialIntervalSeconds
	} else {
		rp.InitialIntervalSeconds = 1
	}

	if retryPolicy != nil && retryPolicy.MaximumAttempts != 0 {
		rp.MaximumAttempts = retryPolicy.MaximumAttempts
	} else {
		rp.MaximumAttempts = 5
	}

	return rp
}
