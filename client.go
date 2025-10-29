package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// الاتصال بالسيرفر على المنفذ 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Type your message:")

	// قراءة الرسائل من المستخدم وإرسالها للسيرفر
	for {
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		// استقبال رد السيرفر
		reply := make([]byte, 1024)
		n, err := conn.Read(reply)
		if err != nil {
			fmt.Println("Server disconnected.")
			return
		}
		fmt.Println("Server:", string(reply[:n]))
	}
}
