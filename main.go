// main.go
package main

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/adapter/onebot11"

	"github.com/milkcandyxxxx/Kumobot/bot"
	_ "github.com/milkcandyxxxx/Kumobot/plugins"
	"log"
)

func main() {
	log.Println("Kumobot 启动中...")
	// 加载配置
	if err := adapter.LoadConfig(); err != nil {
		log.Fatal("加载配置失败:", err)
	}
	// 新建bot实例
	b := bot.NewBot(&adapter.GlobalConfig, adapter.GlobalConfig.Bot.Prefix)
	// 设置适配器
	adp := onebot11.NewOneBotAdapter(*b.Config)
	b.SetAdapter(&bot.ON11{Adapter: adp})
	bot.SetWebhook(b)
	// 注册插件模块
	b.OnEvent(func(event *adapter.Event) {
		b.Dispatch(event)
	})
	b.Runbot()
	log.Println("Kumobot 已启动")

	select {}
}

