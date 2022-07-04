package main

import (
	"flag"
	"wbroker/app"
	dic "wbroker/app/dig"
)

var (
	flagGraph   = flag.Bool("graph", true, "creates dependencies graph")
	flagAppName = flag.String("app-name", "wbroker", "")
	flagConfig  = flag.String("configs-file", "configs/configs.local.yaml", "configs file name")
)

func main() {
	flag.Parse()

	box := dic.NewBox(app.Modules, *flagAppName, *flagConfig)

	app := dic.New(box)

	if *flagGraph {
		app.Graph("graph.dot")
		//return
	}

	app.Run()
}
