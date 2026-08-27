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
	ReadMessage(chan *Event) error
	// CheckHeartbeat(int64, time.Duration) 心跳检查
	// // SendGroupMessage 发送群组消息
	// SendGroupMessage(groupID string, msg string, AtUserID string) error
	// GetSelfInfo() (SelfInfo, error)
	// GetUserInfo(userID string) (UserInfo, error)
	// DeleteMessage(messageId string) error
	// CallAction(action string, params map[string]string) (Response, error)
}

// NewOneBot11Adapter   新建适配器结构体
func InitializeAdapterInformation(agreementName string) *AdapterInfo {
	return &AdapterInfo{
		WsUrl: GlobalConfig.Agreement[agreementName].WebSocket.WsURL,
		// HttpURL:      GlobalConfig.Agreement[agreementName].WsURL,
		Token:        GlobalConfig.Agreement[agreementName].WebSocket.Token,
		Echo:         make(map[string]chan Response),
		ProtocolName: "ON11",
		// HeartbeatTime: -1,心跳检查
	}
}

type AdapterInfo struct {
	WsUrl        string
	HttpURL      string
	Conn         *websocket.Conn
	Token        string
	Mu           sync.RWMutex
	Echo         map[string]chan Response
	ProtocolName string
	// HeartbeatTime int64 心跳检查
}
