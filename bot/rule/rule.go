/**
 * @author milkcandy
 * @date 2026/6/29
 * @description TODO
 */

package rule

import "github.com/milkcandyxxxx/Kumobot/bot"

// Rule 单个规则的函数

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
