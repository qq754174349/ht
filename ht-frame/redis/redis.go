package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/logger"
	"sync"
	"time"
)

var (
	redisInstances sync.Map
	defaultName    string
	config         *Redis
)

type Redis struct {
	Redis map[string]Config `yaml:"redis" mapstructure:"redis"`
}

type Config struct {
	Addr     string
	User     string
	Password string
	DB       int
}

type AutoConfig struct{}

func init() {
	err := autoconfigure.Register(AutoConfig{})
	if err != nil {
		logger.Fatal("redis 自动配置注册失败")
	}
}

func (AutoConfig) Init() error {
	config = &Redis{}
	autoconfigure.ConfigRead(config)
	first := true
	for name, redisCfg := range config.Redis {
		if err := initRedis(redisCfg, name); err != nil {
			return fmt.Errorf("redis[%s]初始化失败: %v", name, err)
		}
		if first {
			defaultName = name
			first = false
		}
	}

	if defaultName == "" {
		return fmt.Errorf("至少需要配置一个Redis数据源")
	}
	return nil
}

func initRedis(cfg Config, name string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 健康检查
	if err := ping(client); err != nil {
		return fmt.Errorf("redis[%s]连接失败: %v", name, err)
	}

	redisInstances.Store(name, client)
	go monitor(name, client, cfg)
	logger.Infof("Redis[%s]已初始化 @ %s", name, cfg.Addr)
	return nil
}

func Get(name ...string) (*redis.Client, error) {
	var instancesName string
	if len(name) == 0 {
		instancesName = defaultName
	} else {
		instancesName = name[0]
	}
	val, ok := redisInstances.Load(instancesName)
	if !ok {
		return nil, fmt.Errorf("redis[%s]未初始化", name)
	}

	client := val.(*redis.Client)
	if err := verify(client); err != nil {
		return nil, fmt.Errorf("redis[%s]连接异常: %v", name, err)
	}
	return client, nil
}

// 内部工具函数
func ping(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}

func monitor(name string, client *redis.Client, cfg Config) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := verify(client); err != nil {
			logger.Warnf("Redis[%s]连接异常: %v", name, err)
			_ = reconnect(name, cfg)
		}
	}
}

func reconnect(name string, cfg Config) error {
	if val, ok := redisInstances.Load(name); ok {
		if client, ok := val.(*redis.Client); ok {
			_ = client.Close()
		}
		redisInstances.Delete(name)
	}
	return initRedis(cfg, name)
}

// verify 检查Redis连接是否健康
func verify(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 1. 基础Ping测试
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping失败: %v", err)
	}

	// 2. 可选：执行实际命令测试
	if err := client.Do(ctx, "SELECT", client.Options().DB).Err(); err != nil {
		return fmt.Errorf("DB选择测试失败: %v", err)
	}

	// 3. 检查连接池状态
	poolStats := client.PoolStats()
	if poolStats.TotalConns == 0 {
		return fmt.Errorf("连接池异常: 无可用连接")
	}

	logger.Debugf("Redis连接正常 (活跃连接: %d)", poolStats.IdleConns)
	return nil
}

func CloseAllRedis() {
	redisInstances.Range(func(key, value interface{}) bool {
		if client, ok := value.(*redis.Client); ok {
			_ = client.Close()
			logger.Infof("Redis[%v]已关闭", key)
		}
		return true
	})
}
