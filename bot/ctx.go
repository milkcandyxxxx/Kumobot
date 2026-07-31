/**
 * @author milkcandy
 * @date 2026/4/15
 * 对于插件写法的二次包装
 * @description TODO
 */

package bot

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/bot/ON11"
)

// GetUserInfo 获取用户信息
//
//	func (c *Ctx) GetUserInfo(id string) (core.UserInfo, error) {
//		return c.bot.GetUserInfo()
//	}

type Ctx struct {
	Bot         *Bot
	ON11        *ON11.ON11
	Event       *adapter.Event
	Ch          chan *adapter.Event
	ChatHistory map[string]chan *adapter.Event
	Plugin      *PluginInfo
}

// ++++++++++++++++++++++++信息发送++++++++++++++++++++++++

// // SendPrivateMessage 发送私聊信息
// func (c *Ctx) SendPrivateMessage(userID string, msg string) (Adapter.CommonResponse, error) {
//
// 	return c.ON11.SendPrivateMessage(userID, msg)
// }

// func (c *Ctx) SendGroupMessage(groupID string, msg string) error {
// 	return c.ON11.SendGroupMessage("", groupID, msg)
// }
// func SendGroupMessageAt() {
// 	return
// }

// func (c *Ctx) SendAt(atUserID string, msg string) error {
// 	if c.Event.DetailType == "group" {
// 		return c.SendGroupMessageAt(atUserID, c.Event.GroupID, msg)
// 	}
// 	return nil
// }

// ++++++++++++++++++++++++信息获取++++++++++++++++++++++++

// // GetUserInfo 获取用户信息
// func (c *Ctx) GetUserInfo(userID string) (Adapter.UserInfo, error) {
// 	return c.Bot.Adapter.GetUserInfo(userID)
// }
// // 修改函数签名，返回具名结构体
// func (c *Ctx) GetLoginInfo() (LoginInfo, error) {
// 	a, _ := c.Bot.Adapter.CallAction("get_login_info", nil)
// 	var data LoginInfo
// 	err := json.Unmarshal(a.Data, &data)
// 	if err != nil {
// 		return LoginInfo{}, err
// 	}
// 	return data, nil
// }
