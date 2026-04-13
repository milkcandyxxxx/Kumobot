/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package plugins

import "github.com/milkcandyxxxx/Kumobot/plugin"

func init() {
	// 指定词语回复
	plugin.OnCommand("/echo", func(ctx *plugin.Ctx) {
		ctx.Send(ctx.ExtractPlainText())
	})
}
