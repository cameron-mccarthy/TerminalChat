package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("error")
	}
	go recieve(conn)
	var message string
	fmt.Print("You have connected.\n<you>")
	for {
		message, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(conn, message)
		fmt.Print("<you>")
	}
}

func recieve(conn net.Conn) {
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("\r\033[K")
		if err != nil {
			fmt.Println("Server disconnected.")
			break
		}
		fmt.Println(strings.Trim(msg, "\n\r"))
		fmt.Print("<you>")
	}
}
