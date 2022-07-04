package grpc

import (
	wb "github.com/AlimByWater/wbroker/external/grpc/wbroker"
	"go.uber.org/dig"
	dic "wbroker/app/dig"
)

// Module is api module
var Module = dic.Module{
	{CreateFunc: NewServer},
	//{CreateFunc: grpc.NewServer},
	{CreateFunc: Adapter},
}

type (
	// AdapterIn DI container
	AdapterIn struct {
		dig.In

		GRPCServer *Server
	}

	// AdapterOut DI container
	AdapterOut struct {
		dig.Out

		GRPCServer dic.Definition `group:"grpc_impl"`
	}
)

// Adapter for external types
func Adapter(in AdapterIn) AdapterOut {
	return AdapterOut{
		GRPCServer: dic.Definition{
			Description:    &wb.WBroker_ServiceDesc,
			Implementation: in.GRPCServer,
		},
	}
}
