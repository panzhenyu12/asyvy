package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// Redis
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	Web struct {
		Addr string `yaml:"addr"`
	} `yaml:"web"`
	Worker struct {
		Concurrency int    `yaml:"concurrency"`
		CacheDir    string `yaml:"cache_dir"`
	} `yaml:"worker"`
}

func GetConfig(confpath string) (*Config, error) {
	var conf *Config
	var err error
	if len(confpath) > 0 {
		conf, err = LoadFromFile(confpath)
		if err != nil {
			return nil, err
		}
	} else {
		conf, err = LoadFromEnv()
		if err != nil {
			return nil, err
		}
	}
	if conf.Worker.Concurrency == 0 {
		conf.Worker.Concurrency = 10
	}
	return conf, nil
}

func LoadFromFile(path string) (*Config, error) {
	var c Config
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func LoadFromEnv() (*Config, error) {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("ASYVY")
	// 通过结构体标签绑定环境变量
	viper.BindEnv("redis.addr")
	viper.BindEnv("redis.password")
	viper.BindEnv("redis.db")
	viper.BindEnv("web.addr")
	viper.BindEnv("worker.concurrency")
	viper.BindEnv("worker.cache_dir")
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
