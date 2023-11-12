package main

import (
	"bufio"
	"fmt"
	"net"
)

func readDataFromReader(reader *bufio.Reader) (chan string, chan int) {
	r := make(chan int)
	s := make(chan string)
	go func() {
		ss := ""
		for {
			if reader.Buffered() == 0 {
				r <- 0
				s <- ss
				return
			}
			receive := make([]byte, 5)
			rlen, err := reader.Read(receive)
			if err != nil {
				r <- 1
				s <- ss
				return
			}
			str := string(receive[:rlen])
			ss += str
		}
	}()
	return s, r
}

func receiveData(conn net.Conn) {
	fmt.Printf("connection from %s is connected\n", conn.RemoteAddr().String())
	for {
		reader := bufio.NewReader(conn)
		reader.Peek(1)
		s, i := readDataFromReader(reader)
		err := <-i
		ss := <-s
		if err == 1 || ss == "" {
			fmt.Printf("connection from %s is closed!\n", conn.RemoteAddr().String())
			conn.Close()
			break
		}
		fmt.Printf("data received: %s\n", ss)
		conn.Write([]byte(fmt.Sprintf("received: %s", ss)))
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
