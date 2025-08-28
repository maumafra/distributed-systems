package main

import (
	"fmt"
	"net"
	//"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	//serverTime = time.Now()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading: ", err)
			return
		}
		data := buffer[:n]
		fmt.Printf("Received: %s\n", data)

		_, err = conn.Write([]byte("Echo: " + string(data)))
		if err != nil {
			fmt.Println("Error writing: ", err)
			return
		}
	}
}

func startServer(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening: ", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server listening on %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	startServer(":8080")
}
