package adapter

// Response 响应体的总接口

type Event struct {
	// ========== 基础字段（所有事件都有） ==========
	ID         string  `json:"id"`          // OB12独有
	Self       BotSelf `json:"self"`        // OB12独有
	Time       int64   `json:"time"`        // 共用，OB12 float64 / OB11 int64
	Type       string  `json:"type"`        // 事件大类，OB12 type / OB11 post_type
	DetailType string  `json:"detail_type"` // 详细类型，OB12 detail_type / OB11 message_type或notice_type
	SubType    string  `json:"sub_type"`    // 子类型，共用
	SelfId     int64   `json:"self_id"`     // OB11 机器人自身 QQ 号

	Sender OB11Sender `json:"sender"`

	// ========== 消息事件字段 ==========
	MessageID interface{}      `json:"message_id"` // 消息ID，OB12 string / OB11 int32
	Message   []MessageSegment `json:"message"`    // 消息段数组，共用；OB11 CQ码时转单text段
	UserID    interface{}      `json:"user_id"`    // 发送者ID，OB12顶层 / OB11可能来自sender
	GroupID   string           `json:"group_id"`   // 群ID，群聊才有
	GuildID   string           `json:"guild_id"`   // OB12独有，频道ID

	// ========== onebots平台扩展 ==========
	AltMessage string `json:"alt_message,omitempty"` // 共用 OB12纯文本（onebots扩展），OB11可用raw_message
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
type Response struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    struct {
		MessageId string  `json:"message_id"`
		Time      float64 `json:"time"`
	} `json:"data"`
	Message string `json:"message"`
	Echo    string `json:"echo"`
}
