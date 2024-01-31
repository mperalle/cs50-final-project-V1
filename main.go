package main

import (
	"io"
	"log"
	"net"
)

const ListAddr string = "localhost:3000"

func main() {

	//start TCP server listening on specified port
	listener, err := net.Listen("tcp", ListAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	//loop for accepting incoming connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		//test with an echo response
		io.Copy(connection, connection)
	}
}
