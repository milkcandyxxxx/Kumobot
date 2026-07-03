/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package bot

import "github.com/milkcandyxxxx/Kumobot/adapter"

// Handler 处理函数（回调函数定义）
type Handler func(ctx *Ctx)

// Matcher 匹配器，用于指令的匹配
type Matcher struct {
	Type        string              // 匹配类型 命令，前缀等等
	Pattern     string              // 匹配所需的关键词 如 / ! 等
	Handler     Handler             // 回调函数
	Rules       []Rule              // 规则
	ChatHistory chan *adapter.Event // 存放后续对话内容
	ChatWith    string              // 和谁在对话
	LifeCycle   bool                // 用于判断是只回复一次还是一场会话
	GroupUserID string              // 判断是那条会话
}

// addMatcher 添加匹配器
func addMatcher(m Matcher) {
	mu.Lock()
	defer mu.Unlock()
	runningPlugin.Matcher = append(runningPlugin.Matcher, m)
}
