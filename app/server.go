package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {

	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Read ", n, " bytes")
	fmt.Println("get request: ", string(req))

	reqURL := getRequestURL(req)
	if reqURL == "/" {
		fmt.Println("Request URL: /, sending 200 OK")
		n, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		fmt.Println("Request URL: ", reqURL, " sending 404 Not Found")
		n, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Wrote ", n, " bytes")
}

func getRequestURL(req []byte) string {

	// request line \r\n headers \r\n request body \r\n
	requestLine := bytes.Split(req, []byte("\r\n"))[0]
	reqURLByte := bytes.Split(requestLine, []byte(" "))[1]
	return string(reqURLByte)
}
