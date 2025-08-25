package integ

import (
	"github.com/indeedeng/iwf/config"
)

const testWorkflowServerPort = "9714"
const testIwfServerPort = "9715"

func createTestConfig(testCfg IwfServiceTestConfig) config.Config {
	cfg := config.Config{
		Api: config.ApiConfig{
			Port:           9715,
			MaxWaitSeconds: 12, // use 12 so that we can test it in the waiting test
			WaitForStateCompletionMigration: config.WaitForStateCompletionMigration{
				SignalWithStartOn: "old",
				WaitForOn:         "old",
			},
			QueryWorkflowFailedRetryPolicy: config.QueryWorkflowFailedRetryPolicy{
				InitialIntervalSeconds: 1,
				MaximumAttempts:        10,
			},
		},
		Interpreter: config.Interpreter{
			VerboseDebug:              false,
			FailAtMemoIncompatibility: !testCfg.DisableFailAtMemoIncompatibility,
			InterpreterActivityConfig: config.InterpreterActivityConfig{
				DefaultHeaders: testCfg.DefaultHeaders,
			},
		},
	}
	if testCfg.S3TestThreshold > 0 {
		externalStorage := config.ExternalStorageConfig{
			Enabled: true,
			ThresholdInBytes: testCfg.S3TestThreshold,
			SupportedStorages: []config.SupportedStorage{
				{
					Status: config.StorageStatusActive,
					StorageId: "s3",
					StorageType: "s3",
					S3Endpoint: "http://localhost:9000",
					S3Bucket: "iwf-test-bucket",
					S3Region: "us-east-1",
					S3AccessKey: "minioadmin",
					S3SecretKey: "minioadmin",
				},
			},
		}
		cfg.ExternalStorage = externalStorage
	}
	return cfg
}
