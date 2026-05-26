/**
 * @author milkcandy
 * @date 2026/5/25
 * @description TODO
 */

package plugins

import (
	"github.com/milkcandyxxxx/Kumobot/bot"
	"time"
)

func init() {
	bot.OnPlugin("time", "无", "milk", "无", "0", "false")
	// 指定词语回复
	bot.OnCommand("time", func(b *bot.Bot) {
		a := time.Now().Format(time.RFC3339)
		b.Send(a)
	})
}
