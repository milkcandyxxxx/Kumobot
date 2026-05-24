/**
 * @author milkcandy
 * @date 2026/5/22
 * @description TODO
 */

package onebot11

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"log"
	"net/http"
	"net/url"
)

// Adapter OneBot11Adapter 适配器结构体
type Adapter struct {
	wsUrl   string
	httpURL string
	conn    *websocket.Conn
}

func (a *Adapter) SendPrivateMessage(userID interface{}, msg string) error {
	// 构建消息段
	payload := map[string]interface{}{
		"action": "send_private_msg",
		"params": map[string]interface{}{
			"user_id": userID,
			"message": msg,
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

func (a *Adapter) SendGroupMessage(groupID string, msg string) error {
	// TODO implement me
	panic("implement me")
}

func (a *Adapter) GetSelfInfo() (adapter.SelfInfRes, error) {
	// TODO implement me
	panic("implement me")
}

func (a *Adapter) GetUserInfo(userID string) (adapter.UserInfo, error) {
	// TODO implement me
	panic("implement me")
}

func (a *Adapter) DeleteMessage(messageId string) error {
	// TODO implement me
	panic("implement me")
}

// NewOneBot11Adapter   新建适配器结构体
func NewOneBotAdapter(wsUrl string, httpUrl string) *Adapter {
	return &Adapter{
		wsUrl:   wsUrl,
		httpURL: httpUrl,
	}
}
func (a *Adapter) Connect() error {
	// 检测地址合法性
	wslAddr, err := url.Parse(a.wsUrl)
	if err != nil {

		log.Println("地址格式错误")
		return err
	}
	header := http.Header{}
	header.Set("Authorization", "Bearer ~bcxCaBq0j4sZGjt")
	conn, _, err := websocket.DefaultDialer.Dial(wslAddr.String(), header)
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
		if _, exists := check["self_id"]; exists {
			// 看事件类型
			var oneBotEvent OneBotEvent
			err = json.Unmarshal(message, &oneBotEvent)
			if err != nil {
				log.Println("解析消息失败")
				return nil, err
			}
			// 为消息
			if oneBotEvent.PostType == "message" {
				var oneBotEvent OneBotEvent
				err = json.Unmarshal(message, &oneBotEvent)
				var oB11BaseMessageEvent OB11BaseMessageEvent
				err = json.Unmarshal(message, &oB11BaseMessageEvent)
				if err != nil {
					log.Println("解析消息失败")
					return nil, err
				}
				// 为私聊
				if oB11BaseMessageEvent.MessageType == "private" {
					var oB11PrivateMessageEvent OB11PrivateMessageEvent
					err = json.Unmarshal(message, &oB11PrivateMessageEvent)
					event := &adapter.Event{
						Time:       oB11PrivateMessageEvent.Time,
						SelfId:     oB11PrivateMessageEvent.SelfId,
						Type:       oB11PrivateMessageEvent.PostType,
						DetailType: oB11PrivateMessageEvent.MessageType,
						UserID:     oB11PrivateMessageEvent.UserId,
						Message:    oB11PrivateMessageEvent.Message,
						AltMessage: oB11PrivateMessageEvent.RawMessage,
						SubType:    oB11PrivateMessageEvent.SubType,
						Sender:     oB11PrivateMessageEvent.Sender,
					}
					return event, nil
				}

			}

		} else {
			var res adapter.Response
			err = json.Unmarshal(message, &res)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
}
