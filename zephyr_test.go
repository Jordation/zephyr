package zephyr

import (
	"net"
	"testing"
)

func TestListen(t *testing.T) {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic("listen err: " + err.Error())
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		panic("accept err: " + err.Error())
	}

	buff := []byte("hello world!")
	_, err = conn.Write(buff)
	if err != nil {
		panic("write errs: " + err.Error())
	}
}

func Test_Zephman(t *testing.T) {
	zeph := New()
	zeph.Run(":3000")
}
