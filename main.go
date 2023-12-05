package main

import (
	"bufio"
	"fmt"
	"net"
)

func authentication(reader *bufio.Reader) (bool, error) {
	s, err := readData(reader)
	if err != nil {
		return false, err
	}
	token := (*s)[1:]
	// logic for auth
	if len(token) > 5 {
		return true, nil
	}
	return false, nil
}

// this function read data per 5 bytes
func readData(reader *bufio.Reader) (*string, error) {
	ss := ""
	for {
		if reader.Buffered() == 0 {
			break
		}
		receive := make([]byte, 5)
		rlen, err := reader.Read(receive)
		if err != nil {
			return nil, err
		}
		str := string(receive[:rlen])
		ss += str
	}
	return &ss, nil
}

func receiveData(conn *net.Conn) {
	fmt.Printf("connection from %s is connected\n", (*conn).RemoteAddr().String())
	defer (*conn).Close()
	// this is authentication part
	reader := bufio.NewReader((*conn))
	reader.Peek(1)
	ok, err := authentication(reader)
	if err != nil {
		fmt.Printf("connection from %s is closed!\n", (*conn).RemoteAddr().String())
		return
	}
	if !ok {
		fmt.Println("Permission denied!")
		(*conn).Write([]byte("#closed"))
		return
	}
	(*conn).Write([]byte("Success!"))
	// pass authentication part
	for {
		reader := bufio.NewReader((*conn))
		reader.Peek(1)
		s, err := readData(reader)
		if err != nil || *s == "" {
			fmt.Printf("connection from %s is closed!\n", (*conn).RemoteAddr().String())
			break
		}
		fmt.Printf("data received: %s\n", *s)
		(*conn).Write([]byte(fmt.Sprintf("received: %s", *s)))
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
		go receiveData(&conn)
	}
}
