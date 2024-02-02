package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const IPAddr string = "localhost:3000"

func main() {

	//Connect to the specified IP address
	connection, err := net.Dial("tcp", IPAddr)
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
}
