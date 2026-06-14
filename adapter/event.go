package adapter

import "encoding/json"

// Response 响应体的总接口
type Response struct {
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data"`
	Echo   string          `json:"echo"`
}

// type ResponseData struct {
// 	UserID   int64  `json:"user_id"`
// 	Nickname string `json:"nickname"`
// }

// Event 上下文总接口
type Event struct { // 共用字段
	Time       int64            // 事件时间戳
	SubType    string           // 事件子类型
	MessageID  string           // 消息唯一ID
	Message    []MessageSegment // 消息段数组
	UserID     string           // 发送者账号
	GroupID    string           // 群组账号
	AltMessage string           // 纯文本
	DetailType string           // 事件细分类型,OB11为MessageType
	// OneBot11 专属字段
	PostType    string     // 事件大类
	MessageType string     // 消息类型
	SelfID      int64      // 机器人自身账号
	Anonymous   any        // 匿名发言信息
	Sender      OB11Sender // 发送者详情

	// OneBot12 专属字段
	ID   string  // 事件标识
	Self BotSelf // 自身账号信息
	Type string  // 事件主类型

	GuildID string // 频道ID
}

type BotSelf struct {
	Platform string `json:"platform"` // OB12 self.platform；OB11填"qq"
	UserID   string `json:"user_id"`  // OB12 self.user_id；OB11 self_id转来
}

// MessageSegment 消息段
// OB12标准，OB11 array格式时相同
type MessageSegment struct {
	Type string                 `json:"type"` // text / image / at / face ...
	Data map[string]interface{} `json:"data"` // 具体数据
}

// GetMessageText 获取纯文本，处理部分平台无AltMessage
func (e *Event) GetMessageText() string {
	if e.AltMessage != "" {
		return e.AltMessage
	}
	var msgAll string
	for _, msg := range e.Message {
		if msg.Type == "text" {
			msgAll += msg.Data["text"].(string)

		}
	}
	return msgAll
}

type OB11Sender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
}

type SelfInfo struct {
	UserID   string // OB11 QQ 号
	Nickname string // OB11 QQ 昵称
}
