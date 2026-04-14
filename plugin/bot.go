/**
 * @author milkcandy
 * @date 2026/4/14
 * @description TODO
 */

package plugin

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/core"
)

// Bot bot适配器
type Bot struct {
	config  *core.Config    // 机器人配置
	adapter adapter.Adapter // 适配器选择
}

func (b *Bot) SendPrivateMessage(userID string, msg string) error {
	return b.adapter.SendPrivateMessage(userID, msg)
}
func (b *Bot) SendGroupMessage(groupID string, msg string) error {
	fmt.Println(3)
	return b.adapter.SendGroupMessage(groupID, msg)
}

// NewBot 新建bot
func NewBot(config *core.Config) *Bot {
	return &Bot{
		config: config,
	}
}
