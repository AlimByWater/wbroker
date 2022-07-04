package configs

import (
	"go.uber.org/dig"
	"gopkg.in/yaml.v3"
)

type (
	Configs struct {
		GRPCServer *GRPCServerConfig `yaml:"grpc"`
	}

	// GRPCServerConfig contains grpc server settings
	GRPCServerConfig struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}

	// NewConfigParams contains raw configs.
	NewConfigParams struct {
		dig.In

		RawConfig []byte `name:"raw_config"`
	}
)

// NewConfig returns configs func
func NewConfig(p NewConfigParams) (*Configs, error) {
	cfg := &Configs{}
	err := yaml.Unmarshal(p.RawConfig, cfg)
	return cfg, err
}
