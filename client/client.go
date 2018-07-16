package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fatih/pool"
)

var (
	WriteNum    = 1000
	WriteThread = 10
	ConnNum     = 10
	//Server      [3]string = [3]string{"10.249.249.171:9069", "10.249.249.172:9069", "10.249.249.173:9069"}[:3]
	Server [1]string = [1]string{"127.0.0.1:8080"}

	p    map[string]pool.Pool
	text = "sadfasdfsafasdfasdfasdffffffffffffffffffffffffffffffffffffffafasdffffffffffffffffffffffffffffffffffffffffffff"
)

func init() {
	p = make(map[string]pool.Pool)
	for _, server := range Server {
		factory := func() (net.Conn, error) {
			return net.Dial("tcp", server)
		}
		pi, _ := pool.NewChannelPool(3, 10, factory)
		p[server] = pi
	}
}

func writeString(writer *bufio.Writer, conn net.Conn) {
	defer conn.Close()
	writer.WriteString(text)
	writer.Flush()
	buffer := make([]byte, 512)
	conn.Read(buffer)
}

func write3copy(done chan<- struct{}) {
	for i := 0; i < WriteNum; i++ {
		start := time.Now()
		var wg sync.WaitGroup
		for _, server := range Server {
			wg.Add(1)
			conn, _ := p[server].Get()
			writer := bufio.NewWriter(conn)
			go func(writer *bufio.Writer, conn net.Conn) {
				defer func() {
					wg.Done()
					conn.Close()
				}()
				writer.WriteString(text)
				writer.Flush()
				buffer := make([]byte, 512)
				conn.Read(buffer)
			}(writer, conn)
		}

		fmt.Println(time.Since(start).Seconds())
		wg.Wait()
	}
	done <- struct{}{}
}

func main() {
	defer func() {
		for _, server := range Server {
			p[server].Close()
		}

	}()
	time.Sleep(time.Second * 20)
	done := make(chan struct{})
	for i := 0; i < WriteThread; i++ {
		go write3copy(done)
	}

	for i := 0; i < WriteThread; i++ {
		<-done
	}
}
