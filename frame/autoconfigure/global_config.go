// Package autoconfigure Package config 全局配置文件
package autoconfigure

import (
	"fmt"
	"log"

	"github.com/qq754174349/ht/ht-frame/config"
	"github.com/spf13/viper"
)

const (
	defaultConfigFileName = "config"
	defaultConfigFileType = "yaml"
)

var (
	appCfg       *config.AppConfig
	initializers []Configuration
)

type Configuration interface {
	Close() error
}

func init() {
	viper.AddConfigPath("configs/")
	viper.SetConfigType(defaultConfigFileType)
	viper.SetConfigName(defaultConfigFileName + "." + defaultConfigFileType)
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	active := viper.GetString("active")
	if active == "nil" {
		log.Fatalf("没有激活的配置")
	}
	log.Printf("Active environment:%s", active)

	viper.SetConfigName(defaultConfigFileName + "-" + active + "." + defaultConfigFileType)
	// 读取配置文件
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	appCfg = &config.AppConfig{}
	err := viper.Unmarshal(appCfg)
	if err != nil {
		log.Fatal("配置文件格式错误")
	}
	config.SetAppCfg(appCfg)
}

func Register(conf ...Configuration) error {
	if len(conf) < 1 {
		return fmt.Errorf("必须注册一个配置初始化器")
	}
	initializers = append(initializers, conf...)
	return nil
}

func Close() {
	for _, v := range initializers {
		err := v.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ConfigRead 配置文件读取，读取统一走这
func ConfigRead(rawVal any) {
	err := viper.Unmarshal(rawVal)
	if err != nil {
		log.Fatal("配置文件格式错误")
	}
}
