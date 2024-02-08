# cs50-final-project-V1 : EasyChatApp

EasyChatApp is a command-line application for instant messaging and file transfering over TCP. 

## Description

EasyChatApp allows anyone to start a TCP server and have a chat room in the terminal accessible by anyone on the same local network. It allows also anyone in the conversation to send files to everyone.

## Getting Started

### Dependencies

* No dependencies, the application only uses packages from the standard library 

### How to use the application

* Download the executable file "EasyChatApp" from this repository

* To start the TCP server use the -l flag and specify the IP address and port number
```
./EasyChatApp -l [IP_address] [port]
```

* To start a TCP client specify the IP address and port number of the server
```
./EasyChatApp [IP_address] [port]
```

* Enter your name and join the chat room

* To send a file to the chat room (all connected clients will receive it), use this command:
```
/send name_of_file
```
