/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package plugins

import (
	"github.com/milkcandyxxxx/Kumobot/bot"
	"log"
)

func init() {
	bot.OnPlugin("ping", "", "1", "false")
	log.Println("注册成功")
	// 指定词语回复
	bot.OnCommand("1", func(b *bot.Bot) {
		b.Send("pong")
	})
}
