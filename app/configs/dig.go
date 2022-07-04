package configs

import dic "wbroker/app/dig"

var Modules = dic.Module{
	{CreateFunc: NewConfig},
}
