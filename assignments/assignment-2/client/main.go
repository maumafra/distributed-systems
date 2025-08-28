package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		fmt.Print("Write: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println("Error writing: ", err)
			return
		}

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading: ", err)
			return
		}

		data := buffer[:n]
		fmt.Printf("Received: %s\n", data)
	}
}

func startClient(port string) {
	conn, err := net.Dial("tcp", "localhost"+port)
	if err != nil {
		fmt.Println("Error connecting to localhost: ", err)
		return
	}
	handleConnection(conn)
}

func main() {
	startClient(":8080")
}
