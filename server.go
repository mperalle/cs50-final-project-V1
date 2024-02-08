package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func startServer() {
	// Start TCP server listening for incoming connections on specified address
	listener, err := net.Listen("tcp", ipAddr+":"+port)

	if err != nil {
		log.Println("IP address or port invalid")
		log.Fatal(err)
	}
	fmt.Println("Server listening on: ", listener.Addr().String())

	defer listener.Close()

	for {
		// Accept incoming connections
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("New connection from: ", connection.RemoteAddr())

		// Add connection to the list of connections
		connectionList[connection] = struct{}{}

		// Start a goroutine to handle the connection
		go connectionHandler(connection)

	}
}

func connectionHandler(connection net.Conn) {

	defer connection.Close()

	// Fetch name from client
	scanner := bufio.NewScanner(connection)
	scanner.Scan()
	name := scanner.Text()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Send notification when the user joins the conversation
	for c := range connectionList {
		if c != connection {
			fmt.Fprintln(c, name+" joined the conversation.")
		} else {
			fmt.Fprintln(c, "You joined the conversation")
		}
	}

	// Continuously read data from the connection
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println(message)

		// Check if command for file transfer received
		if len(message) > len("/send") && message[:len("/send")] == "/send" {

			// Forward the send command to the other connections
			for c := range connectionList {
				if c != connection {
					fmt.Fprintln(c, message)

				}
			}

			fileSizeByte, fileData := readFile(connection)

			// Forward the file data to the other connections
			for c := range connectionList {
				if c != connection {
					_, err := c.Write(fileSizeByte)
					if err != nil {
						log.Fatal(err)
					}
					_, err = c.Write(fileData)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println("File sent to:", c.RemoteAddr())
				}
			}

		} else {
			// Forward message to the other connections
			for c := range connectionList {
				if c != connection {
					fmt.Fprintln(c, name+": "+message)
				}
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)

	} else {
		// Send notification to the other connections when leaving
		for c := range connectionList {
			if c != connection {
				fmt.Fprintln(c, name+" left the conversation.")
			}
		}
		// Log when client disconnects
		log.Println("Client disconnected:", connection.RemoteAddr())

		// Remove connection from connection list
		delete(connectionList, connection)
	}
}
