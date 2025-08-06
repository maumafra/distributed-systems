package main

import (
	"fmt"
	"net"
	"strconv"
)

func factorial(n int) int {
	if n < 0 {
		return 0
	}

	if n == 0 || n == 1 {
		return 1
	}

	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	return res
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		return
	}
	defer ln.Close()
	fmt.Println("Server listening on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		return
	}

	numStr := string(buffer[:n])
	num, err := strconv.Atoi(numStr)
	if err != nil {
		fmt.Println("Error converting: ", err.Error())
		conn.Write([]byte("Invalid number"))
		return
	}

	result := factorial(num)
	response := strconv.Itoa(result)
	conn.Write([]byte(response))
}