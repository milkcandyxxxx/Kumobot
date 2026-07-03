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

	}, rule.OnCommand("ping"))
	bot.OnChat(func(ctx *bot.Ctx) {
		ctx.ON11.Send(fmt.Sprintf("%+v", 111))
		fmt.Printf("通道是%+v", ctx.Ch)
		for {
			aaa := <-ctx.Ch
			fmt.Println("通道内得到的", aaa)
			ctx.ON11.Send(aaa.AltMessage)
		}

	}, rule.OnCommand("aaa"))
}
