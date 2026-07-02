/*
*
  - @author milkcandy
  - @date 2026/4/13
  - @description TODO
*/
package plugins

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/bot"
	"github.com/milkcandyxxxx/Kumobot/bot/rule"
)

func init() {
	bot.OnPlugin("ping", "无", "milk", "无", "2", "false")
	bot.OnMessage(func(ctx *bot.Ctx) {
		ctx.ON11.Send(fmt.Sprintf("%+v", 111))

	}, rule.Group("602297234"), rule.OnCommand("ping"))
}
