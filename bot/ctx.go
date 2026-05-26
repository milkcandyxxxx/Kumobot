/**
 * @author milkcandy
 * @date 2026/4/15
 * 对于插件写法的二次包装
 * @description TODO
 */

package bot

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"strconv"
	"strings"
)

// GetUserInfo 获取用户信息
//
//	func (c *Ctx) GetUserInfo(id string) (core.UserInfo, error) {
//		return c.bot.GetUserInfo()
//	}
type Ctx struct {
	Bot   *Bot
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
	return c.Bot.adapter.SendPrivateMessage(userID, msg)
}

// SendGroupMessageAt 发送群里信息
func (c *Ctx) SendGroupMessageAt(atUserID string, groupID string, msg string) error {
	return c.Bot.adapter.SendGroupMessage(atUserID, groupID, msg)
}
func (c *Ctx) SendGroupMessage(groupID string, msg string) error {
	return c.Bot.adapter.SendGroupMessage("", groupID, msg)
}

// Send 一键发送默认为回复（在哪触发的哪里回复）
func (c *Ctx) Send(msg string) error {
	if c.Bot.Event.DetailType == "private" {
		return c.SendPrivateMessage(c.Bot.Event.UserID, msg)
	}
	if c.Bot.Event.DetailType == "channel" || c.Bot.Event.DetailType == "group" {
		return c.SendGroupMessageAt(c.Bot.Event.UserID, c.Event.GroupID, msg)
	}
	return nil
}
func (c *Ctx) SendAt(atUserID string, msg string) error {
	if c.Event.DetailType == "group" {
		return c.SendGroupMessageAt(atUserID, c.Event.GroupID, msg)
	}
	return nil
}

// ++++++++++++++++++++++++信息获取++++++++++++++++++++++++

// GetUserInfo 获取用户信息
func (c *Ctx) GetUserInfo(userID string) (adapter.UserInfo, error) {
	return c.Bot.adapter.GetUserInfo(userID)
}
func (c *Ctx) CallAction(action string, params map[string]string) (adapter.Response, error) {

	return c.Bot.adapter.CallAction(action, nil)
}
