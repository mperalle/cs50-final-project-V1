package main

import (
	"fmt"
	"log"
	"net"
)

const ListAddr string = "localhost:3000"

// List of connections to keep track of
var ConnectionList []net.Conn

func main() {

	//start TCP server listening on specified port
	listener, err := net.Listen("tcp", ListAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Start server on: ", listener.Addr().String())

	defer listener.Close()

	//loop for accepting incoming connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("New connection from: ", connection.RemoteAddr())
		//add connection to list of connections
		ConnectionList = append(ConnectionList, connection)

		//launch go routine for handling the connection
		go connectionHandler(connection)
	}
}

func connectionHandler(connection net.Conn) {
	//create buffer to read from connection
	buffer := make([]byte, 1024)
	numberRead, err := connection.Read(buffer)
	if err != nil {
		log.Println(err)
	}
	message := buffer[0:numberRead]

	//write to all other connections
	for _, c := range ConnectionList {
		if c != connection {
			c.Write(message)
		}
	}
}
