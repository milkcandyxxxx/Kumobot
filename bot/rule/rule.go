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

// OnlyGroup 指定只能群组触发
func OnlyGroup() bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return ctx.Event.DetailType == "group"
	}
}

// Group 指定什么群触发
func Group(groupID string) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return ctx.Event.GroupID == groupID
	}
}

// User 指定用户触发
func User(userID string) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return ctx.Event.UserID == userID
	}
}

// OnCommand 指定前置词触发
func OnCommand(cmd string) bot.Rule {
	return func(ctx *bot.Ctx) bool {
		return strings.HasPrefix(ctx.Event.AltMessage, cmd)
	}
}

// // OnAdmin 群组管理员触发.,目前不太好触发，因为公共字段没有进行判断的
// func OnAdmin() bot.Rule {
// 	return func(ctx *bot.Ctx) bool {
// 		return ctx.Event
// 	}
// }
