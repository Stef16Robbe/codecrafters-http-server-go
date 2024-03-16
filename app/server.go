package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	// defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatalln("Error accepting connection: ", err.Error())
	}

	fmt.Printf("new conn from: %v\n", conn.RemoteAddr().String())

	_, err = conn.Read([]byte{})

	if err != nil {
		log.Fatalln("Error reading connection: ", err.Error())
	}

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	if err != nil {
		log.Fatalln("Error accepting connection: ", err.Error())
	}
}
