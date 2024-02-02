package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// List of connections to keep track of
// var ConnectionList []net.Conn
var ConnectionList map[net.Conn]struct{} = make(map[net.Conn]struct{})

func main() {

	var IPAddr string
	var port string

	// Retrieve IP address and port from command-line arguments
	switch len(os.Args) {
	case 2:
		IPAddr = os.Args[1]
		port = "0"
		fmt.Println("Random port assigned")
	case 3:
		IPAddr = os.Args[1]
		port = os.Args[2]
	default:
		fmt.Println("main.go [IP address] [port]")
		os.Exit(2)
	}

	// Start TCP server listening on specified port
	listener, err := net.Listen("tcp", IPAddr+":"+port)
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

	// Fetch name from client
	scanner := bufio.NewScanner(connection)
	scanner.Scan()
	name := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// Write to all other connections
	for c := range ConnectionList {
		if c != connection {
			fmt.Fprintln(c, name+" joined the conversation.")
		}
	}

	// Loop to continuously read and write data to the client
	for scanner.Scan() {
		for c := range ConnectionList {
			if c != connection {
				fmt.Fprintln(c, name+": "+scanner.Text())
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
