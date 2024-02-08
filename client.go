package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func startClient() {
	// Connect to the TCP server at the specified address
	connection, err := net.Dial("tcp", ipAddr+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	// Close connection
	defer connection.Close()

	// Start goroutine to read from terminal input and write to the server
	go readStdinWriteConnection(connection)

	// Read data from the server
	readConnection(connection)

}

func readConnection(connection net.Conn) {

	// Create buffer to read data
	buffer := make([]byte, 1024)

	// Continuously read data from the server connection
	for {
		n, err := connection.Read(buffer)

		// Handle error in reading
		if err != nil {
			if err == io.EOF {
				// Connection closed
				break
			} else {
				log.Fatal(err)
			}
		}

		// Check if received message is a file transfer command
		message := string(buffer[:n])
		if len(message) > len("/send") && message[:len("/send")] == "/send" {

			// Fetch the file name
			words := strings.Fields(message)
			fmt.Println("Incoming file transfer...")

			// Read file from the server
			_, fileData := readFile(connection)

			// Save file
			os.WriteFile(words[1], fileData, 0666)
			fmt.Println(words[1], "saved!")

		} else {

			// Print received message
			fmt.Print(message)
		}

	}
}

func readStdinWriteConnection(connection net.Conn) {

	// Create a scanner to read from terminal input
	scanner := bufio.NewScanner(os.Stdin)

	// Ask for a name
	fmt.Println("Welcome to the chat! What's your name? ")

	// Continuously reading from terminal input
	for scanner.Scan() {

		message := scanner.Text()

		// Check if it is a send file command
		if len(message) > len("/send") && message[:len("/send")] == "/send" {

			// Fetch the file name
			words := strings.Fields(message)

			// Read file
			fileData, err := os.ReadFile(words[1])
			if err != nil {
				fmt.Println("Error in opening the file, please check if the file name is correct")
				continue
			}

			// Write send command to the server
			fmt.Fprintln(connection, message)

			// Delay to make sure server has time to start reading before sending file
			time.Sleep(100 * time.Millisecond)

			// Send file to the server
			sendFile(fileData, connection)

		} else {

			// Write message to the server connection
			fmt.Fprintln(connection, message)
		}
	}

	// Handle error for the reading
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
