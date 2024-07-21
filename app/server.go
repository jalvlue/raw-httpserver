package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
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

	handleConnection(conn)
}

// handleConnection reads the request from the connection,
// routes the request, and writes the response back to the connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		return
	}
	fmt.Println("Recieve request, read ", n, " bytes")
	fmt.Println(string(req))

	reqURL := getRequestURL(req)

	n, err = routeRequest(conn, reqURL)
	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		return
	}
	fmt.Println("Send response back, wrote ", n, " bytes")
}

// routeRequest routes the request to the corresponding handler based on the request URL
func routeRequest(conn net.Conn, reqURL string) (int, error) {
	var response []byte

	if isRootEndpoint(reqURL) {
		fmt.Println("Request URL: /, sending 200 OK")
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if isEchoEndpoint(reqURL) {
		fmt.Println("Request URL: /echo, sending 200 OK, and echoing back msg")

		msg := getEchoMsg(reqURL)
		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg))
	} else {
		fmt.Println("Request URL: ", reqURL, " sending 404 Not Found")
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	return conn.Write(response)
}

// isRootEndpoint checks if the request URL is the root endpoint
func isRootEndpoint(url string) bool {
	return url == "/"
}

// isEchoEndpoint checks if the request URL is the echo endpoint
func isEchoEndpoint(url string) bool {
	return strings.HasPrefix(url, "/echo/")
}

// getEchoMsg returns the message to be echoed back for the echo endpoint
func getEchoMsg(url string) string {
	prefixLen := len("/echo/")
	return url[prefixLen:]
}

// getRequestURL parses the request URL from the request
func getRequestURL(req []byte) string {
	// request line \r\n headers \r\n request body \r\n
	requestLine := bytes.Split(req, []byte("\r\n"))[0]
	reqURLByte := bytes.Split(requestLine, []byte(" "))[1]
	fmt.Println("Request URL: ", string(reqURLByte), " len: ", len(reqURLByte))
	return string(reqURLByte)
}
