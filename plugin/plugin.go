/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 插件管理
 */

package plugin

import (
	"log"
	"sync"
)

// Plugin 插件总结构体储存信息
type Plugin struct {
	Name    string
	Version string
	Author  string
	Help    string
	Matcher []Matcher
}

// // Info 插件基础信息结构体
// type Info struct {
//
// }
var (
	plugins       = []*Plugin{} // 存储所有的匹配器注册
	runningPlugin *Plugin
	mu            sync.RWMutex // 加锁（目前是冷加载插件，后续热加载等需要注意）
)

// addPlugin 添加插件
func addPlugin(p *Plugin) {
	mu.Lock()
	defer mu.Unlock()
	plugins = append(plugins, p)
	runningPlugin = p
	log.Println("\n加载插件:", p.Name, "\n版本:", p.Version, "\n作者:", p.Author)
}
func GetPluginName() string {
	return runningPlugin.Name
}

func OnPlugin(info ...string) {
	thisPlugin := &Plugin{
		Name:    info[0],
		Version: info[1],
		Author:  info[2],
		Help:    info[3],
	}
	addPlugin(thisPlugin)
}
