package dic

import "go.uber.org/dig"

type (
	// Provider should be function
	Provider struct {
		CreateFunc interface{}
		Options    []dig.ProvideOption
	}

	// Module is logic set of providers
	Module []Provider

	// Box is composite of modules and invoke func
	Box interface {
		// Main is first invoke function or ordered slice of functions for invoke
		Main() interface{}
		// Modules is full list of providers
		Modules() Module
	}
)

// Append other module to target module
func (m Module) Append(o Module) Module {
	m = append(m, o...)
	return m
}
