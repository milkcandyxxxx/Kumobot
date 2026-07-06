/**
 * @author milkcandy
 * @date 2026/4/14
 * @description TODO
 */

package bot

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/bot/ON11"
	"github.com/robfig/cron/v3"
)

// Bot bot适配器
type Bot struct {
	Config         *adapter.Config // 机器人配置
	Adapter        adapter.Adapter // 适配器选择
	Prefix         string
	Module         []func(event *adapter.Event)
	Event          *adapter.Event
	cronScheduler  *cron.Cron          // 定时任务
	MessageChannel chan *adapter.Event // 读取消息的通道
	ChatHistory    map[string]chan *adapter.Event
}
type Rule func(ctx *Ctx) bool

// GetUserInfo
// func (b *Bot) GetUserInfo() (core.UserInfo, error) {
// 	return b.Adapter.GetUserInfo()
// }

// NewBot 新建bot
func NewBot(config *adapter.Config, prefix string) *Bot {
	return &Bot{
		Config:         config,                               // 配置文件
		Prefix:         prefix,                               // 已废弃，原用于全局触发提示词
		ChatHistory:    make(map[string]chan *adapter.Event), // 初始化
		MessageChannel: make(chan *adapter.Event, 100),       // 消息通道，所有消息都先放进该通道内
	}

}

// OnEvent 注册模块，目前仅有插件模块，后续可能会考虑弃用
func (b *Bot) OnEvent(module func(event *adapter.Event)) {
	b.Module = append(b.Module, module)
}

// Execute 启动bot
func (b *Bot) Execute() {

	for {
		event := <-b.MessageChannel  // 循环接取信息
		if event.Type == "message" { // 是消息类型则传入模块（目前仅有插件模块）
			for _, h := range b.Module {
				h(event)
			}
		}
	}
}

// Runbot 启动bot
func (b *Bot) Runbot() {
	err := b.Adapter.Connect()
	if err != nil {
		return
	}
	go b.Adapter.ReadMessage(b.MessageChannel) // 两个协程同时进行消息读取于接收
	go b.Execute()
	b.StartCron() // 注册定时任务
}

// SetAdapter 设置适配器
func (b *Bot) SetAdapter(a adapter.Adapter) {
	b.Adapter = a
}

// Webhook 暂时用于消息推送的解决方案
var Webhook *Bot

func SetWebhook(b *Bot) {
	Webhook = b
}
func (b *Bot) StartCron() {
	b.cronScheduler = cron.New(cron.WithSeconds())
	mu.RLock()
	for _, p := range plugins {
		for _, m := range p.Matcher {
			if m.Type != "cron" {
				continue
			}
			ctx := &Ctx{}
			switch a := b.Adapter.(type) {
			case *ON11.ON11:
				ctx = &Ctx{Bot: b, ON11: a}
			}
			b.cronScheduler.AddFunc(m.Pattern, func() {
				m.Handler(ctx)
			})

		}
	}
	mu.RUnlock()
	b.cronScheduler.Start()
}
