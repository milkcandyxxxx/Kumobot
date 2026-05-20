// main.go
package main

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/bot"
	"github.com/milkcandyxxxx/Kumobot/core"
	_ "github.com/milkcandyxxxx/Kumobot/plugins"
	"log"
)

func main() {
	log.Println("Kumobot 启动中...")
	// 加载配置
	if err := core.LoadConfig(); err != nil {
		log.Fatal("加载配置失败:", err)
	}
	// 新建bot实例
	b := bot.NewBot(&core.GlobalConfig)
	// 设置适配器
	adp := adapter.NewOneBotAdapter(
		core.GlobalConfig.Onebots.WsURL,
		core.GlobalConfig.Onebots.HttpURL,
		core.GlobalConfig.Bot.Prefix,
	)
	b.SetAdapter(adp)
	// 注册插件
	m := b.NewModule()
	m.add
	b.Runbot()
	log.Println("Kumobot 已启动")

	select {}
}
