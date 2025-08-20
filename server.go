package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type message struct {
	sender net.Conn
	text   string
}

var (
	clients  = make(map[net.Conn]string)
	messages = make(chan message)
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("error")
	}
	go broadcast()
	fmt.Println("listening...")
	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("connection")
	_, ok := clients[conn]
	if !ok {
		clients[conn] = "User: " + conn.RemoteAddr().String()
	}
	messages <- message{conn, clients[conn] + " has connected to the server."}
	for {
		incoming, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		if incoming[0] == '/' {
			incoming = action(incoming[1:], conn)
			fmt.Fprintf(conn, "%s\n", incoming)
		} else {
			messages <- message{conn, incoming}
		}
	}
}

func action(command string, conn net.Conn) string {
	command = strings.Trim(command, "\n\r")
	words := strings.Split(command, " ")
	switch words[0] {
	case "help":
		return help(words)
	case "nick":
		if len(words) != 2 {
			return "<ADMIN>INVALID NICKNAME PARAMENTERS. PLEASE USE ONE WORD, NO SPACES"
		}
		clients[conn] = words[1]
		return "<ADMIN>NICK NAME SET TO " + words[1]
	case "exit":
		conn.Close()
		messages <- message{conn, clients[conn] + " has disconnected"}
		return ""
	default:
		return "Command unknown. Try /help for list of commands."
	}
}

func broadcast() {
	for {
		msg := <-messages
		for conns := range clients {
			if conns != msg.sender {
				fmt.Fprintf(conns, "[%s]: %s\n", clients[msg.sender], msg.text)
			}
		}
	}
}

func help(words []string) string {
	if len(words) == 1 {
		return "List of commands: /help /nick /exit"
	}
	switch words[1] {
	case "nick":
		return "Sets user nickname"
	case "exit":
		return "Closes connection with server"
	case "help":
		return "Gives inoformation about a given command."
	default:
		return "Command unknown. Try /help for command list"
	}
}
