package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 建立 TCP 连接
	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// 从标准输入中读取输入，并将其发送到 TCP 连接
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		// 发送数据
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error sending:", err.Error())
			break
		}
	}
}