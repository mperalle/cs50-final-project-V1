package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// Global variables for IP address and port number
var ipAddr string
var port string

// List of connections to keep track of
var connectionList map[net.Conn]struct{} = make(map[net.Conn]struct{})

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
		ipAddr = flag.Arg(0)
		port = "0"
		fmt.Println("Port not specified, assigned to random number")
	case 2:
		//Assign IP Address and port to the firsts command-line argument after flag
		ipAddr = flag.Arg(0)
		port = flag.Arg(1)
	default:
		fmt.Print("usage: ./main [-l] IP_address [port]\n", "  -l: start server listening on specified socket address\n",
			"port: no port = random port\n")
		os.Exit(2)
	}

	// Start server or client
	if listen {
		startServer()
	} else {
		startClient()
	}
}
