/**
 * @author milkcandy
 * @date 2026/4/15
 * @description TODO
 */

package plugin

import (
	"strconv"
	"strings"
)

func (c *Ctx) Send(msg string) error {
	if c.event.DetailType == "private" {
		return c.bot.SendPrivateMessage(c.event.UserID, msg)
	}
	if c.event.DetailType == "channel" {
		return c.bot.SendGroupMessage(c.event.GroupID, msg)
	}
	return nil
}

// ExtractPlainText 用于获取第一个参数
func (c *Ctx) ExtractPlainText() string {
	return strings.SplitN(c.message, " ", 2)[1]
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

// OnPlugin 注册插件
func OnPlugin(info ...string) {
	thisPlugin := &Plugin{
		Name:      "无",
		Version:   "无",
		Author:    "佚名",
		Help:      "无",
		Priority:  0,
		Exclusive: false,
		Matcher:   nil,
	}
	thisPlugin.Name = info[0]
	thisPlugin.Help = info[1]
	priority, _ := strconv.Atoi(info[2])
	exclusive, _ := strconv.ParseBool(info[3])
	thisPlugin.Priority = priority
	thisPlugin.Exclusive = exclusive

	addPlugin(thisPlugin)
}
