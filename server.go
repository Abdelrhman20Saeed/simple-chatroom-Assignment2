package main

import (
	"fmt"
	"net"
)

func main() {
	// استمع على المنفذ 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Client disconnected.")
			return
		}
		message := string(buffer[:n])
		fmt.Println("Received:", message)
		conn.Write([]byte("Message received\n"))
	}
}
