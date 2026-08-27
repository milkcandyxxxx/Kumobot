/**
 * @author milkcandy
 * @date 2026/4/15
 * 对于插件写法的二次包装
 * @description TODO
 */

package bot

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/adapter/onebot11"
)

// GetUserInfo 获取用户信息
//
//	func (c *Ctx) GetUserInfo(id string) (core.UserInfo, error) {
//		return c.bot.GetUserInfo()
//	}

type Ctx struct {
	Bot         *Bot
	Event       *adapter.Event
	Ch          chan *adapter.Event
	ChatHistory map[string]chan *adapter.Event
}

//	func (c *Ctx) ON11() *ON11.ON11 {
//		return &ON11.ON11{
//			Adapter: c.Bot.Adapter.(*ON11.ON11).Adapter,
//			Event:   c.Event,
//		}
//	}
func (c *Ctx) send() {
	if c.Bot.AdapterName == "onebot11" {
		c.Bot.Adapter.(*onebot11.Adapter).SendPrivateMessage(73808768244, "aaa")
	}

}
