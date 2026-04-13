/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 调度器，判断所有匹配
 */

package plugin

import (
	"github.com/milkcandyxxxx/Kumobot/core"
	"strings"
)

// Dispatch 调度主函数实现
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
	for _, m := range matchers {
		matched := false
		switch m.Type {
		case "startswith":
			matched = strings.HasPrefix(ctx.message, m.Pattern)
		case "cmd":
			matched = isCmd(ctx.message, m.Pattern)

		}

		if matched {
			m.Handler(ctx)
		}
	}
}

// cmd 类型的匹配规则
func isCmd(msg string, cmd string) bool {
	return strings.HasPrefix(msg, cmd)
}
