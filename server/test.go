package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// 创建本地 TCP 服务器
	ln, err := net.Listen("tcp", "127.0.0.1:3333")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Println("Local TCP server started on 127.0.0.1:3333")

	// 在新建的 goroutine 中处理 TCP 连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	// 从连接中读取数据，并打印收到的内容
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
			break
		}

		data := buffer[:n]
		fmt.Printf("Received data: %s\n", string(data))
	}
}