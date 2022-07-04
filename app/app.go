package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	grpcOrig "google.golang.org/grpc"
	"net"
	"wbroker/app/api/grpc"
	"wbroker/app/broker"
	configs "wbroker/app/configs"
	dic "wbroker/app/dig"
)

var Modules = grpc.Module.
	Append(configs.Modules).
	Append(broker.Modules).
	Append(dic.Module{
		{CreateFunc: NewApp},
	})

type App struct {
	grpcServer *grpcOrig.Server
	cfg        *configs.Configs
}

// NewApp creates consul service registrars and returns new app instance
func NewApp(grpc *grpcOrig.Server, cfg *configs.Configs) dic.App {
	return &App{
		grpcServer: grpc,
		cfg:        cfg,
	}
}

// Run is app main function
func (a *App) Run(ctx context.Context) error {
	logrus.Info("Starting app...")

	grpcLis, err := net.Listen("tcp", addr(a.cfg.GRPCServer.Host, a.cfg.GRPCServer.Port))
	if err != nil {
		return err
	}

	go func() {
		if err := a.grpcServer.Serve(grpcLis); err != nil {
			if err != grpcOrig.ErrServerStopped {
				panic(fmt.Sprintf("grpc: Serve return unexpected error: %s", err))
			}
		}
	}()

	logrus.Info("App successfully started 🐣")

	<-ctx.Done()

	return nil
}

func addr(host string, port int) string {
	if len(host) == 0 {
		return fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("%s:%d", host, port)
}
