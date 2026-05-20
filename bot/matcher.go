/**
 * @author milkcandy
 * @date 2026/4/13
 * @description TODO
 */

package bot

// Handler 处理函数（回调函数定义）
type Handler func(b *Bot)

// Matcher 匹配器，用于指令的匹配
type Matcher struct {
	Type    string  // 匹配类型 命令，前缀等等
	Pattern string  // 匹配所需的关键词 如 / ! 等
	Handler Handler // 回调函数
}

// addMatcher 添加匹配器
func addMatcher(m Matcher) {
	mu.Lock()
	defer mu.Unlock()
	runningPlugin.Matcher = append(runningPlugin.Matcher, m)
}
