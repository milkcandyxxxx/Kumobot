/**
 * @author milkcandy
 * @date 2026/4/14
 * @description TODO
 */

package bot

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"strings"
)

// Bot bot适配器
type Bot struct {
	Config  *adapter.Config // 机器人配置
	adapter adapter.Adapter // 适配器选择
	Prefix  string
	Module  []func(event *adapter.Event)
	Event   *adapter.Event
}

// GetUserInfo
// func (b *Bot) GetUserInfo() (core.UserInfo, error) {
// 	return b.adapter.GetUserInfo()
// }

// NewBot 新建bot
func NewBot(config *adapter.Config, prefix string) *Bot {
	return &Bot{
		Config: config,
		Prefix: prefix,
	}
}

// OnEvent 注册事件监听函数
func (b *Bot) OnEvent(module func(event *adapter.Event)) {
	b.Module = append(b.Module, module)
}

func (b *Bot) Execute() {
	err := b.adapter.Connect()
	if err != nil {
		return
	}

	for {
		msg, _ := b.adapter.ReadMessage()
		event, ok := msg.(*adapter.Event)

		if !ok {
			continue
		}
		if event.Type == "message" {

			if !strings.HasPrefix(event.AltMessage, b.Prefix) {
				continue
			}
			event.AltMessage = event.AltMessage[len(b.Prefix):]
			for _, h := range b.Module {
				h(event)
			}
		}
	}
}
func (b *Bot) Runbot() {
	go b.Execute()
}

// SetAdapter 设置适配器
func (b *Bot) SetAdapter(a adapter.Adapter) {
	b.adapter = a
}

// func (b *Bot) DeleteMessage(messageId string) error {
// 	return b.adapter.DeleteMessage(messageId)
// }

// func (b *Bot) CallAction(action string, params map[string]string) (adapter.Response, error) {
// 	return b.adapter.(*onebot11.Adapter).CallAction(action, params)
// }

// 暂时用于消息推送的解决方案
var Webhook *Bot

func SetWebhook(b *Bot) {
	Webhook = b
}

// func SendGroupMessage(groupID string, msg string) error {
// 	return Webhook.adapter.SendGroupMessage("1910399844", groupID, msg)
// }

// func SendPrivateMessage(userID string, msg string) error {
// 	return Webhook.adapter.SendPrivateMessage(userID, msg)
// }
