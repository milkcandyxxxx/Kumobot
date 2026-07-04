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
	"math/rand"
	"strconv"
	"time"
)

func init() {
	bot.OnPlugin("ping", "无", "milk", "无", "2", "false")
	bot.OnMessage(func(ctx *bot.Ctx) {
		ctx.ON11.Send(fmt.Sprintf("%+v", 111))

	}, rule.OnCommand("ping"))
	bot.OnChat(func(ctx *bot.Ctx) {
		rand.Seed(time.Now().UnixNano())
		a := rand.Intn(10)
		for {
			ctx.ON11.Send(fmt.Sprintf("%s", "输入数字"))
			b := <-ctx.Ch
			num, _ := strconv.Atoi(b.AltMessage)
			if a > num {
				ctx.ON11.Send("太小了")
			} else if a < num {
				ctx.ON11.Send("太大了")
			} else if a == num {
				ctx.ON11.Send("猜对了")
				return
			}
		}
	}, rule.OnCommand("aaa"))
}
