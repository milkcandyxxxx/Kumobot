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
	"strconv"
	"strings"
)

// GetUserInfo 获取用户信息
//
//	func (c *Ctx) GetUserInfo(id string) (core.UserInfo, error) {
//		return c.bot.GetUserInfo()
//	}
type Ctx struct {
	Bot  *Bot
	ON11 *ON11
}
type ON11 struct {
	*onebot11.Adapter
	Event *adapter.Event
}

// ExtractPlainText 用于获取第一个参数
func (b *Bot) ExtractPlainText() string {
	return strings.SplitN(b.Event.GetMessageText(), " ", 2)[1]
}

// ++++++++++++++++++++++++匹配机制++++++++++++++++++++++++

// OnCommand 单匹配
func OnCommand(cmd string, h Handler) {
	addMatcher(Matcher{
		Type:    "cmd",
		Pattern: cmd,
		// Priority:  0,
		// Exclusive: false,
		Handler: h,
	})
}

// OnRegex 正则匹配
func OnRegex(regex string, h Handler) {
	addMatcher(Matcher{
		Type:    "regex",
		Pattern: regex,
		Handler: h,
	})
}

// ++++++++++++++++++++++++插件管理++++++++++++++++++++++++

// OnPlugin 注册插件
func OnPlugin(info ...string) {
	thisPlugin := &Plugin{
		Name:      "无",
		Version:   "无",
		Author:    "佚名",
		Help:      "无",
		Priority:  0,
		Exclusive: false,
	}
	thisPlugin.Name = info[0]
	thisPlugin.Help = info[3]
	thisPlugin.Version = info[1]
	thisPlugin.Author = info[2]
	priority, _ := strconv.Atoi(info[4])
	exclusive, _ := strconv.ParseBool(info[5])
	thisPlugin.Priority = priority
	thisPlugin.Exclusive = exclusive
	addPlugin(thisPlugin)
}

// ++++++++++++++++++++++++信息发送++++++++++++++++++++++++

// SendPrivateMessage 发送私聊信息
func (c *Ctx) SendPrivateMessage(userID string, msg string) error {

	return c.ON11.SendPrivateMessage(userID, msg)
}

// SendGroupMessageAt 发送群里信息
func (o *ON11) SendGroupMessageAt(atUserID string, groupID string, msg string) error {
	return o.SendGroupMessage(atUserID, groupID, msg)
}

// func (c *Ctx) SendGroupMessage(groupID string, msg string) error {
// 	return c.ON11.SendGroupMessage("", groupID, msg)
// }
// func SendGroupMessageAt() {
// 	return
// }

// Send 一键发送默认为回复（在哪触发的哪里回复）
func (o *ON11) Send(msg string) error {
	if o.Event.DetailType == "private" {
		return o.SendPrivateMessage(o.Event.UserID, msg)
	}
	if o.Event.DetailType == "channel" || o.Event.DetailType == "group" {
		return o.SendGroupMessageAt(o.Event.UserID, o.Event.GroupID, msg)
	}
	return nil
}

// func (c *Ctx) SendAt(atUserID string, msg string) error {
// 	if c.Event.DetailType == "group" {
// 		return c.SendGroupMessageAt(atUserID, c.Event.GroupID, msg)
// 	}
// 	return nil
// }

// ++++++++++++++++++++++++信息获取++++++++++++++++++++++++

// // GetUserInfo 获取用户信息
// func (c *Ctx) GetUserInfo(userID string) (adapter.UserInfo, error) {
// 	return c.Bot.adapter.GetUserInfo(userID)
// }

// LoginInfo 自身信息
type LoginInfo struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

// // 修改函数签名，返回具名结构体
// func (c *Ctx) GetLoginInfo() (LoginInfo, error) {
// 	a, _ := c.Bot.adapter.CallAction("get_login_info", nil)
// 	var data LoginInfo
// 	err := json.Unmarshal(a.Data, &data)
// 	if err != nil {
// 		return LoginInfo{}, err
// 	}
// 	return data, nil
// }
