package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var IPAddr string
var port string

// List of connections to keep track of
// var ConnectionList []net.Conn
var ConnectionList map[net.Conn]struct{} = make(map[net.Conn]struct{})

func main() {

	// Declare flag for starting the server
	var listen bool
	flag.BoolVar(&listen, "l", false, "start TCP server")
	flag.Parse()
	fmt.Println(flag.Args())

	// Retrieve IP address and port from command-line arguments
	switch len(flag.Args()) {
	case 1:
		//Assign IP Address to the first command-line argument after flag
		IPAddr = flag.Arg(0)
		port = "0"
		fmt.Println("Port not specified, assigned to random number")
	case 2:
		//Assign IP Address and port to the firsts command-line argument after flag
		IPAddr = flag.Arg(0)
		port = flag.Arg(1)
	default:
		fmt.Print("usage: ./main [-l] IP_address [port]\n", "  -l: start server listening on specified socket address\n",
			"port: no port = random port\n")
		os.Exit(2)
	}

	if listen {
		startServer()
	} else {
		startClient()
	}

}

func startServer() {
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
		} else {
			fmt.Fprintln(c, "You joined the conversation")
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
		fmt.Println("reading standard input:", err)
	} else {
		// Write to all other connections
		for c := range ConnectionList {
			if c != connection {
				fmt.Fprintln(c, name+" left the conversation.")
			}
		}
	}

}

func startClient() {
	// Connect to the specified socket address
	connection, err := net.Dial("tcp", IPAddr+":"+port)
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
