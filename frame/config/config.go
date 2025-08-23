package config

const (
	Logger = "log"
	WEB    = "web"
	MYSQL  = "mysql"
	REDIS  = "redis"
	CONSUL = "consul"
)

var appCfg *AppConfig

type AppConfig struct {
	Active  string
	AppName string `yaml:"app_name" json:"app_name" mapstructure:"app_name"`
}

func SetAppCfg(appConfig *AppConfig) {
	appCfg = appConfig
}

func GetAppCfg() *AppConfig {
	return appCfg
}
