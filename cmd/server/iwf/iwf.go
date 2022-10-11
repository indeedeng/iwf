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
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/api"
	temporalapi "github.com/cadence-oss/iwf-server/service/api/temporal"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"github.com/urfave/cli"
	"go.temporal.io/sdk/client"
	"log"
	"strings"
	"sync"
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

func start(c *cli.Context) {
	configPath := c.GlobalString("config")
	config, err := service.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Unable to load config for path %v because of error %v", configPath, err)
	}
	services := getServices(c)

	// The client is a heavyweight object that should be created once per process.
	var unifiedClient api.UnifiedClient
	if config.Backend.Temporal != nil {
		temporalClient, err := client.Dial(client.Options{
			HostPort:  config.Backend.Temporal.HostPort,
			Namespace: config.Backend.Temporal.Namespace,
		})
		if err != nil {
			log.Fatalf("Unable to connect to Temporal because of error %v", err)
		}
		unifiedClient = temporalapi.NewTemporalClient(temporalClient)

		for _, svcName := range services {
			go launchTemporalService(svcName, config, unifiedClient, temporalClient)
		}
	} else {
		panic("only support Temporal today")
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func launchTemporalService(svcName string, config *service.Config, unifiedClient api.UnifiedClient, temporalClient client.Client) {
	switch svcName {
	case serviceAPI:
		svc := api.NewService(unifiedClient)
		log.Fatal(svc.Run(fmt.Sprintf(":%v", config.Api.Port)))
	case serviceInterpreter:
		interpreter := temporal.NewInterpreterWorker(temporalClient)
		interpreter.Start()
	default:
		log.Printf("Invalid service: %v", svcName)
	}
}

func getServices(c *cli.Context) []string {
	val := strings.TrimSpace(c.String("services"))
	tokens := strings.Split(val, ",")

	if len(tokens) == 0 {
		log.Fatal("No services specified for starting")
	}

	var services []string
	for _, token := range tokens {
		t := strings.TrimSpace(token)
		services = append(services, t)
	}

	return services
}
