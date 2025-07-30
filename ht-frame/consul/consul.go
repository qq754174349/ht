package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/qq754174349/ht/ht-frame/autoconfigure"
	log "github.com/qq754174349/ht/ht-frame/logger"
	"time"
)

var (
	config *Consul
	client *api.Client
)

type AutoConfig struct{}

type Consul struct {
	Consul Config
}

type Config struct {
	Addr string
}

func init() {
	err := autoconfigure.Register(AutoConfig{})
	if err != nil {
		log.Fatal("consul 自动配置注册失败")
	}
}

func (AutoConfig) Init() error {
	var err error
	client, err = api.NewClient(&api.Config{Address: GetConfig().Consul.Addr})
	if err != nil {
		return fmt.Errorf("创建 Consul 客户端失败: %v", err)
	}
	return nil
}

func (AutoConfig) Close() error {
	return nil
}

func GetConfig() Consul {
	if config == nil {
		config = &Consul{}
		autoconfigure.ConfigRead(config)
	}
	return *config
}

func Register(reg *api.AgentServiceRegistration) error {
	go monitorConsulAndRegister(reg)
	return nil
}

func Deregister(id string) {
	client.Agent().ServiceDeregister(id)
}

func monitorConsulAndRegister(reg *api.AgentServiceRegistration) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		_, err := client.Agent().Self()
		if err != nil {
			log.Warn("[Consul] 不可用，等待恢复...")
			continue
		}

		services, err := client.Agent().Services()
		if err != nil {
			log.Warn("[Consul] 获取服务列表失败：", err)
			continue
		}

		if _, ok := services[reg.ID]; ok {
			// 已注册
			continue
		}

		// 注册服务
		if err := client.Agent().ServiceRegister(reg); err != nil {
			log.Warnf("[Consul] %s注册失败：%s", reg.Name, err)
		} else {
			log.Infof("[Consul] %s注册成功", reg.Name)
		}
	}
}
