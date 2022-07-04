package dic

import (
	"context"
	"github.com/sirupsen/logrus"
)

// TODO implement graceful shutdown
// TODO implement box for testing

var defaultModules = Module{
	{CreateFunc: RegisterGRPC},
	{CreateFunc: NewContext},
}

type (
	// WBroker box provides minimal required modules for initing
	WBroker struct {
		modules  Module
		appName  string
		cfgFile  string
		exitCode int
	}

	// App interface define app interface required for box
	App interface {
		Run(context.Context) error
	}
)

// NewBox returns new wbroker box
func NewBox(appModules Module, appName, cfgFile string) Box {
	b := &WBroker{
		appName: appName,
		cfgFile: cfgFile,
	}

	b.modules = appModules.
		Append(defaultModules).
		Append(Module{
			{CreateFunc: NewConfigFunc(b.cfgFile)},
		})

	return b
}

// Modules returns all modules: box + app modules
func (w *WBroker) Modules() Module {
	return w.modules
}

// Main returns function or functions slice for ordered call
func (w *WBroker) Main() interface{} {
	return []interface{}{
		w.mainFunc(),
	}
}

func (w *WBroker) mainFunc() interface{} {
	return func(ctx context.Context, app App) {
		if err := app.Run(ctx); err != nil {
			logrus.Errorf("Failed run app: %s", err.Error())
			// for most cases app should be restarted
			// let's exit with non-zero code
			w.exitCode = 3
		}
	}
}

func NewContext() context.Context {
	return context.Background()
}
