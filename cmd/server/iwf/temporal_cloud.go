package iwf

import (
	"context"
	"crypto/tls"
	"fmt"

	cloudservicev1 "go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NOTE: from https://github.com/temporalio/cloud-samples-go/blob/main/client/temporal/client.go

type TemporalCloudNamespaceClientInput struct {
	// The temporal cloud namespace to connect to (required) for e.g. "prod.a2dd6"
	Namespace string `required:"true"`

	// The auth to use for the client, defaults to local
	Auth AuthType

	MetricsHandler client.MetricsHandler
}

func GetTemporalCloudNamespaceClient(
	ctx context.Context, input *TemporalCloudNamespaceClientInput,
) (client.Client, error) {

	opts := client.Options{
		Namespace: input.Namespace,
	}
	err := input.Auth.apply(&opts)
	if err != nil {
		return nil, err
	}
	
	return client.Dial(opts)
}

type AuthType interface {
	apply(options *client.Options) error
}

// ApiKeyAuth is an implementation of above AuthType
type ApiKeyAuth struct {
	// The api key to use for the client
	APIKey string
}

func (a *ApiKeyAuth) apply(options *client.Options) error {

	c, err := NewConnectionWithAPIKey(a.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create cloud api connection: %w", err)
	}
	resp, err := c.CloudService().GetNamespace(context.Background(), &cloudservicev1.GetNamespaceRequest{
		Namespace: options.Namespace,
	})
	if err != nil {
		return fmt.Errorf("failed to get namespace %s: %w", options.Namespace, err)
	}
	if resp.GetNamespace().GetEndpoints().GetGrpcAddress() == "" {
		return fmt.Errorf("namespace %q has no grpc address", options.Namespace)
	}
	options.HostPort = resp.GetNamespace().GetEndpoints().GetGrpcAddress()
	options.Credentials = client.NewAPIKeyStaticCredentials(a.APIKey)
	options.ConnectionOptions = client.ConnectionOptions{
		TLS: &tls.Config{},
		DialOptions: []grpc.DialOption{
			grpc.WithUnaryInterceptor(
				func(
					ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn,
					invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
				) error {
					return invoker(
						metadata.AppendToOutgoingContext(ctx, "temporal-namespace", options.Namespace),
						method,
						req,
						reply,
						cc,
						opts...,
					)
				},
			),
		},
	}
	return nil
}
