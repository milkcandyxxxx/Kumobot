/**
 * @author milkcandy
 * @date 2026/4/13
 * @description 框架内置的插件函数
 */

package plugin

// OnCommand 单匹配
func OnCommand(cmd string, h Handler) {
	addMatcher(Matcher{
		Type:      "cmd",
		Pattern:   cmd,
		Priority:  0,
		Exclusive: false,
		Handler:   h,
	})
}
