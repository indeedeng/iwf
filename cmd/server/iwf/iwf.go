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
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/client/workflowState"
	temporalimpl "github.com/cadence-oss/iwf-server/service/interpreter/temporalImpl"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"io/ioutil"
	"log"
	"net/http"
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
		startTestWebhookEndpoint()
		break
	case "interpreter-temporal":
		startTemporalWorker()
	default:
		log.Printf("Invalid service: %v", service)
	}
}

func startTemporalWorker() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world", worker.Options{})

	w.RegisterWorkflow(temporalimpl.Interpreter)
	w.RegisterActivity(temporalimpl.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
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

func startTestWebhookEndpoint() {
	http.HandleFunc("/", logIncomingRequest)

	fmt.Printf("Starting server for testing...\n")
	// TODO make api endpoint port configurable
	if err := http.ListenAndServe(":8801", nil); err != nil {
		log.Fatal(err)
	}
}

func logIncomingRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Printf("Path not supported: %v", r.URL.Path)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		var body []byte
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("[Failed to read request body]: %v", err.Error())
		}

		log.Printf("[Test server incoming request]: %v, URL: %v", string(body), r.URL.Path)
		w.WriteHeader(http.StatusOK)
		break
	case "GET":
		// below is for test only, will be removed

		runTestRestApi()
		runTestTemporalWorkflow()
		w.Write([]byte("hello, world"))
		w.WriteHeader(http.StatusOK)
		break
	default:
		log.Printf("Only POST methods are supported.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func runTestTemporalWorkflow() {

	fmt.Println("test temporal workflow")

	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "hello_world_workflowID",
		TaskQueue: "hello-world",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, temporalimpl.Interpreter, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}

func runTestRestApi() {
	client := openapi.NewAPIClient(&openapi.Configuration{})
	req := client.DefaultApi.ApiV1WorkflowStateStartPost(context.Background())
	wfType := "123"
	resp, httpResp, err := req.WorkflowStateStartRequest(openapi.WorkflowStateStartRequest{
		WorkflowType: &wfType,
	}).Execute()
	fmt.Println("test REST API", resp.GetCommandRequest(), httpResp, err)
}
