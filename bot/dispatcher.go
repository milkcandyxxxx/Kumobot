/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 调度器，判断所有匹配
 */

package bot

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/bot/ON11"
	"sort"
	"strings"
)

// CreateChathistoryChannel 创建多轮会话通道
func (b *Bot) CreateChathistoryChannel(ctx *Ctx) {
	// 以用户id+群组id作为字典
	GroupUserID := fmt.Sprintf("%s%s", ctx.Event.UserID, ctx.Event.GroupID)
	b.ChatHistory[GroupUserID] = make(chan *adapter.Event)
	ctx.Ch = b.ChatHistory[GroupUserID]
	ctx.ChatHistory = b.ChatHistory
	fmt.Println("刚赋值的ctx", ctx.ChatHistory)
	fmt.Println(ctx.Ch, "xxxxxxxxxxxxx已经创建")
}

// Dispatch 插件模块调度器
func (b *Bot) Dispatch(event *adapter.Event) {
	fmt.Println("bot中的字典", b.ChatHistory)
	// 尝试取值，看有没有多轮会话在等待
	ChatHistory, ok := b.ChatHistory[fmt.Sprintf("%s%s", event.UserID, event.GroupID)]
	// 将信息发送至指定通道
	if ok {
		ChatHistory <- event
		return
	}
	ctx := &Ctx{}
	// 枚举所有支持的协议类型
	switch a := b.Adapter.(type) {
	case *ON11.ON11:
		a.Event = event
		// 这里ctx中的ON11不能直接使用 a ，
		// a中的event的地址类型，如果写入就会导致新消息到来时，
		// 所有ctx的event都会变化，所以这里写入值而非地址
		ctx = &Ctx{
			Bot: b, ON11: &ON11.ON11{
				Adapter: a.Adapter,
				Event:   event,
			}, Event: event,
		}
	}
	mu.Lock()
	defer mu.Unlock()
	// 按照优先级排序
	sort.Slice(b.plugins, func(i, j int) bool {
		return b.plugins[i].Priority > b.plugins[j].Priority
	})
	// 遍历插件
	for _, p := range b.plugins {
		// 遍历插件规则
		for _, m := range p.Respond {
			if checkMatcher(ctx, m) {
				if m.LifeCycle {
					b.CreateChathistoryChannel(ctx)
					go m.HandlerFlow(ctx)
				} else {
					go m.HandlerFlow(ctx)
				}
			}
		}
		if p.Exclusive {
			return
		}
	}
}
func checkMatcher(ctx *Ctx, matcher *Responder) bool {
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
