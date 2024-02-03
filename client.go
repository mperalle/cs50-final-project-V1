package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println(message)
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
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
