// Copyright (c) 2021 Cadence workflow OSS organization
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package iwf

import (
	"fmt"
	isvc "github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/api"
	cadenceapi "github.com/indeedeng/iwf/service/api/cadence"
	temporalapi "github.com/indeedeng/iwf/service/api/temporal"
	"github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/common/log"
	"github.com/indeedeng/iwf/service/common/log/loggerimpl"
	"github.com/indeedeng/iwf/service/common/log/tag"
	"github.com/indeedeng/iwf/service/interpreter/cadence"
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	apiv1 "github.com/uber/cadence-idl/go/proto/api/v1"
	"github.com/urfave/cli"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/converter"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	cclient "go.uber.org/cadence/client"
	"go.uber.org/cadence/compatibility"
	"go.uber.org/cadence/encoded"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/grpc"
	rawLog "log"
	"strings"
	"sync"
	"time"
)

const serviceAPI = "api"
const serviceInterpreter = "interpreter"

// BuildCLI is the main entry point for the iwf server
func BuildCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "iwf service"
	app.Usage = "iwf service"
	app.Version = "beta"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config/development.yaml",
			Usage: "config path is a path relative to root, or an absolute path",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{""},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "services",
					Value: fmt.Sprintf("%s, %s", serviceAPI, serviceInterpreter),
					Usage: "start services/components in this project",
				},
			},
			Usage:  "start iwf notification service",
			Action: start,
		},
	}
	return app
}

const DefaultCadenceDomain = "default"
const DefaultCadenceHostPort = "127.0.0.1:7833"

func start(c *cli.Context) {
	configPath := c.GlobalString("config")
	config, err := config.NewConfig(configPath)
	if err != nil {
		rawLog.Fatalf("Unable to load config for path %v because of error %v", configPath, err)
	}
	zapLogger, err := config.Log.NewZapLogger()
	if err != nil {
		rawLog.Fatalf("Unable to create a new zap logger %v", err)
	}
	logger := loggerimpl.NewLogger(zapLogger)

	services := getServices(c)

	// The client is a heavyweight object that should be created once per process.
	var unifiedClient api.UnifiedClient
	if config.Interpreter.Temporal != nil {
		var metricHandler client.MetricsHandler
		if config.Interpreter.Temporal.Prometheus != nil {
			pscope := newPrometheusScope(*config.Interpreter.Temporal.Prometheus, logger)
			metricHandler = sdktally.NewMetricsHandler(pscope)
		}

		temporalClient, err := client.Dial(client.Options{
			HostPort:       config.Interpreter.Temporal.HostPort,
			Namespace:      config.Interpreter.Temporal.Namespace,
			MetricsHandler: metricHandler,
		})
		if err != nil {
			rawLog.Fatalf("Unable to connect to Temporal because of error %v", err)
		}
		unifiedClient = temporalapi.NewTemporalClient(temporalClient, config.Interpreter.Temporal.Namespace, converter.GetDefaultDataConverter(), false)

		for _, svcName := range services {
			go launchTemporalService(svcName, *config, unifiedClient, temporalClient, logger)
		}
	} else if config.Interpreter.Cadence != nil {
		hostPort := DefaultCadenceHostPort
		domain := DefaultCadenceDomain
		if config.Interpreter.Cadence.HostPort != "" {
			hostPort = config.Interpreter.Cadence.HostPort
		}
		if config.Interpreter.Cadence.Domain != "" {
			domain = config.Interpreter.Cadence.Domain
		}
		serviceClient, closeFunc, err := BuildCadenceServiceClient(hostPort)
		if err != nil {
			rawLog.Fatalf("Unable to connect to Cadence because of error %v", err)
		}
		cadenceClient, err := BuildCadenceClient(serviceClient, domain)
		if err != nil {
			rawLog.Fatalf("Unable to connect to Cadence because of error %v", err)
		}
		unifiedClient = cadenceapi.NewCadenceClient(domain, cadenceClient, serviceClient, encoded.GetDefaultDataConverter(), closeFunc)

		for _, svcName := range services {
			go launchCadenceService(svcName, *config, unifiedClient, serviceClient, domain, closeFunc, logger)
		}
	} else {
		panic("must provide either Cadence or Temporal config")
	}

	// TODO improve the waiting with process signal
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func launchTemporalService(svcName string, config config.Config, unifiedClient api.UnifiedClient, temporalClient client.Client, logger log.Logger) {
	switch svcName {
	case serviceAPI:
		svc := api.NewService(config, unifiedClient, logger.WithTags(tag.Service(svcName)))
		rawLog.Fatal(svc.Run(fmt.Sprintf(":%v", config.Api.Port)))
	case serviceInterpreter:
		interpreter := temporal.NewInterpreterWorker(config, temporalClient, isvc.TaskQueue, converter.GetDefaultDataConverter())
		interpreter.Start()
	default:
		rawLog.Fatalf("Invalid service: %v", svcName)
	}
}

func launchCadenceService(
	svcName string,
	config config.Config,
	unifiedClient api.UnifiedClient,
	service workflowserviceclient.Interface,
	domain string,
	closeFunc func(),
	logger log.Logger) {
	switch svcName {
	case serviceAPI:
		svc := api.NewService(config, unifiedClient, logger.WithTags(tag.Service(svcName)))
		rawLog.Fatal(svc.Run(fmt.Sprintf(":%v", config.Api.Port)))
	case serviceInterpreter:
		interpreter := cadence.NewInterpreterWorker(config, service, domain, isvc.TaskQueue, closeFunc)
		interpreter.Start()
	default:
		rawLog.Fatalf("Invalid service: %v", svcName)
	}
}

func getServices(c *cli.Context) []string {
	val := strings.TrimSpace(c.String("services"))
	tokens := strings.Split(val, ",")

	if len(tokens) == 0 {
		rawLog.Fatal("No services specified for starting")
	}

	var services []string
	for _, token := range tokens {
		t := strings.TrimSpace(token)
		services = append(services, t)
	}

	return services
}

const _cadenceFrontendService = "cadence-frontend"
const _cadenceClientName = "cadence-client"

func BuildCadenceClient(service workflowserviceclient.Interface, domain string) (cclient.Client, error) {
	return cclient.NewClient(
		service,
		domain,
		&cclient.Options{
			FeatureFlags: cclient.FeatureFlags{
				WorkflowExecutionAlreadyCompletedErrorEnabled: true,
			},
		}), nil
}

func BuildCadenceServiceClient(hostPort string) (workflowserviceclient.Interface, func(), error) {

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: _cadenceClientName,
		Outbounds: yarpc.Outbounds{
			_cadenceFrontendService: {Unary: grpc.NewTransport().NewSingleOutbound(hostPort)},
		},
	})

	if dispatcher != nil {
		if err := dispatcher.Start(); err != nil {
			rawLog.Fatal("Failed to create outbound transport channel", err)
		}
	}

	if dispatcher == nil {
		rawLog.Fatal("No RPC dispatcher provided to create a connection to Cadence Service")
	}

	clientConfig := dispatcher.ClientConfig(_cadenceFrontendService)
	return compatibility.NewThrift2ProtoAdapter(
			apiv1.NewDomainAPIYARPCClient(clientConfig),
			apiv1.NewWorkflowAPIYARPCClient(clientConfig),
			apiv1.NewWorkerAPIYARPCClient(clientConfig),
			apiv1.NewVisibilityAPIYARPCClient(clientConfig),
		), func() {
			dispatcher.Stop()
		}, nil
}

// tally sanitizer options that satisfy Prometheus restrictions.
// This will rename metrics at the tally emission level, so metrics name we
// use maybe different from what gets emitted. In the current implementation
// it will replace - and . with _
var (
	safeCharacters = []rune{'_'}

	sanitizeOptions = tally.SanitizeOptions{
		NameCharacters: tally.ValidCharacters{
			Ranges:     tally.AlphanumericRange,
			Characters: safeCharacters,
		},
		KeyCharacters: tally.ValidCharacters{
			Ranges:     tally.AlphanumericRange,
			Characters: safeCharacters,
		},
		ValueCharacters: tally.ValidCharacters{
			Ranges:     tally.AlphanumericRange,
			Characters: safeCharacters,
		},
		ReplacementCharacter: tally.DefaultReplacementCharacter,
	}
)

func newPrometheusScope(c prometheus.Configuration, logger log.Logger) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				logger.Error("error in prometheus reporter", tag.Error(err))
			},
		},
	)
	if err != nil {
		logger.Fatal("error creating prometheus reporter", tag.Error(err))
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sanitizeOptions,
		Prefix:          "temporal_samples",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)

	logger.Info("prometheus metrics scope created")
	return scope
}
