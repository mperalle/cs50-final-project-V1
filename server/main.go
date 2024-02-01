package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const ListAddr string = "localhost:3000"

// List of connections to keep track of
// var ConnectionList []net.Conn
var ConnectionList map[net.Conn]struct{} = make(map[net.Conn]struct{})

func main() {

	// Start TCP server listening on specified port
	listener, err := net.Listen("tcp", ListAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server listening on: ", listener.Addr().String())

	defer listener.Close()

	// Loop for accepting incoming connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("New connection from: ", connection.RemoteAddr())
		// Add connection to list of connections
		ConnectionList[connection] = struct{}{}

		// Launch go routine for handling the connection
		go connectionHandler(connection)

		//test
		fmt.Println(ConnectionList)
	}
}

func connectionHandler(connection net.Conn) {

	defer connection.Close()

	// Ask for name
	var name string
	fmt.Fprintln(connection, "Welcome to the chat! What's your name?")
	fmt.Fscan(connection, &name)

	// Write to all connections
	for c := range ConnectionList {
		fmt.Fprintln(c, name, "joined the conversation.")
	}

	// Create buffer to read from connection
	buffer := make([]byte, 1024)

	// Loop to continuously read and write data
	for {
		// Waiting for data to be read in the connection
		numberRead, err := connection.Read(buffer)
		// Error handling
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected: ", connection.RemoteAddr())
			} else {
				log.Println("Error in reading: ", err)
			}
			// Delete connection from ConnectionList
			delete(ConnectionList, connection)

			// Write to all other connections
			for c := range ConnectionList {
				fmt.Fprintln(c, name, "left the conversation.")
			}

			return
		}

		message := string(buffer[0:numberRead])

		// Write to all other connections
		for c := range ConnectionList {
			if c != connection {
				fmt.Fprint(c, name, ": ", message)
			}
		}
	}
}
