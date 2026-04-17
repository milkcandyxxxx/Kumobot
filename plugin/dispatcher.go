/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 调度器，判断所有匹配
 */

package plugin

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/core"
	"sort"
	"strings"
)

// Dispatch 插件模块调度器
func (b *Bot) Dispatch(event *core.Event) {
	if event.Type != "message" {
		return
	}
	ctx := &Ctx{
		event:   event,
		bot:     b,
		message: event.GetMessageText(),
	}
	mu.Lock()
	defer mu.Unlock()
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Priority > plugins[j].Priority
	})
	for i, p := range plugins {
		fmt.Printf("%d,%v", i, p)
	}
	for _, p := range plugins {
		for _, m := range p.Matcher {
			// 默认为匹配
			matched := false
			// 依据类型判断是否匹配
			switch m.Type {
			case "startswith":
				matched = strings.HasPrefix(ctx.message, m.Pattern)
			case "cmd":
				matched = isCmd(ctx.message, m.Pattern)
			}
			// 匹配则执行
			if matched {
				m.Handler(ctx)
			}
			// 判断是否独家（向下传递）
			if p.Exclusive {
				return
			}
		}
	}
}

// isCmd 类型的匹配规则
func isCmd(msg string, cmd string) bool {

	return strings.HasPrefix(msg, cmd)
}
