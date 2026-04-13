package core

type BotSelf struct {
	Platform string `json:"platform"` // 平台名称
	UserID   string `json:"user_id"`  // 机器人 ID
}

// MessageSegment 消息段（OneBot 12 标准）
// OneBot 12 的消息是由多个消息段组成的数组
type MessageSegment struct {
	Type string                 `json:"type"` // 消息段类型：text, image, at, face 等
	Data map[string]interface{} `json:"data"` // 消息段数据
}

// Event OneBot 12 标准事件结构
/* message.private 私聊消息标准
{
    "id": "b6e65187-5ac0-489c-b431-53078e9d2bbb",
    "self": {
        "platform": "qq",
        "user_id": "123234"
    },
    "time": 1632847927.599013,
    "type": "message",
    "detail_type": "private",
    "sub_type": "",
    "message_id": "6283",
    "message": [
        {
            "type": "text",
            "data": {
                "text": "OneBot is not a bot"
            }
        },
        {
            "type": "image",
            "data": {
                "file_id": "e30f9684-3d54-4f65-b2da-db291a477f16"
            }
        }
    ],
    "alt_message": "OneBot is not a bot[图片]",
    "user_id": "123456788"
}
*/
type Event struct {
	// ========== 基础字段（所有事件都有） ==========
	ID         string  `json:"id"` // 事件 ID
	Self       BotSelf `json:"self"`
	Time       int64   `json:"time"`        // 事件时间戳（秒）
	Type       string  `json:"type"`        // 事件类型：message, notice, request, meta
	DetailType string  `json:"detail_type"` // 详细类型
	SubType    string  `json:"sub_type"`    // 子类型
	// 机器人自身信息

	// ========== 消息事件字段 ==========
	MessageID string           `json:"message_id"` // 消息 ID（用于撤回、引用等）
	Message   []MessageSegment `json:"message"`    // 消息内容（消息段数组）
	UserID    string           `json:"user_id"`    // 发送者 ID
	GroupID   string           `json:"group_id"`   // 群 ID（群聊消息）

	GuildID string `json:"guild_id"`

	// ==========onebots平台提供的字段非标准onebots==========
	AltMessage string `json:"alt_message,omitempty"` // 纯文本消息（onebots 扩展）
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

// GetPlatform 获取平台名称
func (e *Event) GetPlatform() string {
	return e.Self.Platform
}
