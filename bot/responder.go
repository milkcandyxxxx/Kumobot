/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package bot

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/adapter"
)

// Handler 处理函数（回调函数定义）
type Handler func(ctx *Ctx)

// Responder 事件响应器
type Responder struct {
	Type        string                         // 匹配类型 命令，前缀等等
	Pattern     string                         // 匹配所需的关键词 如 / ! 等
	Rules       []Rule                         // 规则
	ChatHistory chan *adapter.Event            // 存放后续对话内容
	LifeCycle   bool                           // 用于判断是只回复一次还是一场会话
	GroupUserID string                         // 判断 是那条会话
	Handler     Handler                        // 回调函数
	Pre         []func(r *Responder, ctx *Ctx) // 前置函数
	Post        []func(r *Responder, ctx *Ctx) // 后置函数
	Timing      string                         // 定时任务的时间
}

// HandlerFlow 回调函数的执行流程，前置函数和后置函数
func (r *Responder) HandlerFlow(ctx *Ctx) {
	fmt.Println("传入HandlerFlow的ctx", ctx)
	if r.Pre != nil {
		for _, pre := range r.Pre {
			pre(r, ctx)
		}
	}
	r.Handler(ctx)
	if r.Post != nil {
		for _, post := range r.Post {
			post(r, ctx)
		}
	}
}

// // addRespond 添加匹配器，为每个插件添加不同和适配器
// func (b *Bot) addRespond(m *Responder) {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
// 	runningPlugin.Respond = append(runningPlugin.Respond, m)
// }
