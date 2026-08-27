/**
 * @author milkcandy
 * @date 2026/4/14
 * @description TODO
 */

package bot

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/adapter/onebot11"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
)

// Bot bot适配器
type Bot struct {
	Config      *adapter.Config // 机器人配置
	Adapter     adapter.Adapter // 适配器选择
	AdapterName string          // 协议名称
	// Prefix         string                         // 全局命令前缀
	Module         []func(event *adapter.Event)   // 模块列表
	Event          *adapter.Event                 // 事件结构体
	cronScheduler  *cron.Cron                     // 定时任务
	MessageChannel chan *adapter.Event            // 读取消息的通道
	ChatHistory    map[string]chan *adapter.Event // 消息历史
	Plugins        []*PluginInfo                  // 插件列表
	Mu             sync.RWMutex                   // 写锁
}
type Rule func(ctx *Ctx) bool

// GetUserInfo
// func (b *Bot) GetUserInfo() (core.UserInfo, error) {
// 	return b.Adapter.GetUserInfo()
// }

// NewBot 新建bot
func NewBot(config *adapter.Config) *Bot {
	return &Bot{
		Config:         config,                               // 配置文件
		ChatHistory:    make(map[string]chan *adapter.Event), // 初始化
		MessageChannel: make(chan *adapter.Event, 100),       // 消息通道，所有消息都先放进该通道内
	}

}

// NewOneBot11Adapter 新建适配器
func (b *Bot) NewOneBot11Adapter(c adapter.Config) adapter.Adapter {
	info := adapter.InitializeAdapterInformation("onebot11")
	// 手动包一层 onebot11.Adapter
	b.Adapter = &onebot11.Adapter{
		AdapterInfo: info,
	}
	fmt.Printf("%+v", b.Adapter)
	b.AdapterName = "onebot11"
	return &onebot11.Adapter{
		AdapterInfo: info,
	}
}

// OnEvent 注册模块，目前仅有插件模块，后续可能会考虑弃用
func (b *Bot) OnEvent(module func(event *adapter.Event)) {
	b.Module = append(b.Module, module)
}

// Execute 启动接收消息
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
	err := b.Adapter.Connect() // 先连接
	if err != nil {
		return
	}
	go b.Adapter.ReadMessage(b.MessageChannel) // 两个协程同时进行消息读取于接收
	go b.Execute()                             // 接收信息
	// go b.Adapter.CheckHeartbeat(15, 10*time.Second) 心跳检查
	b.StartCron() // 注册定时任务
}

// SetAdapter 设置适配器
func (b *Bot) SetAdapter(a adapter.Adapter) {
	b.Adapter = a
}

// StartCron 启动定时任务
func (b *Bot) StartCron() {
	b.cronScheduler = cron.New(cron.WithSeconds())
	for _, p := range b.Plugins {
		for _, m := range p.Respond {
			if m.Type != "Cron" {
				continue
			}
			ctx := &Ctx{}
			// switch a := b.Adapter.(type) {
			// case *ON11.ON11:
			ctx = &Ctx{Bot: b}
			// }
			b.cronScheduler.AddFunc(m.Timing, func() {
				m.Handler(ctx)
			})

		}
	}
	b.cronScheduler.Start()
}

// OnPlugin 注册插件
func (b *Bot) OnPlugin(p *PluginInfo) *PluginInfo {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	b.Plugins = append(b.Plugins, p)
	log.Println("\n加载插件:", p.Name, "\n版本:", p.Version, "\n作者:", p.Author, "\n帮助:", p.Help, "\n优先级:", p.Priority, "\n独家:", p.Exclusive)
	return p
}

// OnON11Plugin 注册插件
// func (b *Bot) OnON11Plugin(p *ON11.Plugin) *ON11.Plugin {
// 	b.Mu.Lock()
// 	defer b.Mu.Unlock()
// 	b.Plugins = append(b.Plugins, p.PluginInfo)
// 	log.Println("\n加载插件:", p.Name, "\n版本:", p.Version, "\n作者:", p.Author, "\n帮助:", p.Help, "\n优先级:", p.Priority, "\n独家:", p.Exclusive)
// 	return p
// }

// RegisterPlugin bot注册插件
func (b *Bot) RegisterPlugin(p Plugin) {
	p.Register(b)
}

// ON11 用于bot之际获取适配器方法（定时任务等情况下没有ctx上下文使用）
func (b *Bot) ON11() *onebot11.Adapter {
	return b.Adapter.(*onebot11.Adapter)
}
