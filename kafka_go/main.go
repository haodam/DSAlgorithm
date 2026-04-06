package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// startServer()
	fmt.Println(os.Args)
	if os.Args[1] == "server" {
		startServer()
	} else {
		clientConnect()
	}
}

func startServer() {
	ln, _ := net.Listen("tcp", ":1234")
	conn, _ := ln.Accept() // Block until can
	conn.Close()
}

func clientConnect() {
	conn, _ := net.Dial("tcp", ":1234")

	// Read from stdin in a line
	rd := bufio.NewReader(os.Stdin)
	line, err := rd.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Printf("Send to server: %s\n", line)

	// Write to server
	streamWr := bufio.NewWriter(conn)
	streamWr.WriteByte(byte(len(line)))
	streamWr.WriteString(line)
	streamWr.Flush()

	conn.Close()
}
