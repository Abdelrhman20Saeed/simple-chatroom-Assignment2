package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Message struct {
	Sender string
	Text   string
	When   time.Time
}

var (
	clients = make(map[net.Conn]string)
	history []Message
	mutex   sync.Mutex
)

func broadcast(sender net.Conn, msg Message) {
	mutex.Lock()
	history = append(history, msg)
	for client := range clients {
		if client != sender {
			fmt.Fprintf(client, "[%s] %s: %s\n", msg.When.Format("15:04:05"), msg.Sender, msg.Text)
		}
	}
	mutex.Unlock()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	
	conn.Write([]byte("Enter your name: "))
	nameRaw, _ := reader.ReadString('\n')
	name := strings.TrimSpace(nameRaw)
	if name == "" {
		name = "Anonymous"
	}

	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	fmt.Printf("%s joined the chat\n", name)
	conn.Write([]byte("--- Chat history ---\n"))
	mutex.Lock()
	for _, m := range history {
		fmt.Fprintf(conn, "[%s] %s: %s\n", m.When.Format("15:04:05"), m.Sender, m.Text)
	}
	mutex.Unlock()
	conn.Write([]byte("--------------------\n"))

	for {
		textRaw, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%s left the chat\n", name)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			return
		}

		text := strings.TrimSpace(textRaw)
		if text == "" {
			continue
		}
		if text == "exit" {
			fmt.Fprintf(conn, "Goodbye, %s!\n", name)
			fmt.Printf("%s left the chat\n", name)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			return
		}

		msg := Message{Sender: name, Text: text, When: time.Now()}
		broadcast(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("✅ Server started on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
