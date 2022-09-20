package integ

import (
	"go.temporal.io/sdk/client"
	"log"
)

func createTemporalClient() client.Client {
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("unable to connect to Temporal %v", err)
	}
	return temporalClient
}
