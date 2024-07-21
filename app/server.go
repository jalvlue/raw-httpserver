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
	defer conn.Close()

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("get request, read ", n, " bytes")
	fmt.Println(string(req))

	reqURL := getRequestURL(req)

	n, err = routeRequest(conn, reqURL)
	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("send response back, wrote ", n, " bytes")
}

func routeRequest(conn net.Conn, reqURL string) (n int, err error) {
	if isRootEndpoint(reqURL) {
		fmt.Println("Request URL: /, sending 200 OK")
		n, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if isEchoEndpoint(reqURL) {
		fmt.Println("Request URL: /echo, sending 200 OK, and echoing back msg")

		msg := getEchoMsg(reqURL)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg)
		n, err = conn.Write(buf.Bytes())
	} else {
		fmt.Println("Request URL: ", reqURL, " sending 404 Not Found")
		n, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	return
}

func isRootEndpoint(url string) bool {
	return url == "/"
}

func isEchoEndpoint(url string) bool {
	urlComponents := strings.Split(url, "/")
	return urlComponents[1] == "echo"
}

func getEchoMsg(url string) string {
	return url[6:]
}

func getRequestURL(req []byte) string {
	// request line \r\n headers \r\n request body \r\n
	requestLine := bytes.Split(req, []byte("\r\n"))[0]
	reqURLByte := bytes.Split(requestLine, []byte(" "))[1]
	fmt.Println("Request URL: ", string(reqURLByte), " len: ", len(reqURLByte))
	return string(reqURLByte)
}
