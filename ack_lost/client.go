package main

import (
	"log"
	"net"
	"time"
)

func main() {
	var c chan struct{}
	for i := 0; i < 4; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:8081", time.Second*20)
		if err != nil {
			log.Println(err)
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				time.Sleep(time.Millisecond * 50)
				continue
			}
			return
		}

		go func(conn net.Conn) {
			log.Println(conn.LocalAddr())
			n, err := conn.Write([]byte("hello, world"))
			log.Printf("write %d bytes, nil: %v", n, err)

			for {
				time.Sleep(time.Second)
			}
			conn.Close()
		}(conn)
	}

	<-c
}
