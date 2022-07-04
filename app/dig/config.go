package dic

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type (
	Configs struct {
		dig.Out

		GRPCServer *GRPCServerConfig `yaml:"grpc"`

		RawConfig []byte `name:"raw_config"`
	}

	// GRPCServerConfig contains grpc server settings
	GRPCServerConfig struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
)

// NewRawConfigFunc returns config provider which using rawconfig
func NewRawConfigFunc(serviceName string, interactive bool, rawcfg []byte) func() (Configs, error) {
	return func() (Configs, error) {
		return newConfig(rawcfg)
	}
}

// NewConfigFunc returns configs constructor func
func NewConfigFunc(fileName string) func() (Configs, error) {
	return func() (Configs, error) {
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			return Configs{}, err
		}
		return newConfig(b)
	}
}

func newConfig(body []byte) (Configs, error) {
	cfg := Configs{}
	err := yaml.Unmarshal(body, &cfg)
	if err != nil {
		return Configs{}, err
	}

	cfg.RawConfig = body

	logrus.Info(cfg.GRPCServer.Port)

	return cfg, nil
}
