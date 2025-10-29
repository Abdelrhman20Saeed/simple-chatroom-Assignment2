package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func connectToServer() net.Conn {
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			log.Println("⚠️ Can't connect to server, retrying in 2s...")
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println("✅ Connected to server!")
		return conn
	}
}

func main() {
	conn := connectToServer()
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	
	go func() {
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				log.Println("❌ Disconnected from server, reconnecting...")
				conn = connectToServer()
				serverReader = bufio.NewReader(conn)
				continue
			}
			fmt.Print(message)
		}
	}()

	
	for {
		textRaw, _ := reader.ReadString('\n')
		text := strings.TrimSpace(textRaw)
		if text == "" {
			continue
		}
		if text == "exit" {
			fmt.Println("👋 Exiting chat...")
			conn.Write([]byte(text + "\n"))
			return
		}
		conn.Write([]byte(text + "\n"))
	}
}
