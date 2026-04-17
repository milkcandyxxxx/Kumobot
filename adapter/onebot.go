package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/milkcandyxxxx/Kumobot/core"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// OneBotAdapter 适配器结构体
type OneBotAdapter struct {
	wsUrl   string
	httpURL string
	conn    *websocket.Conn
	module  []func(event *core.Event)
	prefix  string
}

// NewOneBotAdapter 新建适配器
func NewOneBotAdapter(wsUrl string, httpUrl string, prefix string) *OneBotAdapter {
	return &OneBotAdapter{
		wsUrl:   wsUrl,
		httpURL: httpUrl,
		module:  []func(event *core.Event){},
		prefix:  prefix,
	}
}

// 接口实现

// Connect 连接
func (a *OneBotAdapter) Connect() error {
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
	go a.readMessage()
	return nil
}

// Disconnect 断开连接
func (a *OneBotAdapter) Disconnect() error {
	// 避免未连接就断开
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

// readMessage 读取信息
func (a *OneBotAdapter) readMessage() {
	for {
		_, message, err := a.conn.ReadMessage()
		if err != nil {
			log.Println("获取消息失败")
		}
		var event core.Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("解析消息失败")
			continue
		}
		fmt.Printf("%+v\n", event)

		if !strings.HasPrefix(event.AltMessage, a.prefix) {
			continue
		}

		event.AltMessage = event.AltMessage[1:]

		for _, h := range a.module {
			h(&event)
		}
	}
}
func (a *OneBotAdapter) SendPrivateMessage(userID string, msg string) error {
	// 构建消息段
	payload := map[string]interface{}{
		"user_id":     userID,
		"detail_type": "private",
		"message": []map[string]interface{}{
			{
				"type": "text",
				"data": map[string]interface{}{"text": msg},
			},
		},
	}

	// 序列化
	body, _ := json.Marshal(payload)

	// 发送 POST 请求
	resp, err := http.Post(a.httpURL+"/send_message", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("onebot status error: %d", resp.StatusCode)
	}
	return nil
}
func (a *OneBotAdapter) SendGroupMessage(groupID string, msg string) error {
	// TODO: 实现 HTTP 调用

	// 构建消息段
	payload := map[string]interface{}{
		"group_id":    groupID,
		"detail_type": "group",
		"message": []map[string]interface{}{
			{
				"type": "text",
				"data": map[string]interface{}{"text": msg},
			},
		},
	}

	// 序列化
	body, _ := json.Marshal(payload)

	// 发送 POST 请求
	resp, err := http.Post(a.httpURL+"/send_message", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
