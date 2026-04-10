package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strings"
	"time"
)

var config Config
var httpClient *resty.Client // 使用resty库来刚方便的发送请求

// 配置文件结构体
type Config struct {
	Onebots OnebotsConfig `yaml:"onebots"`
}
type OnebotsConfig struct {
	WsUrl   string `mapstructure:"ws_url"`
	HttpUrl string `mapstructure:"http_url"`
}

// // Event 信息数据,目前仅为简单学习（并不全）
// type SimpleEvent struct {
// 	ID         string `json:"id"`          // 事件 ID
// 	Type       string `json:"type"`        // 事件类型：message, notice, request
// 	DetailType string `json:"detail_type"` // 详细类型：private, group, channel
// 	UserID     string `json:"user_id"`     // 发送者 ID
// 	GroupID    string `json:"group_id"`    // 群 ID
// 	Message    string `json:"alt_message"` // 消息内容（纯文本，onebots 扩展）
// 	Platform   string `json:"platform"`    // 平台：kook, qq, discord（onebots 扩展）
// }

// Event标准的onebot12
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

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	log.Println("云宝~启动中！o((>ω< ))o")
	err := loadConfig()
	fmt.Println(config)
	httpClient = resty.New()
	if err != nil {
		log.Fatal("配置文件解析失败", err)
	}
	// 启动websocket连接
	go connectWebSocket(config.Onebots.WsUrl)
	select {}
}

// connectWebSocket websocket连接
func connectWebSocket(addr string) {
	for {
		// 检测地址合法性
		addrStruct, err := url.Parse(addr)
		if err != nil {
			log.Println("地址解析失败", err)
			continue
		}
		conn, _, err := websocket.DefaultDialer.Dial(addrStruct.String(), nil)
		if err != nil {
			log.Println("建立连接失败", err)
			log.Println("云宝将在5秒后尝试重新连接！~(￣▽￣)~*")
			time.Sleep(5 * time.Second)
			continue
		}
		readWebSocket(conn)
		conn.Close()
		log.Println("连接意外中断，云宝将在5秒后尝试重新连接！=￣ω￣=")
	}
}

// readWebSocket 读取消息
func readWebSocket(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取信息失败", err.Error())
			return
		}
		log.Println(string(message))
		event, err := parseEvent(message)
		if err != nil {
			log.Println("解析事件失败", err)
		}
		handleEvent(event)
	}
}

// parseEvent 解析消息[]byts为json
func parseEvent(data []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// handleEvent 事件处理
func handleEvent(event *Event) {
	// 目前只处理事件
	if event.Type != "message" {
		return
	}
	// 简单的连接测试echo
	msg := processCommand(event.AltMessage)
	sendMessage(event, msg)
}

// processCommand 命令处理
func processCommand(msg string) string {
	//  首尾空格
	msg = strings.TrimSpace(msg)
	// 是否以"/"开头
	if !strings.HasPrefix(msg, "/echo") {
		return ""
	}
	returnedMsg := msg[5:]

	return returnedMsg
}

// sendMessage 发送消息
func sendMessage(event *Event, msg string) {
	// {
	// 	"action": "send_message",
	// 	"params": {
	// 	"detail_type": "private",
	// 		"user_id": "123445667",
	// 		"message": [
	// {
	// "type": "text",
	// "data": {
	// "text": "嗨～"
	// }
	// }
	// ]
	// },
	// "echo": "1234"
	// }
	// 动作请求的标准格式，其中"action": "send_message",可以变成 "/send_message"在http中，"echo": "1234"为服务器回复，主要用于webs连接，可选
	url := config.Onebots.HttpUrl + "/send_message"
	var body map[string]interface{}
	// 定义消息结构体
	message := MessageSegment{
		Type: "text",
		Data: map[string]interface{}{
			"text": msg,
		},
	}
	// message_group := MessageSegment{
	// 	Type: "text",
	// 	Data: map[string]interface{}{
	// 		"text": msg + "群组",
	// 	},
	// }
	// 如果是私聊
	if event.DetailType == "private" {
		body = map[string]interface{}{
			"detail_type": event.DetailType,
			"user_id":     event.UserID,
			"message":     []MessageSegment{message},
		}
		// 这里一张少东西，后面研究一下
		// onbots里没有严格按照kook的guild_id（服务器）channel_id（频道），而是统一归为group（注：但实际是不对等的，因为数据接收时类型是"detail_type":"channel，）
	} else if event.DetailType == "channel" {
		body = map[string]interface{}{
			"detail_type": "group",
			"group_id":    event.GroupID,
			"message":     []MessageSegment{message},
		}
	}
	// 发送消息
	resp, err := httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		log.Println("发送错误", err)
	}
	if resp.StatusCode() != 200 {
		log.Println("发送失败", resp.String())
	} else {
		log.Println("发送成功", message.Data["text"])
	}
}
