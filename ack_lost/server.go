package main

import (
	"log"
	"net"
)

func main() {
	var c chan struct{}
	_, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	<-c
}
