/**
 * @author milkcandy
 * @date 2026/5/23
 * @description TODO
 */

package onebot11

import "github.com/milkcandyxxxx/Kumobot/adapter"

type OB11BaseMetaEvent struct {
}

type OneBotEvent struct {
	Time     int64  `json:"time"`      // 事件发生的时间戳
	SelfId   int64  `json:"self_id"`   // 机器人自身 QQ 号
	PostType string `json:"post_type"` // 事件类型
}

// OB11BaseMessageEvent 所有消息事件的基类
type OB11BaseMessageEvent struct {
	OneBotEvent
	// 具体消息类型在子类中定义
	MessageType string                   `json:"message_type"` // 消息类型
	MessageId   int64                    `json:"message_id"`   // 消息 ID
	UserId      int64                    `json:"user_id"`      // 发送者 QQ 号
	Message     []adapter.MessageSegment `json:"message"`      // 消息内容
	RawMessage  string                   `json:"raw_message"`  // 原始消息内容
}

// OB11Sender 发送者信息

// OB11PrivateMessageEvent 私聊消息事件
type OB11PrivateMessageEvent struct {
	OB11BaseMessageEvent
	SubType string             `json:"sub_type"`
	Sender  adapter.OB11Sender `json:"sender"`
}
