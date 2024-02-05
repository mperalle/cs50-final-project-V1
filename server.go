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

func startServer() {
	// Start TCP server listening on specified port
	listener, err := net.Listen("tcp", ipAddr+":"+port)
	//listener, err := net.Listen("tcp", ":"+port)
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
		connectionList[connection] = struct{}{}

		// Launch go routine for handling the connection
		go connectionHandler(connection)

		//test
		fmt.Println(connectionList)
	}
}

func connectionHandler(connection net.Conn) {

	defer connection.Close()

	// Fetch name from client
	scanner := bufio.NewScanner(connection)
	scanner.Scan()
	name := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// Write to all other connections
	for c := range connectionList {
		if c != connection {
			fmt.Fprintln(c, name+" joined the conversation.")
		} else {
			fmt.Fprintln(c, "You joined the conversation")
		}
	}

	// Loop to continuously read and write data to the client
	for scanner.Scan() {
		message := scanner.Text()
		for c := range connectionList {
			if c != connection {
				if len(message) > len("/send") && message[:len("/send")] == "/send" {
					fmt.Fprintln(c, message)
					forwardFile(c, connection)
				}

				fmt.Fprintln(c, name+": "+message)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("reading standard input:", err)
	} else {
		// Write to all other connections
		for c := range connectionList {
			if c != connection {
				fmt.Fprintln(c, name+" left the conversation.")
			}
		}
		// Log when client disconnects
		log.Println("Client disconnected:", connection.RemoteAddr())
	}
}

func forwardFile(writeConnection net.Conn, readConnection net.Conn) {
	buffer := make([]byte, 1024)
	var fileData []byte
	for {
		n, err := readConnection.Read(buffer)
		if err == io.EOF {
			fmt.Println("connection closed")
			break
		} else if strings.Contains(string(buffer[:n]), "END") {
			fileData = append(fileData, buffer[:n-3]...)
			fmt.Println("File reading done!")
			break
		}
		fileData = append(fileData, buffer...)
	}

	writeConnection.Write(fileData)
	writeConnection.Write([]byte("END"))
	fmt.Println("File sent to:", writeConnection.RemoteAddr())

}
