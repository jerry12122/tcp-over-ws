package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	// 建立websocket listener
	upgrader := websocket.Upgrader{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrade to websocket error:", err)
			return
		}
		defer conn.Close()

		// 接收websocket訊息，建立tcp連接並寫入訊息
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read from websocket error:", err)
				return
			}
			tcpConn, err := net.Dial("tcp", "127.0.0.1:3333")
			if err != nil {
				fmt.Println("dial tcp error:", err)
				continue
			}
			_, err = tcpConn.Write(msg)
			if err != nil {
				fmt.Println("write to tcp error:", err)
				continue
			}
			tcpConn.Close()
		}
	})
	http.ListenAndServe(":8082", nil)
}