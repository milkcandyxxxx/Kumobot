/**
 * @author milkcandy
 * @date 2026/4/15
 * @description TODO
 */

package plugins

import (
	"github.com/milkcandyxxxx/Kumobot/plugin"
	"time"
)

func init() {
	plugin.OnPlugin("time", "milkcandyxxxx", "1", "2")
	plugin.OnCommand("time", func(ctx *plugin.Ctx) {
		now := time.Now()
		ctx.Send(now.Format("2006-01-02 15"))
	})
}
