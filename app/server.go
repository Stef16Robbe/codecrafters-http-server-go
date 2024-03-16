package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
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

func checkFileExists(dir, filename string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
		// log.Fatalln("Folder does not exist.")
	}

	// assuming dir ends in /
	_, err = os.Stat(dir + filename)
	// fmt.Println("exists or not:", os.IsNotExist(err), dir+filename)
	return !os.IsNotExist(err)
}

func readFile(dir, filename string) string {
	content, err := os.ReadFile(dir + filename)
	if err != nil {
		log.Fatalln("Err:", err.Error())
	}

	return string(content)
}

func handleRequest(conn net.Conn, dir string) {
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
	} else if strings.Contains(requestHeaders.StartLine.Path, "/files/") {
		if dir == "" {
			log.Fatalf("Give up dir!")
		}
		// assuming just file name not a path
		filename := strings.Split(requestHeaders.StartLine.Path, "/files/")[1]
		if filename == "" {
			log.Fatalf("Incorrect filename!")
		}
		if checkFileExists(dir, filename) {
			filecontent := readFile(dir, filename)
			res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n%v", len(filecontent), filecontent)
			_, err = conn.Write([]byte(res))
		} else {
			_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
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

	dir := flag.String("directory", "", "help message for flag n")
	flag.Parse()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln("Error accepting connection: ", err.Error())
		}

		go handleRequest(conn, *dir)
	}

}
