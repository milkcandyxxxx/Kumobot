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
// func (c *Ctx) GetUserInfo(id string) (core.UserInfo, error) {
// 	return c.bot.GetUserInfo()
// }

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
func (b *Bot) SendPrivateMessage(userID string, msg string) error {
	return b.adapter.SendPrivateMessage(userID, msg)
}

// SendGroupMessageAt 发送群里信息
func (b *Bot) SendGroupMessageAt(atUserID string, groupID string, msg string) error {
	return b.adapter.SendGroupMessage(atUserID, groupID, msg)
}
func (b *Bot) SendGroupMessage(groupID string, msg string) error {
	return b.adapter.SendGroupMessage("", groupID, msg)
}

// Send 一键发送默认为回复（在哪触发的哪里回复）
func (b *Bot) Send(msg string) error {
	if b.Event.DetailType == "private" {
		return b.SendPrivateMessage(b.Event.UserID, msg)
	}
	if b.Event.DetailType == "channel" || b.Event.DetailType == "group" {
		return b.SendGroupMessageAt(b.Event.UserID, b.Event.GroupID, msg)
	}
	return nil
}
func (b *Bot) SendAt(atUserID string, msg string) error {
	if b.Event.DetailType == "group" {
		return b.SendGroupMessageAt(atUserID, b.Event.GroupID, msg)
	}
	return nil
}

// ++++++++++++++++++++++++信息获取++++++++++++++++++++++++

// GetUserInfo 获取用户信息
func (b *Bot) GetUserInfo(userID string) (adapter.UserInfo, error) {
	return b.adapter.GetUserInfo(userID)
}

// GetSelfInfo 获取自身数据
func (b *Bot) GetSelfInfo() (adapter.SelfInfRes, error) {
	return b.adapter.GetSelfInfo()
}
