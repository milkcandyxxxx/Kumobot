/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 上下文配置，自动选择平台群聊等
 */

package plugin

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/core"
	"strings"
)

// TODO 目前还未直线自动配置（适配器等）

// SetAdapter 设置适配器
func (b *Bot) SetAdapter(adapter adapter.Adapter) {

	b.adapter = adapter

}

// Ctx 消息上下文
type Ctx struct {
	event   *core.Event // 原始事件
	message string      // 消息纯文本
	bot     *Bot        // 机器人实例
}

func (c *Ctx) Send(msg string) error {
	if c.event.DetailType == "private" {
		return c.bot.SendPrivateMessage(c.event.UserID, msg)
	}
	if c.event.DetailType == "channel" {
		return c.bot.SendGroupMessage(c.event.GroupID, msg)
	}
	return nil
}

// ExtractPlainText 用于获取第一个参数
func (c *Ctx) ExtractPlainText() string {
	return strings.SplitN(c.message, " ", 2)[1]
}
