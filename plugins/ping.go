/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package plugins

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/plugin"
)

func init() {
	plugin.OnPlugin("ping", "", "1", "false")
	// 指定词语回复
	plugin.OnRegex("124\\d{4}111", func(ctx *plugin.Ctx) {
		ctx.Send("pong")
		res, _ := ctx.GetSelfInfo()
		fmt.Println(res)
	})
}
