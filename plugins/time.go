/**
 * @author milkcandy
 * @date 2026/5/25
 * @description TODO
 */

package plugins

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/bot"
)

func init() {
	bot.OnPlugin("time", "无", "milk", "无", "0", "false")
	// 指定词语回复
	bot.OnCommand("time", func(c *bot.Ctx) {
		d, _ := c.ON11.GetForwardMsg("1047161647")
		fmt.Printf("@@@@@@@@@@%#v\n", d.Messages[1].Message[0].Data)
	})
}
