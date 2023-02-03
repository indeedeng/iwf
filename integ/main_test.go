package integ

import (
	"context"
	"flag"
	"fmt"
	"github.com/indeedeng/iwf/cmd/server/iwf"
	"github.com/indeedeng/iwf/service/common/ptr"
	"go.temporal.io/sdk/client"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/.gen/go/shared"
	"log"
	"os"
	"testing"
	"time"
)

// TODO move starting Cadence/Temporal workflow worker and iwf workflow handler here
func TestMain(m *testing.M) {
	flag.Parse()
	var err error

	if *temporalIntegTest {
		var temporalClient client.Client
		for i := 0; i < *dependencyWaitSeconds; i++ {
			temporalClient, err = client.Dial(client.Options{
				HostPort:  *temporalHostPort,
				Namespace: testNamespace,
			})
			if err != nil {
				fmt.Println("wait for Temporal to be up...last err: ", err)
				time.Sleep(time.Second)

			} else {
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
			_, _, err = iwf.BuildCadenceServiceClient(iwf.DefaultCadenceHostPort)
			if err != nil {
				fmt.Println("wait for Cadence to be up...last err: ", err)
				time.Sleep(time.Second)
			} else {
				break
			}
		}
		if err != nil {
			log.Fatalf("cannot connnect to Cadence service%v", err)
		}
		fmt.Println("connected to Cadence service")

		var closeFunc func()
		var serviceClient workflowserviceclient.Interface
		serviceClient, closeFunc, err = iwf.BuildCadenceServiceClient(iwf.DefaultCadenceHostPort)
		defer closeFunc()
		for i := 0; i < *dependencyWaitSeconds; i++ {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
			_, err = serviceClient.DescribeDomain(ctx, &shared.DescribeDomainRequest{
				Name: ptr.Any(iwf.DefaultCadenceDomain),
			})
			if err != nil {
				fmt.Println("wait for Cadence domain to be ready...", err)
				time.Sleep(time.Second)
			} else {
				break
			}
		}
		if err != nil {
			log.Fatalf("Cadence service is not ready %v", err)
		}
		fmt.Println("Cadence service is now ready")

		code := m.Run()
		fmt.Println("finished running integ test with status code", code)
		os.Exit(code)
	}
}
