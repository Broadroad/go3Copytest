package main

import (
	"time"
	"fmt"
	"log"
	"net"
	"bufio"
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


func HeartBeating(conn net.Conn, readerChannel chan byte,timeout int) {
		select {
		case <-readerChannel:
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			break
		case <-time.After(time.Second*5):
			conn.Close()
		}
 
	}

func handleMessage(n []byte,mess chan byte){
	for _ , v := range n{
		mess <- v
	}
	close(mess)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(1000 * time.Second))
	for {
		message1, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err, message1)
			break
		}

		//message := make(chan byte)
//		go HeartBeating(conn, message, 10)
//		go handleMessage(message, message)
		fmt.Println(message1, "kang")
		var resp []byte = []byte("You are welcome. I'm server.")
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println("Write error:", err)
		}
		fmt.Println("connetion end")
	}
}
