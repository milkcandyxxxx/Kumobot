package onebot12

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Adapter OneBot12Adapter 适配器结构体
type Adapter struct {
	wsUrl   string
	httpURL string
	conn    *websocket.Conn
}

// NewOneBotAdapter  新建适配器结构体
func NewOneBotAdapter(wsUrl string, httpUrl string) *Adapter {
	return &Adapter{
		wsUrl:   wsUrl,
		httpURL: httpUrl,
	}
}

// 接口实现

// Connect 连接
func (a *Adapter) Connect() error {
	// 检测地址合法性
	wslAddr, err := url.Parse(a.wsUrl)
	if err != nil {

		log.Println("地址格式错误")
		return err
	}
	conn, _, err := websocket.DefaultDialer.Dial(wslAddr.String(), nil)
	if err != nil {
		log.Println("连接失败")
		return err
	}
	a.conn = conn
	return nil
}

// Disconnect 断开连接
func (a *Adapter) Disconnect() error {
	// 避免未连接就断开
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// ReadMessage 读取信息
func (a *Adapter) ReadMessage() (interface{}, error) {
	for {

		_, message, err := a.conn.ReadMessage()
		fmt.Println(string(message))
		if err != nil {
			log.Println("获取消息失败")
		}
		// 先放入万能字典，查看是事件还是动作响应
		var check map[string]interface{}
		err = json.Unmarshal(message, &check)
		if err != nil {
			log.Println("解析消息失败", err)
			continue
		}
		if _, exists := check["type"]; exists {
			var One12Event Event
			err = json.Unmarshal(message, &One12Event)
			if err != nil {
				log.Println("解析消息失败", err)
				continue
			}
			event := &adapter.Event{
				ID:         One12Event.ID,
				Type:       One12Event.Type,
				Self:       One12Event.Self,
				Time:       One12Event.Time,
				DetailType: One12Event.DetailType,
				SubType:    One12Event.SubType,
				MessageID:  One12Event.MessageID,
				Message:    One12Event.Message,
				UserID:     One12Event.UserID,
				GuildID:    One12Event.GuildID,
				GroupID:    One12Event.GroupID,
				AltMessage: One12Event.AltMessage,
			}
			return event, nil
		} else {
			var response adapter.Response
			err = json.Unmarshal(message, &response)
			fmt.Println(response)
			if err != nil {
				log.Println("解析消息失败")
				return nil, err
			}
			return &response, nil
		}

	}
}

// SendPrivateMessage 发送私聊消息
func (a *Adapter) SendPrivateMessage(userID string, msg string) error {
	// 构建消息段
	payload := map[string]interface{}{
		"action": "send_message",
		"params": map[string]interface{}{
			"user_id":     userID,
			"detail_type": "private",
			"message": []map[string]interface{}{
				{
					"type": "text",
					"data": map[string]interface{}{"text": msg},
				},
			},
		},
	}

	// 序列化
	body, _ := json.Marshal(payload)
	err := a.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return err
	}
	return nil
	// // 发送 POST 请求 TODO 暂留
	// resp, err := http.Post(a.httpURL+"/send_message", "application/json", bytes.NewBuffer(body))
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("onebot status error: %d", resp.StatusCode)
	// }
	// return nil
}

// SendGroupMessage 发送群聊消息
func (a *Adapter) SendGroupMessage(groupID string, msg string) error {
	// 构建消息段
	payload := map[string]interface{}{
		"action": "send_message",
		"params": map[string]interface{}{
			"group_id":    groupID,
			"detail_type": "group",
			"message": []map[string]interface{}{
				{
					"type": "text",
					"data": map[string]interface{}{"text": msg},
				},
			},
		},
	}
	// 序列化
	body, _ := json.Marshal(payload)
	// // 发送 POST 请求 TODO 暂留
	// resp, err := http.Post(a.httpURL+"/send_message", "application/json", bytes.NewBuffer(body))
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()
	// return nil
	err := a.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return err
	}
	return nil
}

// GetSelfInfo 获取机器人自身信息
func (a *Adapter) GetSelfInfo() (adapter.SelfInfRes, error) {
	res, err := http.Post(a.httpURL+"/get_self_info", "application/json", nil)
	if err != nil {
		return adapter.SelfInfRes{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var selfinfo adapter.SelfInfRes
	err = json.Unmarshal(body, &selfinfo)
	return selfinfo, nil
}

// GetUserInfo 获取用户信息
func (a *Adapter) GetUserInfo(userID string) (adapter.UserInfo, error) {
	return adapter.UserInfo{}, nil
}

func (a *Adapter) DeleteMessage(messageId string) error {
	payload := map[string]interface{}{
		"action": "delete_message",
		"params": map[string]interface{}{
			"message_id": messageId,
		},
	}
	// 序列化
	body, _ := json.Marshal(payload)
	err := a.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return err
	}
	return nil
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
