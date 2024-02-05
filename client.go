package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func startClient() {
	// Connect to the specified socket address
	connection, err := net.Dial("tcp", ipAddr+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	go readInputAndWriteConnection(connection)

	readConnection(connection)

}

func readConnection(connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := connection.Read(buffer)
		message := string(buffer[:n])
		if len(message) > len("/send") && message[:len("/send")] == "/send" {
			words := strings.Fields(message)
			receiveFile(words[1], connection)

		} else {
			fmt.Println("message received:", message)
		}
		if err == io.EOF {
			break
		}
	}
}

func readInputAndWriteConnection(connection net.Conn) {
	// Read input from terminal and send to connection
	scanner := bufio.NewScanner(os.Stdin)
	// Ask for name
	fmt.Println("Welcome to the chat! What's your name? ")

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Fprintln(connection, message)
		if len(message) > len("/send") && message[:len("/send")] == "/send" {
			words := strings.Fields(message)
			sendFile(words[1], connection)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func sendFile(name string, c net.Conn) {
	fileData, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	c.Write(fileData)
	c.Write([]byte("END"))
	fmt.Println("File sent!")

}

func receiveFile(name string, c net.Conn) {
	buffer := make([]byte, 1024)
	var fileData []byte
	for {
		n, err := c.Read(buffer)
		if err == io.EOF {
			fmt.Println("connection closed")
			break
		} else if strings.Contains(string(buffer[:n]), "END") {
			fileData = append(fileData, buffer[:n-3]...)
			fmt.Println("File received from ")
			break
		}
		fileData = append(fileData, buffer...)
	}
	os.WriteFile("Transfered"+name, fileData, 0666)
	fmt.Println("Transfered"+name, "saved!")
}
