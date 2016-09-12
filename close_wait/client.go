package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.AF_UNSPEC)
	if err != nil {
		log.Fatal(err)
	}

	lsa, err := ResolveSockaddr("127.0.0.1:49622")
	if err != nil {
		log.Fatal("parse local addr: ", err.Error())
	}

	if err := syscall.Bind(fd, lsa); err != nil {
		log.Fatal("Bind: ", err.Error())
	}

	if err := connect(fd, "127.0.0.1:8081"); err != nil {
		log.Fatal("connect: ", err.Error())
	}

	syscall.Close(fd)
}

func connect(fd int, addr string) error {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return fmt.Errorf("resolve addr err:%s", err.Error())
	}

	raw := syscall.RawSockaddrInet4{}
	raw.Family = syscall.AF_INET
	p := (*[2]byte)(unsafe.Pointer(&raw.Port))
	p[0] = byte((uint16(taddr.Port)) >> 8)
	p[1] = byte((uint16(taddr.Port)))
	for i, ip := range taddr.IP.To4() {
		raw.Addr[i] = ip
	}

	_, _, e1 := syscall.Syscall(syscall.SYS_CONNECT, uintptr(fd), uintptr(unsafe.Pointer(&raw)), 16)
	switch e1 {
	case 0:
		return nil
	default:
		return e1
	}
}

func ResolveSockaddr(addr string) (syscall.Sockaddr, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, errors.New("non-IPv4 address")
	}

	sa := new(syscall.SockaddrInet4)
	for i := 0; i < net.IPv4len; i++ {
		sa.Addr[i] = ip[i]
	}
	sa.Port, _ = strconv.Atoi(port)

	return sa, nil
}
