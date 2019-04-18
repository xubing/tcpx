package main

import (
	"bytes"
	"encoding/json"
	"errorX"
	"fmt"
	"net"
	"tcpx"
)

var packx = tcpx.NewPackx(tcpx.JsonMarshaller{})

func main() {
	conn, err := net.Dial("tcp", "localhost:7171")
	if err != nil {
		panic(err)
	}
	received := Receive(conn)
	go func() {
		for {
			buf := <-received
			var message tcpx.Message
			var receivedString string
			fmt.Println(buf)
			message, e := packx.Unpack(buf, &receivedString)
			if e != nil {
				panic(errorx.Wrap(e))
			}
			fmt.Println("收到服务端消息块:", smartPrint(message))
			fmt.Println("服务端消息:", receivedString)
		}
	}()
	buf, e := packx.Pack(1, "hello,I am client xiao ming")
	if e != nil {
		panic(e)
	}
	conn.Write(buf)
	conn.Write(buf)
	conn.Write(buf)
	conn.Write(buf)
	select {}
}

func Receive(conn net.Conn) <-chan []byte {
	var info = make([]byte, 16, 16)
	var content []byte
	var received = make(chan []byte, 200)
	go func() {
		for {
			conn.Read(info)
			length, e := packx.LengthOf(info)
			if e != nil {
				panic(e)
			}
			content = make([]byte, length)
			conn.Read(content)

			received <- append(info, content ...)
		}
	}()
	return received
}

// fucking buffer
func Receive2(conn net.Conn) <-chan []byte {
	var buffer = bytes.NewBuffer(nil)
	var received = make(chan []byte, 0)
	go func() {
		for {
			buffer.Reset()
			buffer.ReadFrom(conn)
			received <- buffer.Bytes()
		}
	}()
	return received
}

func smartPrint(src interface{}) string {
	buf, _ := json.MarshalIndent(src, "  ", "  ")
	return string(buf)
}
