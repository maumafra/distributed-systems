package main

import (
	"fmt"
	"bufio"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to localhost: ", err.Error())
		return
	}
	defer conn.Close()
	fmt.Print("Enter a number: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	_, err = conn.Write([]byte(input))
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		return
	}

	response := string(buffer[:n])
	fmt.Println("Server response: ", response)

	_, err = strconv.Atoi(response)
	if err != nil {
		fmt.Println("Error")
	}

}