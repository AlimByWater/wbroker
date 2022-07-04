package dic

import (
	"os"
	"reflect"

	"go.uber.org/dig"
)

// Engine provides additional functionality for dig DI container
type Engine struct {
	container  *dig.Container
	errHandler func(error)
	main       interface{}
}

// New returns engine with specific box
func New(box Box) *Engine {
	e := &Engine{
		container: dig.New(),
		main:      box.Main(),
	}
	asserMain(e.main)
	e.provide(box.Modules())
	return e
}

// SetErrorHandler sets default error handler, error handler called on provide invoke errors
func (e *Engine) SetErrorHandler(h func(error)) *Engine {
	e.errHandler = h
	return e
}

// Run invoke main function from container
func (e Engine) Run() {
	t := reflect.TypeOf(e.main)
	switch t.Kind() {
	case reflect.Func:
		e.handleErr(e.container.Invoke(e.main))
	case reflect.Slice:
		for _, fn := range e.main.([]interface{}) {
			e.handleErr(e.container.Invoke(fn))
		}
	}
}

// Invoke provides abbility to access to DI container
func (e Engine) Invoke(function interface{}) Engine {
	e.handleErr(e.container.Invoke(function))
	return e
}

// Graph write dependencies graph to file
func (e Engine) Graph(filename string) {
	file, err := os.Create(filename)
	e.handleErr(err)
	defer file.Close()
	e.handleErr(dig.Visualize(e.container, file))
}

func (e Engine) provide(m Module) {
	for _, c := range m {
		e.handleErr(e.container.Provide(c.CreateFunc, c.Options...))
	}
}

func (e Engine) handleErr(err error) {
	if err != nil {
		if e.errHandler == nil {
			panic(err)
		}
		e.errHandler(err)
	}
}

func asserMain(f interface{}) {
	t := reflect.TypeOf(f)
	switch t.Kind() {
	case reflect.Func:
	case reflect.Slice:
		for _, fn := range f.([]interface{}) {
			ft := reflect.TypeOf(fn)
			if ft.Kind() != reflect.Func {
				panic("box main slice should contains only functions")
			}
		}
	default:
		panic("box.Main must return function or functions slice, gotten " + t.Kind().String())
	}
}
