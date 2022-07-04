package broker

import (
	"go.uber.org/dig"
	dic "wbroker/app/dig"
)

var Modules = dic.Module{
	{CreateFunc: NewBroker},
	{CreateFunc: Adapter},
}

type (
	// AdapterIn DI container
	AdapterIn struct {
		dig.In

		Broker *ActualBroker
	}

	// AdapterOut DI container
	AdapterOut struct {
		dig.Out

		Service Service
	}
)

// Adapter for external types
func Adapter(in AdapterIn) AdapterOut {
	return AdapterOut{
		Service: in.Broker,
	}
}
