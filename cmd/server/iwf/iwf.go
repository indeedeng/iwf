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
	"github.com/cadence-oss/iwf-server/service/api"
	temporalimpl "github.com/cadence-oss/iwf-server/service/interpreter/temporalImpl"
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
			Name:  "root, r",
			Value: ".",
			Usage: "root directory of execution environment",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "config",
			Usage: "config dir is a path relative to root, or an absolute path",
		},
		cli.StringFlag{
			Name:  "env, e",
			Value: "development",
			Usage: "runtime environment",
		},
		cli.StringFlag{
			Name:  "zone, az",
			Value: "",
			Usage: "availability zone",
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
			Usage: "start iwf notification service",
			Action: func(c *cli.Context) {
				var wg sync.WaitGroup
				services := getServices(c)

				for _, service := range services {
					wg.Add(1)
					go launchService(service, c)
				}

				wg.Wait()
			},
		},
	}
	return app
}

func launchService(service string, c *cli.Context) {
	switch service {
	case "api":
		router := api.NewRouter()
		// TODO use port number from config
		log.Fatal(router.Run(":8801"))
	case "interpreter-temporal":
		interpreter := temporalimpl.NewInterpreterWorker()
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
