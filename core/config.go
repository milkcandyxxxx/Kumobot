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

// SelfInfRes 机器人自身信息结构体
type SelfInfRes struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    struct {
		UserId          string `json:"user_id"`
		UserName        string `json:"user_name"`
		Platform        string `json:"platform"`
		UserDisplayname string `json:"user_displayname"`
	} `json:"data"`
	Message string `json:"message"`
}
type UserInfo struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    struct {
		UserId   string `json:"user_id"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int    `json:"age"`
		Area     string `json:"area"`
	} `json:"data"`
	Message string `json:"message"`
}

// GlobalConfig 全局配置
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
