package config

import (
	"time"

	"github.com/qq754174349/ht/ht-frame/autoconfigure"
)

var config *Web

type Web struct {
	Web Config
}

type Config struct {
	Port    string
	Timeout time.Duration
}

func Get() *Web {
	if config == nil {
		config = &Web{}
		autoconfigure.ConfigRead(config)
	}
	return config
}
