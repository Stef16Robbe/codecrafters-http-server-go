package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type StartLine struct {
	Method  Method
	Path    string
	Version string
}

type Method int

const (
	GET Method = iota
	POST
	PUT
	PATCH
)

func parseStartline(data string) StartLine {
	sep := strings.Split(data, " ")
	method, err := MethodString(sep[0])
	if err != nil {
		log.Fatalln("Error parsing HTTP method:", err)
	}

	return StartLine{
		Method:  method,
		Path:    sep[1],
		Version: sep[2],
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatalln("Failed to bind to port 4221")
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatalln("Error accepting connection: ", err.Error())
	}
	defer conn.Close()

	fmt.Println("new conn from: ", conn.RemoteAddr().String())

	headers := make([]byte, 1024)

	_, err = conn.Read(headers)

	if err != nil {
		log.Fatalln("Error reading connection: ", err.Error())
	}

	splitted := strings.Split(string(headers), "\r\n")
	startLinestr := splitted[0]
	startLine := parseStartline(startLinestr)

	if startLine.Path == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Contains(startLine.Path, "echo") {
		echo := strings.Split(startLine.Path, "/echo/")[1:][0]
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(echo), echo)
		_, err = conn.Write([]byte(res))
	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		log.Fatalln("Error responding to connection: ", err.Error())
	}
}
