/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 插件管理
 */

package bot

import "fmt"

// PluginInfo 插件模块信息
type PluginInfo struct {
	Name      string
	Version   string
	Author    string
	Help      string
	Priority  int          // 优先级，1-10越小越优先
	Exclusive bool         // 是否独家（不允许其他插件再次触发）
	Respond   []*Responder // 事件响应
}
type Plugin interface {
	Register(b *Bot)
}

// OnMessage 注册单次回复
func (p *PluginInfo) OnMessage(h Handler) *Responder {
	r := &Responder{
		Type:    "Message",
		Handler: h,
	}
	p.Respond = append(p.Respond, r)
	return r
}

// OnChat 注册一轮会话
func (p *PluginInfo) OnChat(h Handler) *Responder {
	r := &Responder{
		Type:      "Chat",
		Handler:   h,
		LifeCycle: true,
	}
	p.Respond = append(p.Respond, r)
	return r
}

// OnRule 添加规则
func (r *Responder) OnRule(rule ...Rule) *Responder {
	r.Rules = rule
	fmt.Println(r)
	return r
}

// OnCron 定时任务注册（6位cron表达式：秒 分 时 日 月 周）
func (p *PluginInfo) OnCron(spec string, h Handler) *Responder {
	r := &Responder{
		Type:    "Cron",
		Handler: h,
		Timing:  spec,
	}
	p.Respond = append(p.Respond, r)
	return r
}

// ++++++++++++++++++++++++定时任务++++++++++++++++++++++++
// OnCron 定时任务注册（6位cron表达式：秒 分 时 日 月 周）
// func (b *Bot) OnCron(spec string, h Handler) {
// 	b.addMatcher(Responder{
// 		Type:    "cron",
// 		Pattern: spec,
// 		Handler: h,
// 	})
// }
