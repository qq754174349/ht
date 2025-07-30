package config

import "github.com/qq754174349/ht/ht-frame/autoconfigure"

var config *Config

type Config struct {
	Orm Orm `yaml:"orm" mapstructure:"orm"`
}

type Orm struct {
	Mysql map[string]Mysql `yaml:"mysql" mapstructure:"mysql"`
}

type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func GetConfig() *Config {
	if config != nil {
		return config
	}
	config = &Config{}
	autoconfigure.ConfigRead(config)
	return config
}
