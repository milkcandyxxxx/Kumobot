/**
 * @author milkcandy
 * @date 2026/4/15
 * 对于bot用法的二次包装
 * @description TODO
 */

package bot

import (
	"strconv"
	"strings"
)

func (b *Bot) Send(msg string) error {
	if b.event.DetailType == "private" {
		return b.SendPrivateMessage(b.event.UserID, msg)
	}
	if b.event.DetailType == "channel" {
		return b.SendGroupMessage(b.event.GroupID, msg)
	}
	return nil
}

// GetUserInfo 获取用户信息
// func (c *Ctx) GetUserInfo(id string) (core.UserInfo, error) {
// 	return c.bot.GetUserInfo()
// }

// ExtractPlainText 用于获取第一个参数
func (b *Bot) ExtractPlainText() string {
	return strings.SplitN(b.event.GetMessageText(), " ", 2)[1]
}

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
	thisPlugin.Help = info[1]
	priority, _ := strconv.Atoi(info[2])
	exclusive, _ := strconv.ParseBool(info[3])
	thisPlugin.Priority = priority
	thisPlugin.Exclusive = exclusive
	addPlugin(thisPlugin)
}
