package mail

import (
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/logger"
)

var config *Mail

type Mail struct {
	Mail Config `yaml:"mail" json:"mail"`
}

type Config struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

func init() {
	err := autoconfigure.Register(Configuration{})
	if err != nil {
		logger.Fatal("邮件配置初始化失败", err)
	}
}

type Configuration struct {
}

func (Configuration) Init() error {
	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = &Mail{}
		autoconfigure.ConfigRead(config)
	}
	return &config.Mail
}
