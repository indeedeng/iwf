package integ

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/cmd/server/iwf"
	"github.com/indeedeng/iwf/service"
	"go.temporal.io/sdk/client"
	"log"
	"testing"
	"time"
)

// TODO move starting Cadence/Temporal workflow worker and iwf workflow handler here
func TestMain(m *testing.M) {
	var err error

	if *temporalIntegTest {
		var temporalClient client.Client
		for i := 0; i < *dependencyWaitSeconds; i++ {
			fmt.Println("wait for Temporal to be up...last err: ", err)
			time.Sleep(time.Second)
			temporalClient, err = client.Dial(client.Options{
				HostPort:  *temporalHostPort,
				Namespace: testNamespace,
			})
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Fatalf("unable to connect to Temporal %v", err)
		}
		temporalClient.Close()
		fmt.Println("connected to Temporal namespace")
	}

	if *cadenceIntegTest {
		for i := 0; i < *dependencyWaitSeconds; i++ {
			fmt.Println("wait for Cadence to be up...last err: ", err)
			time.Sleep(time.Second)

			_, _, err = iwf.BuildCadenceServiceClient(iwf.DefaultCadenceHostPort)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Fatalf("cannot connnect to Cadence service%v", err)
		}
		fmt.Println("connected to Cadence service")

		serviceClient, closeFunc, err := iwf.BuildCadenceServiceClient(iwf.DefaultCadenceHostPort)
		defer closeFunc()
		for i := 0; i < *dependencyWaitSeconds; i++ {
			fmt.Println("wait for Cadence domain/Search attributes to be ready...")
			time.Sleep(time.Second)
			resp, err := serviceClient.GetSearchAttributes(context.Background())
			ready := false
			if err == nil {
				for key, _ := range resp.GetKeys() {
					// NOTE: this is the last one we registered in init-ci-cadence.sh
					if key == service.SearchAttributeIwfWorkflowType {
						ready = true
						break
					}
				}
			}
			if ready {
				break
			}
		}
		if err != nil {
			log.Fatalf("Cadence service is not ready %v", err)
		}
		fmt.Println("Cadence service is now ready")
	}
}
