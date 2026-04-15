/**
 * @author milkcandy
 * @date 2026/4/15
 * @description TODO
 */

package plugin

import "strings"

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
		Type:      "cmd",
		Pattern:   cmd,
		Priority:  0,
		Exclusive: false,
		Handler:   h,
	})
}
