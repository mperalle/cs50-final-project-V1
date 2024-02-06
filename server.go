package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func startServer() {
	// Start TCP server listening on specified port
	//listener, err := net.Listen("tcp", ipAddr+":"+port)
	listener, err := net.Listen("tcp", ":"+port)
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
		if len(message) > len("/send") && message[:len("/send")] == "/send" {
			fmt.Println("/send command received...and transfering to all connections...")
			for c := range connectionList {
				if c != connection {
					fmt.Println("/send command transfering to:", c.RemoteAddr())
					fmt.Fprintln(c, message)
				}
			}
			forwardFile(connection)

		} else {
			for c := range connectionList {
				if c != connection {

					fmt.Fprintln(c, name+": "+message)
				}
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

func forwardFile(connection net.Conn) {
	buffer := make([]byte, 1024)
	var fileData []byte
	fileSizeByte := make([]byte, 8)

	_, err := io.ReadFull(connection, fileSizeByte)
	if err != nil {
		log.Fatal(err)
	}

	fileSizeUint64 := binary.LittleEndian.Uint64(fileSizeByte)

	fmt.Println("Got file size! :", fileSizeUint64)
	reader := io.LimitReader(connection, int64(fileSizeUint64))

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			} else {
				fmt.Println("File reading complete!")
				break
			}
		}
		dataRead := buffer[:n]
		fileData = append(fileData, dataRead...)
	}

	for c := range connectionList {
		if c != connection {
			fmt.Println("Write fileSizeByte...")
			_, err = c.Write(fileSizeByte)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Write fileData...")
			_, err = c.Write(fileData)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("File sent to:", c.RemoteAddr())
		}
	}
}
