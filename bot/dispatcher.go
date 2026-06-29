/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 调度器，判断所有匹配
 */

package bot

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/bot/ON11"
	"regexp"
	"sort"
	"strings"
)

// Dispatch 插件模块调度器
func (b *Bot) Dispatch(event *adapter.Event) {

	if event.Type != "message" {
		return
	}
	ctx := &Ctx{}
	switch a := b.Adapter.(type) {
	case *ON11.ON11:
		a.Event = event
		ctx = &Ctx{Bot: b, ON11: a, Event: event}
	}
	mu.Lock()
	defer mu.Unlock()
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Priority > plugins[j].Priority
	})
	for _, p := range plugins {
		for _, m := range p.Matcher {
			if checkMatcher(ctx, m) {
				go m.Handler(ctx)
			}
			// 默认为匹配
			// matched := false
			// // 依据类型判断是否匹配
			// switch m.Type {
			// // 以什么什么开头
			// case "startswith":
			// 	matched = strings.HasPrefix(m.Pattern, event.GetMessageText())
			// // 指令
			// case "cmd":
			// 	matched = isCmd(m.Pattern, event.GetMessageText())
			// // 正则
			// case "regex":
			// 	matched = isRegex(m.Pattern, event.GetMessageText())
			// // 以什么什么结尾
			// case "endswith":
			// 	matched = strings.HasSuffix(m.Pattern, event.GetMessageText())
			// }
			// // 匹配则执行
			// if matched {
			// 	go m.Handler(ctx)
			// }
			// // 判断是否独家（向下传递）
			// if p.Exclusive {
			// 	return
			// }
		}
	}
}
func checkMatcher(ctx *Ctx, matcher Matcher) bool {
	for _, rule := range matcher.Rules {
		if !rule(ctx) {
			return false
		}
	}
	return true
}

// isCmd 类型的匹配规则
func isCmd(cmd string, msg string) bool {

	return strings.HasPrefix(msg, cmd)
}
func isRegex(regex string, msg string) bool {
	match, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return match.MatchString(msg)
}
