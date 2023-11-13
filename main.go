package main

import (
	"bufio"
	"fmt"
	"net"
)

func readData(reader *bufio.Reader) string {
	ss := ""
	for {
		if reader.Buffered() == 0 {
			break
		}
		receive := make([]byte, 5)
		rlen, err := reader.Read(receive)
		if err != nil {
			return ""
		}
		str := string(receive[:rlen])
		ss += str
	}
	return ss
}

func receiveData(conn net.Conn) {
	fmt.Printf("connection from %s is connected\n", conn.RemoteAddr().String())
	for {
		reader := bufio.NewReader(conn)
		reader.Peek(1)
		s := readData(reader)
		if s == "" {
			fmt.Printf("connection from %s is closed!\n", conn.RemoteAddr().String())
			conn.Close()
			break
		}
		fmt.Printf("data received: %s\n", s)
		conn.Write([]byte(fmt.Sprintf("received: %s", s)))
	}
}

func main() {
	server, err := net.Listen("tcp", ":1234")
	fmt.Println("server is started")
	defer server.Close()
	if err != nil {
		panic(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go receiveData(conn)
	}
}
