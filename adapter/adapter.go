package adapter

// 多平台适配器（目前只是一个平台，但是项目推荐先写出接口）

type Adapter interface {
	// Connect 连接
	Connect() error
	// Disconnect 断开连接
	Disconnect() error
	// SendPrivateMessage 发送私人消息
	SendPrivateMessage(userID interface{}, msg string) error
	// ReadMessage 读取消息
	ReadMessage() (interface{}, error)
	// SendGroupMessage 发送群组消息
	SendGroupMessage(groupID string, msg string) error
	GetSelfInfo() (SelfInfRes, error)
	GetUserInfo(userID string) (UserInfo, error)
	DeleteMessage(messageId string) error
}
