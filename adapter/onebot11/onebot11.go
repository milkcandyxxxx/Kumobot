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
	"strconv"
)

// Adapter OneBot11Adapter 适配器结构体
type Adapter struct {
	wsUrl   string
	httpURL string
	conn    *websocket.Conn
	token   string
}

// SendPrivateMessage 发送私聊信息
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

// SendGroupMessage  发送群组信息
func (a *Adapter) SendGroupMessage(atUserID string, groupID string, msg string) error {
	fmt.Println("触发SendGroupMessage")
	// 构建消息段,因为总event获取的都是string，但是onebot11协议数字类型都是int
	at, _ := strconv.ParseInt(atUserID, 10, 64)
	g, _ := strconv.ParseInt(groupID, 10, 64)
	payload := map[string]interface{}{}
	if atUserID == "" {
		payload = map[string]interface{}{
			"action": "send_group_msg",
			"params": map[string]interface{}{
				"group_id": g,
				"message": []interface{}{
					map[string]interface{}{
						"type": "text",
						"data": map[string]interface{}{
							"text": msg,
						},
					},
				},
			},
		}
	} else {
		payload = map[string]interface{}{
			"action": "send_group_msg",
			"params": map[string]interface{}{
				"group_id": g,
				"message": []interface{}{
					map[string]interface{}{
						"type": "at",
						"data": map[string]interface{}{
							"qq": at,
						},
					},
					map[string]interface{}{
						"type": "text",
						"data": map[string]interface{}{
							"text": msg,
						},
					},
				},
			},
		}
	}

	// 序列化
	body, _ := json.Marshal(payload)
	err := a.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return err
	}
	return nil
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
func NewOneBotAdapter(c adapter.Config) *Adapter {
	return &Adapter{
		wsUrl:   c.Onebots.WsURL,
		httpURL: c.Onebots.HttpURL,
		token:   c.Onebots.Token,
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
	header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
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
				switch oB11BaseMessageEvent.MessageType {
				case "private":
					var oB11PrivateMessageEvent OB11PrivateMessageEvent
					err = json.Unmarshal(message, &oB11PrivateMessageEvent)
					event := oB11PrivateMessageEvent.ToAdapterEvent()

					fmt.Println(event)
					return event, nil
				case "group":
					var oB11GroupMessageEvent OB11GroupMessageEvent
					err = json.Unmarshal(message, &oB11GroupMessageEvent)
					event := oB11GroupMessageEvent.ToAdapterEvent()
					fmt.Println(event)
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
