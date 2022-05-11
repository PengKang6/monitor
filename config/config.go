package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var configs = Configs{}

type Configs struct {
	InfluxDBConfig
	JWTConfig
	ServerConfig
}

type InfluxDBConfig struct {
	// Store the URL of your InfluxDB instance
	Url    string
	Bucket string
	Org    string
	// 此token是基于monitor字符生成的，在8086面板找到的，踩坑。。
	Token string
}

type JWTConfig struct {
	JWTSecret  string
	ExpireTime int
	Issuer     string
}

type ServerConfig struct {
	Port                string
	UpstreamInstances   []string
	DownstreamInstances []string
}

func LoadConfig(confDir string) *Configs {
	viper.AddConfigPath(confDir)
	viper.SetConfigName("monitor")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("unable to read config: %s", err))
	}
	//映射到config配置中
	if err := viper.Unmarshal(&configs); err != nil {
		panic(fmt.Errorf("unable to decode into struct：  %s", err))
	}
	return &configs
}

func GetConfig() *Configs {
	return &configs
}
