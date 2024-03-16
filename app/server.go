package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type RequestHeaders struct {
	StartLine StartLine
	Host      string
	UserAgent string
}

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

func parseUserAgent(data string) string {
	if len(data) == 0 {
		return ""
	}
	return strings.Split(data, " ")[1]
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	fmt.Println("new conn from: ", conn.RemoteAddr().String())

	headers := make([]byte, 1024)

	_, err := conn.Read(headers)

	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
	}

	var requestHeaders RequestHeaders
	splitted := strings.Split(string(headers), "\r\n")
	requestHeaders.StartLine = parseStartline(splitted[0])
	requestHeaders.UserAgent = parseUserAgent(splitted[2])

	if requestHeaders.StartLine.Path == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Contains(requestHeaders.StartLine.Path, "echo") {
		echo := strings.Split(requestHeaders.StartLine.Path, "/echo/")[1:][0]
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(echo), echo)
		_, err = conn.Write([]byte(res))
	} else if requestHeaders.StartLine.Path == "/user-agent" {
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(requestHeaders.UserAgent), requestHeaders.UserAgent)
		_, err = conn.Write([]byte(res))
	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		log.Fatalln("Error responding to connection: ", err.Error())
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatalln("Failed to bind to port 4221")
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln("Error accepting connection: ", err.Error())
		}

		go handleRequest(conn)
	}

}
