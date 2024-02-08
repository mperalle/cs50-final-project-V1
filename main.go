package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Declare global variables for IP address and port number
var ipAddr string
var port string

// Initialyze the list of connections to keep track of
var connectionList map[net.Conn]struct{} = make(map[net.Conn]struct{})

func main() {

	// Declare flag to start the server
	var listen bool
	flag.BoolVar(&listen, "l", false, "start TCP server")
	flag.Parse()

	//Retrieve IP address and port from the command-line arguments
	switch len(flag.Args()) {
	case 1:
		//Assign IP Address to the first command-line argument after flag
		if net.ParseIP(flag.Arg(0)) != nil {
			ipAddr = flag.Arg(0)
			port = ""
		} else {
			ipAddr = ""
			port = flag.Arg(0)
		}
	case 2:

		//Assign IP Address and port to the firsts command-line argument after flag
		ipAddr = flag.Arg(0)
		port = flag.Arg(1)
	default:
		fmt.Println("usage: ./main [-l] [IP_address] [port]")
		os.Exit(2)
	}

	// Check the flag to start the server or the client
	if listen {
		startServer()
	} else {
		startClient()
	}
}

func sendFile(fileData []byte, c net.Conn) {

	// Get and convert the file size into a slice of 8 bytes
	fileSizeUint64 := uint64(len(fileData))
	fileSizeByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(fileSizeByte, fileSizeUint64)

	// Send the file size to the connection
	_, err := c.Write(fileSizeByte)
	if err != nil {
		log.Fatal(err)
	}

	// Send the file to the connection
	_, err = c.Write(fileData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File sent!")

}

func readFile(c net.Conn) ([]byte, []byte) {

	buffer := make([]byte, 1024)
	var fileData []byte
	fileSizeByte := make([]byte, 8)

	// Read the next 8 bytes corresponding to the file size
	_, err := io.ReadFull(c, fileSizeByte)
	if err != nil {
		log.Fatal(err)
	}

	// Create a reader to read the exact file bytes
	fileSizeUint64 := binary.LittleEndian.Uint64(fileSizeByte)
	reader := io.LimitReader(c, int64(fileSizeUint64))

	// Read incoming file data into fileData variable
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		dataRead := buffer[:n]
		fileData = append(fileData, dataRead...)
	}

	return fileSizeByte, fileData
}
