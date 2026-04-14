/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package plugins

import (
	"github.com/milkcandyxxxx/Kumobot/plugin"
	"log"
)

func init() {
	plugin.OnPlugin("ping", "0.0.1", "milkcandy", "pong")
	log.Println(plugin.GetPluginName())
	// 指定词语回复
	plugin.OnCommand("ping", func(ctx *plugin.Ctx) {
		ctx.Send("pong")
	})
}
