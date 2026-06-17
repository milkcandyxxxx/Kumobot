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
	"time"
)

// Adapter OneBot11Adapter 适配器结构体
type Adapter struct {
	*adapter.AdapterInfo
}

func NewOneBotAdapter(c adapter.Config) *Adapter {
	return &Adapter{
		AdapterInfo: adapter.NewAdapter(c),
	}
}

// CallAction 执行特定方法用于处理某些协议特定的api
func (a *Adapter) CallAction(action string, params map[string]interface{}) (adapter.Response, error) {
	payload := map[string]interface{}{
		"action": action,
		"params": params,
	}
	echo := time.Now().UnixNano()
	echoStr := strconv.FormatInt(echo, 10)
	payload["echo"] = echoStr
	ch := make(chan adapter.Response, 1)
	a.Mu.Lock()
	a.Echo[echoStr] = ch
	a.Mu.Unlock()
	defer func() {
		a.Mu.Lock()
		delete(a.Echo, echoStr)
		a.Mu.Unlock()
	}()
	// 序列化
	body, _ := json.Marshal(payload)
	err := a.Conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return adapter.Response{}, err
	}
	select {
	case resp := <-ch:
		fmt.Println("返还值", resp)
		return resp, nil
	}
}

// SendPrivateMessage 发送私聊信息
func (a *Adapter) SendPrivateMessage(userID interface{}, msg interface{}) (int32, error) {
	// 构建消息段
	params := map[string]interface{}{
		"user_id": userID,
		"message": msg,
	}
	res, err := a.CallAction("send_private_msg", params)
	if err != nil {
		return 0, err
	}
	var Data struct {
		MsgID int32 `json:"message_id"`
	}
	if err := json.Unmarshal(res.Data, &Data); err != nil {
		return 0, err
	}
	return Data.MsgID, nil
}

// SendGroupMessage  发送群组信息
func (a *Adapter) SendGroupMessage(groupID string, msg interface{}) (int32, error) {
	// 构建消息段,因为总event获取的都是string，但是onebot11协议数字类型都是int

	group, _ := strconv.ParseInt(groupID, 10, 64)
	params := map[string]interface{}{
		"group_id": group,
		"message":  msg,
	}

	res, err := a.CallAction("send_group_msg", params)
	if err != nil {
		return 0, err
	}
	var Data struct {
		MsgID int32 `json:"message_id"`
	}
	if err := json.Unmarshal(res.Data, &Data); err != nil {
		return 0, err
	}
	return Data.MsgID, nil
}

// DeleteMsg  撤回消息
func (a *Adapter) DeleteMsg(msgID int32) error {
	params := map[string]interface{}{
		"message_id": msgID,
	}
	_, err := a.CallAction("delete_msg", params)
	if err != nil {
		return err
	}
	return nil
}

// GetMsgData GetMsg的返回结构体
type GetMsgData struct {
	Time        int32              `json:"time"`
	MessageType string             `json:"message_type"`
	MessageID   int32              `json:"message_id"`
	RealID      int32              `json:"real_id"`
	Sender      adapter.OB11Sender `json:"sender"`
	Message     interface{}        `json:"message"`
}

// GetForwardMsgData get_forward_msg 接口data统一结构体
type GetForwardMsgData struct {
	Messages []struct {
		Message []struct {
			Type string                 `json:"type"`
			Data map[string]interface{} `json:"data"`
		} `json:"message"`
	} `json:"messages"`
}

// GetMsg 获取消息
func (a *Adapter) GetMsg(msgID int32) (GetMsgData, error) {
	params := map[string]interface{}{
		"message_id": msgID,
	}
	res, err := a.CallAction("get_msg", params)
	if err != nil {
		return GetMsgData{}, err
	}
	var Data GetMsgData
	err = json.Unmarshal(res.Data, &Data)
	if err != nil {
		return GetMsgData{}, err
	}
	return Data, nil
}

// GetForwardMsg get_forward_msg
func (a *Adapter) GetForwardMsg(iD string) (GetForwardMsgData, error) {
	params := map[string]interface{}{
		"message_id": iD,
	}
	res, err := a.CallAction("get_forward_msg", params)
	if err != nil {
		return GetForwardMsgData{}, err
	}
	var Data GetForwardMsgData
	err = json.Unmarshal(res.Data, &Data)
	return Data, nil

}

// Connect 连接ws
func (a *Adapter) Connect() error {
	// 检测地址合法性
	wslAddr, err := url.Parse(a.WsUrl)
	if err != nil {

		log.Println("地址格式错误")
		return err
	}
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	conn, _, err := websocket.DefaultDialer.Dial(wslAddr.String(), header)
	if err != nil {
		log.Println("连接失败")
		return err
	}
	a.Conn = conn
	return nil
}

// Disconnect 断开连接
func (a *Adapter) Disconnect() error {
	// 避免未连接就断开
	if a.Conn != nil {
		return a.Conn.Close()
	}
	return nil
}

// ReadMessage 读取信息
func (a *Adapter) ReadMessage() (interface{}, error) {
	for {
		_, message, err := a.Conn.ReadMessage()
		log.Println(string(message))
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
			a.Mu.Lock()
			ch, ok := a.Echo[res.Echo]
			a.Mu.Unlock()
			if ok {
				ch <- res
				return res, nil
			}
			return res, nil
		}
	}
}
