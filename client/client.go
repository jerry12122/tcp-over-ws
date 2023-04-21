package main

import (
	"fmt"
	"net"
	"time"
	"io"
	"github.com/gorilla/websocket"
)

func main() {
	wsURL := "ws://ws.local:8082/"
	tcpAddr := "127.0.0.1:3333"

	// 建立tcp listener
	listener, err := net.Listen("tcp", tcpAddr)
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
		wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			fmt.Println("dial websocket error:", err)
			conn.Close()
			continue
		}
		fmt.Println("websocket connection established")

		// 建立tcp讀取goroutine
		go func() {
			defer wsConn.Close()
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println("read from tcp error:", err)
					conn.Close()
					fmt.Println("tcp connection closed")
					time.Sleep(10 * time.Second)
					if isConnClosed(conn) {
						fmt.Println("websocket connection closed")
						return
					}
					continue
				}
				err = wsConn.WriteMessage(websocket.TextMessage, buf[:n])
				if err != nil {
					fmt.Println("write to websocket error:", err)
					conn.Close()
					fmt.Println("tcp connection closed")
					time.Sleep(10 * time.Second)
					if isConnClosed(conn) {
						fmt.Println("websocket connection closed")
						return
					}
					continue
				}
			}
		}()
	}
}
func isConnClosed(conn net.Conn) bool {
    var one []byte
    // 设置一个短的读取超时时间
    conn.SetReadDeadline(time.Now().Add(time.Millisecond))
    // 尝试从连接中读取一个字节
    _, err := conn.Read(one)
    if err != nil && err != io.EOF {
        return true
    }
    return false
}