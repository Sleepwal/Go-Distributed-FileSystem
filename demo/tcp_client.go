package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. 创建TCP连接
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
}
