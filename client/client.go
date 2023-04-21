package main

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

func main() {
	// 建立tcp listener
	listener, err := net.Listen("tcp", "127.0.0.1:3333")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// 監聽tcp連接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		// 建立websocket連接
		wsConn, _, err := websocket.DefaultDialer.Dial("ws://ws.local:8082/", nil)
		if err != nil {
			fmt.Println("dial websocket error:", err)
			conn.Close()
			continue
		}
		fmt.Println("websocket connection established")

		// 建立tcp讀取goroutine
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println("read from tcp error:", err)
					conn.Close()
					wsConn.Close()
					fmt.Println("websocket connection closed")
					return
				}
				err = wsConn.WriteMessage(websocket.TextMessage, buf[:n])
				if err != nil {
					fmt.Println("write to websocket error:", err)
					conn.Close()
					wsConn.Close()
					fmt.Println("websocket connection closed")
					return
				}
			}
		}()
	}
}