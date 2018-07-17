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
	WriteNum    = 100000
	WriteThread = 1
	ConnNum     = 10
	//Server      [3]string = [3]string{"10.249.249.171:9069", "10.249.249.172:9069", "10.249.249.173:9069"}[:3]
	Server [1]string = [1]string{"127.0.0.1:8080"}

	p    map[string]pool.Pool
	text = "sadfasdfsafasdfasdfasdffffffffffffffffffffffffffffffffffffffafasdffffffffffffffffffffffffffffffffffffffffffff\n"
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
		wg := &sync.WaitGroup{}
		for _, server := range Server {
			wg.Add(1)
			conn, _ := p[server].Get()
			writer := bufio.NewWriter(conn)
			go func(writer *bufio.Writer, conn net.Conn, wg *sync.WaitGroup) {
				defer func() {
					wg.Done()
					conn.Close()
				}()
				writer.WriteString(text)
				writer.Flush()
				buffer := make([]byte, 100)
				n, err := conn.Read(buffer)
				if err != nil{
					fmt.Println(err, n)
				}
					
				fmt.Println(string(buffer))
			}(writer, conn, wg)
		}

		fmt.Println(time.Since(start).String())
		wg.Wait()
	}
	time.Sleep(time.Second)
	done <- struct{}{}
}

func main() {
	defer func() {
		for _, server := range Server {
			p[server].Close()
		}

	}()
	done := make(chan struct{})
	for i := 0; i < WriteThread; i++ {
		go write3copy(done)
	}

	for i := 0; i < WriteThread; i++ {
		<-done
	}
}
