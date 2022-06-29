package grpc

import (
	"testing"
	"wbroker/app/broker"
)

func TestGRPCServer(t *testing.T) Server {
	return Server{
		broker: broker.NewBroker(),
	}
}
