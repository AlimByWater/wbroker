package wbroker

import (
	"fmt"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// ProviderOut is dig adapter for created clients
type ProviderOut struct {
	dig.Out

	WBrokerClient WBrokerClient
}

// ProvideClients provides gRPC-client for communications service with default timeouts
func ProvideClients() (ProviderOut, error) {
	return ProvideClientsWithTimeout(0)() // for adding manual dial provider
}

// ProvideClientsWithTimeout provides gRPC-client for communications service with specified timeouts
func ProvideClientsWithTimeout(timeout time.Duration) func() (ProviderOut, error) {
	return func() (ProviderOut, error) {
		dial, err := grpc.Dial(":24005", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return ProviderOut{}, fmt.Errorf("dial: %w", err)
		}

		return ProviderOut{
			WBrokerClient: NewWBrokerClient(dial),
		}, nil
	}
}
