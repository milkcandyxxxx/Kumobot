/**
 * @author milkcandy
 * @date 2026/5/23
 * @description TODO
 */

package onebot11

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"strconv"
)

// OneBotEvent 所有事件的基类，定义了共有字段。
type OneBotEvent struct {
	Time     int64  `json:"time"`      // 事件发生的时间戳
	SelfID   int64  `json:"self_id"`   // 机器人自身 QQ 号
	PostType string `json:"post_type"` // 事件类型
}

// OB11BaseMessageEvent 所有消息事件的基类
type OB11BaseMessageEvent struct {
	OneBotEvent
	// 具体消息类型在子类中定义
	MessageType string                   `json:"message_type"` // 消息类型
	MessageID   int64                    `json:"message_id"`   // 消息 ID
	UserID      int64                    `json:"user_id"`      // 发送者 QQ 号
	Message     []adapter.MessageSegment `json:"message"`      // 消息内容
	RawMessage  string                   `json:"raw_message"`  // 原始消息内容
}

// OB11PrivateMessageEvent 私聊消息事件
type OB11PrivateMessageEvent struct {
	OB11BaseMessageEvent
	SubType string             `json:"sub_type"`
	Sender  adapter.OB11Sender `json:"sender"`
}

// ToAdapterEvent 将 OB11 的私聊事件通用的 adapter.Event
func (o *OB11PrivateMessageEvent) ToAdapterEvent() *adapter.Event {
	return &adapter.Event{
		Time:       o.Time,
		SelfID:     o.SelfID,
		Type:       o.PostType,
		DetailType: o.MessageType,
		MessageID:  strconv.FormatInt(o.MessageID, 10),
		Message:    o.Message,
		UserID:     strconv.FormatInt(o.UserID, 10),
		AltMessage: o.RawMessage,
		SubType:    o.SubType,
		Sender:     o.Sender,
	}
}

// OB11GroupMessageEvent 群聊消息事件
type OB11GroupMessageEvent struct {
	OB11BaseMessageEvent
	Sender    adapter.OB11Sender `json:"sender"`
	Anonymous any                `json:"anonymous"`
	GroupID   int64              `json:"group_id"`
}

func (o *OB11GroupMessageEvent) ToAdapterEvent() *adapter.Event {
	return &adapter.Event{
		Time:       o.Time,
		SelfID:     o.SelfID,
		Type:       o.PostType,
		DetailType: o.MessageType,
		MessageID:  strconv.FormatInt(o.MessageID, 10),
		Message:    o.Message,
		UserID:     strconv.FormatInt(o.UserID, 10),
		AltMessage: o.RawMessage,
		Sender:     o.Sender,
		Anonymous:  o.Anonymous,
		GroupID:    strconv.FormatInt(o.GroupID, 10),
	}
}
