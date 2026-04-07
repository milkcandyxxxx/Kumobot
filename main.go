package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

func main() {
	log.Println("云宝~启动中！o((>ω< ))o")
	wsAddress := "ws://127.0.0.1:6727/kook/46609/onebot/v12"
	// 启动websocket连接
	go connectWebSocket(wsAddress)
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
			continue
		}
		readWebSocket(conn)
		conn.Close()
	}
}
func readWebSocket(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取信息失败", err.Error())
			return
		}

		log.Println(string(message))

	}
}
