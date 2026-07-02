/**
 * @author milkcandy
 * @date 2026/6/29
 * @description TODO
 */

package rule

import (
	"github.com/milkcandyxxxx/Kumobot/bot"
	"strings"
)

// 规则总体的判断器，例如全部满足，部分满足等

// AllRule 全部规则都要满足
func AllRule(rules ...bot.Rule) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		for _, rule := range rules {
			if !rule(ctx) {
				return false
			}
		}
		return true
	}
}

// AnyRule 满足任意一个
func AnyRule(rules ...bot.Rule) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		for _, rule := range rules {
			if rule(ctx) {
				return true
			}
		}
		return false
	}
}
func OnlyGroup() bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return ctx.Event.DetailType == "group"
	}
}
func Group(groupID string) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return ctx.Event.GroupID == groupID
	}
}
func User(userID string) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return ctx.Event.UserID == userID
	}
}
func OnCommand(cmd string) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return strings.HasPrefix(ctx.Event.AltMessage, cmd)
	}
}
