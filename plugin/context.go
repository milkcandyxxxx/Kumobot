/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 上下文配置，自动选择平台群聊等
 */

package plugin

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/core"
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
