/**
 * @author milkcandy
 * @date 2026/6/15
 * @description TODO
 */

package ON11

import (
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"github.com/milkcandyxxxx/Kumobot/adapter/onebot11"
	"strconv"
)

type ON11 struct {
	*onebot11.Adapter
	Event *adapter.Event
}

// // SendGroupMessage  发送群里信息
// func (o *ON11) SendGroupMessagea(groupID string, msg string) (int32, error) {
// 	return o.SendGroupMessage(groupID, msg)
// }

// SendGroupMessageAt  发送群聊信息，并at回复
func (o *ON11) SendGroupMessageAt(groupID string, msg string) (int32, error) {
	// 构造at消息段
	UserID, _ := strconv.ParseInt(o.Event.UserID, 10, 64)
	msgAt := []map[string]interface{}{
		{
			"type": "at",
			"data": map[string]interface{}{
				"qq": UserID,
			},
		},
		{
			"type": "text",
			"data": map[string]interface{}{
				"text": msg,
			},
		},
	}
	return o.SendGroupMessage(groupID, msgAt)
}

// Send 一键发送默认为回复（在哪触发的哪里回复，默认at）
func (o *ON11) Send(msg string) (int32, error) {
	if o.Event.DetailType == "private" {
		return o.SendPrivateMessage(o.Event.UserID, msg)
	}
	if o.Event.DetailType == "channel" || o.Event.DetailType == "group" {
		return o.SendGroupMessageAt(o.Event.GroupID, msg)
	}
	return 0, nil
}
