/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 插件管理
 */

package bot

import (
	"log"
	"sync"
)

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

// // Info 插件基础信息结构体
// type Info struct {
//
// }
var (
	plugins       = []*PluginInfo{} // 存储所有的匹配器注册
	runningPlugin *PluginInfo
	mu            sync.RWMutex // 加锁（目前是冷加载插件，后续热加载等需要注意）
)

// addPlugin 添加插件
func (b *Bot) addPlugin(p *PluginInfo) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.plugins = append(b.plugins, p)
	runningPlugin = p
	log.Println("\n加载插件:", p.Name, "\n版本:", p.Version, "\n作者:", p.Author, "\n帮助:", p.Help, "\n优先级:", p.Priority, "\n独家:", p.Exclusive)
}
