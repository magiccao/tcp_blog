package main

import (
	"log"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Temporary() {
				time.Sleep(time.Millisecond * 10)
				continue
			}

			log.Fatal(err)
		}

		conn.Close()
	}
}
