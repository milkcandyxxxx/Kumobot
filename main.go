// main.go
package main

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/core"
	"github.com/milkcandyxxxx/Kumobot/plugin"

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
	b := plugin.NewBot(&core.GlobalConfig)
	// 设置适配器
	adp := adapter.NewOneBotAdapter(
		core.GlobalConfig.Onebots.WsURL,
		core.GlobalConfig.Onebots.HttpURL,
	)
	// bot绑定适配器
	b.SetAdapter(adp)
	// 注册插件
	adp.OnEvent(func(event *core.Event) {
		if event.Type != "message" {
			return
		}
		b.Dispatch(event)
	})

	if err := adp.Connect(); err != nil {
		log.Fatal("连接失败:", err)
	}
	log.Println("Kumobot 已启动")

	select {}
}
