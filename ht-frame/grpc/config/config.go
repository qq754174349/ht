package config

import "github.com/qq754174349/ht/ht-frame/autoconfigure"

var config *GRpc

type GRpc struct {
	Grpc Config `json:"grpc" yaml:"grpc"`
}

type Config struct {
	Port int `json:"port" yaml:"port"`
}

func Get() *GRpc {
	if config == nil {
		config = &GRpc{}
		autoconfigure.ConfigRead(config)
	}
	return config
}
