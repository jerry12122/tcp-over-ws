package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	wsAddr := "127.0.0.1:8082"
	tcpAddr := "127.0.0.1:3333"

	// 建立websocket listener
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrade to websocket error:", err)
			return
		}

		// 建立tcp連接，並保留
		tcpConn, err := net.Dial("tcp", tcpAddr)
		if err != nil {
			fmt.Println("dial tcp error:", err)
			conn.Close()
			return
		}

		// 保留 WebSocket 和 TCP 連線
		type connection struct {
			conn      *websocket.Conn
			tcpConn   net.Conn
			writeChan chan []byte
			closed    bool
		}

		c := &connection{
			conn:      conn,
			tcpConn:   tcpConn,
			writeChan: make(chan []byte, 100),
			closed:    false,
		}

		// 開始讀取 WebSocket 訊息，並寫入 TCP 連線
		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
						fmt.Println("read from websocket error:", err)
					}
					c.closed = true
					break
				}

				// 將訊息寫入到 TCP 連線中
				_, err = tcpConn.Write(msg)
				if err != nil {
					fmt.Println("write to tcp error:", err)
					c.closed = true
					break
				}
			}

			// 關閉 TCP 連線
			tcpConn.Close()
		}()

		// 開始讀取 TCP 回傳訊息，並寫入 WebSocket 連線
		go func() {
			for {
				buf := make([]byte, 1024)
				n, err := tcpConn.Read(buf)
				if err != nil {
					if !c.closed {
						fmt.Println("read from tcp error:", err)
					}
					break
				}

				// 將訊息寫入到 WebSocket 連線中
				select {
				case c.writeChan <- buf[:n]:
				default:
					fmt.Println("websocket write channel full")
				}
			}

			// 關閉 WebSocket 連線
			conn.Close()
		}()

		// 開始從 writeChan 中讀取訊息，並寫入到 WebSocket 連線中
		for {
			select {
			case msg := <-c.writeChan:
				err := conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Println("write to websocket error:", err)
					c.closed = true
					break
				}
			case <-time.After(10 * time.Second):
				if c.closed {
					break
				}
			}

			if c.closed {
				break
			}
		}
	})

	fmt.Println("websocket server start at", wsAddr)
	http.ListenAndServe(wsAddr, nil)
}