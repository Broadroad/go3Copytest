package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	addr := "0.0.0.0:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Println("rpc listening", addr)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 512)
		conn.Read(buffer)

		var resp []byte = []byte("You are welcome. I'm server.")
		_, err := conn.Write(resp)
		if err != nil {
			fmt.Println("Write error:", err)
		}
		fmt.Println("connetion end")
	}
}
