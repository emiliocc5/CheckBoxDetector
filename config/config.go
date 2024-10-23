package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// ServiceConfiguration Struct for load envs.
type ServiceConfiguration struct {
	General struct {
		Environment string `yaml:"ENVIRONMENT" env-default:"local" env-description:"Service Environment"`
	} `yaml:"General"`
	File struct {
		Path string `yaml:"FILE_PATH" env-default:"path_to_file" env-description:"Path where file is located"`
	} `yaml:"File"`
}

// NewConfigService Return filled ServiceConfiguration.
func NewConfigService() *ServiceConfiguration {
	cfg := &ServiceConfiguration{}

	return cfg
}

func (cfg *ServiceConfiguration) LoadFromEnvs() {
	if err := cleanenv.ReadConfig(".config.yml", cfg); err != nil {
		panic(err)
	}
}
