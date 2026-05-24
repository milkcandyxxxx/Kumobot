/**
 * @author milkcandy
 * @date 2026/5/23
 * @description TODO
 */

package onebot12

import "github.com/milkcandyxxxx/Kumobot/adapter"

type Event struct {
	// ========== 基础字段（所有事件都有） ==========
	ID         string          `json:"id"` // 事件 ID
	Self       adapter.BotSelf `json:"self"`
	Time       int64           `json:"time"`        // 事件时间戳（秒）
	Type       string          `json:"type"`        // 事件类型：message, notice, request, meta
	DetailType string          `json:"detail_type"` // 详细类型
	SubType    string          `json:"sub_type"`    // 子类型
	// 机器人自身信息
	// ========== 消息事件字段 ==========
	MessageID string                   `json:"message_id"` // 消息 ID（用于撤回、引用等）
	Message   []adapter.MessageSegment `json:"message"`    // 消息内容（消息段数组）
	UserID    string                   `json:"user_id"`    // 发送者 ID
	GroupID   string                   `json:"group_id"`   // 群 ID（群聊消息）
	GuildID   string                   `json:"guild_id"`
	// ==========onebots平台提供的字段非标准onebots==========
	AltMessage string `json:"alt_message,omitempty"` // 纯文本消息（onebots 扩展）
}

// GetPlatform 获取平台名称
func (e *Event) GetPlatform() string {
	return e.Self.Platform
}
