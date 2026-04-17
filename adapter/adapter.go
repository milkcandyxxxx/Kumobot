package adapter

import "github.com/milkcandyxxxx/Kumobot/core"

// 多平台适配器（目前只是一个平台，但是项目推荐先写出接口）

type Adapter interface {
	// Connect 连接
	Connect() error
	// Disconnect 断开连接
	Disconnect() error
	// SendPrivateMessage 发送私人消息
	SendPrivateMessage(userID string, msg string) error
	// SendGroupMessage 发送群组消息
	SendGroupMessage(groupID string, msg string) error
	// OnEvent 注册事件回调
	OnEvent(module func(event *core.Event))
}

// OnEvent 注册事件监听函数
func (a *OneBotAdapter) OnEvent(module func(event *core.Event)) {
	a.module = append(a.module, module)
}
