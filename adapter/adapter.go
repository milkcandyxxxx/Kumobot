package adapter

import (
	"github.com/gorilla/websocket"
	"sync"
)

// 多平台适配器（目前只是一个平台，但是项目推荐先写出接口）

type Adapter interface {
	// Connect 连接
	Connect() error
	// Disconnect 断开连接
	Disconnect() error
	// // SendPrivateMessage 发送私人消息
	// SendPrivateMessage(userID interface{}, msg string) error
	// ReadMessage 读取消息
	ReadMessage() (interface{}, error)
	// // SendGroupMessage 发送群组消息
	// SendGroupMessage(groupID string, msg string, AtUserID string) error
	// GetSelfInfo() (SelfInfo, error)
	// GetUserInfo(userID string) (UserInfo, error)
	// DeleteMessage(messageId string) error
	// CallAction(action string, params map[string]string) (Response, error)
}

// NewOneBot11Adapter   新建适配器结构体
func NewAdapter(c Config) *AdapterInfo {
	return &AdapterInfo{
		WsUrl:   c.Onebots.WsURL,
		HttpURL: c.Onebots.HttpURL,
		Token:   c.Onebots.Token,
		Echo:    make(map[string]chan Response),
	}
}

type AdapterInfo struct {
	WsUrl   string
	HttpURL string
	Conn    *websocket.Conn
	Token   string
	Mu      sync.RWMutex
	Echo    map[string]chan Response
}
