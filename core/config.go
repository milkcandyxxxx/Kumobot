package core

import (
	"github.com/spf13/viper"
)

// Config 配置文件总结构体
type Config struct {
	Onebots OnebotsConfig `mapstructure:"onebots"`
	Bot     BotConfig     `mapstructure:"bot"`
}

// OnebotsConfig http及webst地址结构体
type OnebotsConfig struct {
	WsURL   string `mapstructure:"ws_url"`
	HttpURL string `mapstructure:"http_url"`
}

// BotConfig 机器人配置结构体
type BotConfig struct {
	Name   string   `mapstructure:"name"`
	Prefix string   `mapstructure:"prefix"` // 命令前缀，如 "/"
	Admins []string `mapstructure:"admins"` // 管理员列表
}

// 全局配置
var GlobalConfig Config

// LoadConfig 读取配置文件并写入全局配置
func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&GlobalConfig)
}
