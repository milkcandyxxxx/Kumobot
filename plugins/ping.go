/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package plugins

import (
	"github.com/milkcandyxxxx/Kumobot/bot"
)

func init() {
	bot.OnPlugin("ping", "无", "milk", "无", "2", "true")
	// 指定词语回复
	bot.OnCommand("1", func(b *bot.Bot) {
		b.Send("pong")
	})
}
