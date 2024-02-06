package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func startClient() {
	// Connect to the specified socket address
	connection, err := net.Dial("tcp", ipAddr+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	go readInputAndWriteConnection(connection)

	readConnection(connection)

}

func readConnection(connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := connection.Read(buffer)
		message := string(buffer[:n])
		if len(message) > len("/send") && message[:len("/send")] == "/send" {
			words := strings.Fields(message)
			fmt.Println("/send command received for:", words[1])
			receiveFile(words[1], connection)

		} else {
			fmt.Print(message)
		}
		if err == io.EOF {
			break
		}
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
		if len(message) > len("/send") && message[:len("/send")] == "/send" {
			words := strings.Fields(message)
			sendFile(words[1], connection)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func sendFile(name string, c net.Conn) {
	fileData, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	fileSizeUint64 := uint64(len(fileData))
	fileSizeByte := make([]byte, 8)

	fmt.Println("File size:", fileSizeUint64)

	binary.LittleEndian.PutUint64(fileSizeByte, fileSizeUint64)
	_, err = c.Write(fileSizeByte)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Write(fileData)
	if err != nil {
		log.Fatal(err)
	}
	// c.Write([]byte("END"))
	fmt.Println("File sent!")

}

func receiveFile(name string, c net.Conn) {
	buffer := make([]byte, 1024)
	var fileData []byte
	fileSizeByte := make([]byte, 8)

	_, err := io.ReadFull(c, fileSizeByte)
	if err != nil {
		log.Fatal(err)
	}

	fileSizeUint64 := binary.LittleEndian.Uint64(fileSizeByte)
	fmt.Println("Got file size! :", fileSizeUint64)

	reader := io.LimitReader(c, int64(fileSizeUint64))

	fmt.Println("Reading incoming file...")
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
	os.WriteFile("Transfered"+name, fileData, 0666)
	fmt.Println("Transfered"+name, "saved!")
}
