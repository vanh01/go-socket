package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

func main() {
	connection, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	fmt.Printf("connected to %s\n", connection.RemoteAddr().String())
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("text to send: ")
		input := make([]byte, 1024)
		ilen, err := reader.Read(input)
		if err != nil {
			panic(err)
		}
		input = input[:ilen]
		s := string(input)
		s = strings.Trim(s, "\n")
		_, err = connection.Write([]byte(s))
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(connection)
		reader.Peek(1)
		response := readData(reader)
		if response == "#closed" {
			fmt.Printf("connection to %s is closed!\n", connection.RemoteAddr().String())
			break
		}
		fmt.Printf("sever reply: %s\n", response)
	}
}
