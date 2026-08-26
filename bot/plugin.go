/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 插件管理
 */

package bot

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

// OnMessage 为插件注册功能
func (p *PluginInfo) OnMessage(h Handler) *Responder {
	r := &Responder{Handler: h}
	p.Respond = append(p.Respond, r)
	return r
}
func (p *PluginInfo) OnChat(h Handler) *Responder {
	r := &Responder{
		Type:      "Chat",
		Handler:   h,
		LifeCycle: true,
	}
	p.Respond = append(p.Respond, r)
	return r
}
