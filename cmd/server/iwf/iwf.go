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
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"go.temporal.io/sdk/client"
	"log"
	"strings"
	"sync"

	"github.com/urfave/cli"
)

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
					Value: "api, interpreter-temporal",
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

	// The client is a heavyweight object that should be created once per process.
	temporalClient, err := client.Dial(client.Options{
		HostPort:  config.Temporal.HostPort,
		Namespace: config.Temporal.Namespace,
	})
	if err != nil {
		log.Fatalf("Unable to connect to Temporal because of error %v", err)
	}

	var wg sync.WaitGroup
	services := getServices(c)

	for _, service := range services {
		wg.Add(1)
		go launchService(service, config, temporalClient, c)
	}

	wg.Wait()
}
func launchService(service string, config *service.Config, temporalClient client.Client, c *cli.Context) {
	switch service {
	case "api":
		svc := api.NewService(temporalClient)
		log.Fatal(svc.Run(fmt.Sprintf(":%v", config.Api.Port)))
	case "interpreter-temporal":
		interpreter := temporal.NewInterpreterWorker(temporalClient)
		interpreter.Start()
	default:
		log.Printf("Invalid service: %v", service)
	}
}

func getServices(c *cli.Context) []string {
	val := strings.TrimSpace(c.String("services"))
	tokens := strings.Split(val, ",")

	if len(tokens) == 0 {
		log.Fatal("No services specified for starting")
	}

	services := []string{}
	for _, token := range tokens {
		t := strings.TrimSpace(token)
		services = append(services, t)
	}

	return services
}
