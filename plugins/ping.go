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
)

func init() {
	bot.OnPlugin("ping", "无", "milk", "无", "2", "false")
	bot.OnCommand("ping", func(ctx *bot.Ctx) {
		info, _ := ctx.ON11.GetLoginInfo()
		ctx.ON11.Send(fmt.Sprintf("%+v", info))
	})

}
