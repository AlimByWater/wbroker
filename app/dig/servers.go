package dic

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	// Definition represents grpc implementation information
	Definition struct {
		Description    *grpc.ServiceDesc // Could be found in generated _grpc.pb.go file
		Implementation interface{}       // Server that embeds Unimplemented***Server from generated _grpc.pb.go file
	}

	GRPCServerParams struct {
		dig.In

		// GRPC configs for gRPC-server
		GRPC *GRPCServerConfig `optional:"true"`
		// GRPCServer common gRPC-server
		GPRCServer *grpc.Server
		// GRPCDefinitions grpc servers implementations
		GRPCDefinitions []Definition `group:"grpc_impl"`
	}
)

func RegisterGRPC(p GRPCServerParams) interface{} {
	if p.GRPC != nil {
		logrus.Infof("Create gRPC server, port: %d", p.GRPC.Port)
		for _, i := range p.GRPCDefinitions {
			p.GPRCServer.RegisterService(i.Description, i.Implementation)
		}

		reflection.Register(p.GPRCServer)

		return p.GPRCServer
	}

	return p.GPRCServer
}
