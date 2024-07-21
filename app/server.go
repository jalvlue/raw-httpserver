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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
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

	n, err = routeRequest(conn, req)
	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		return
	}
	fmt.Println("Send response back, wrote ", n, " bytes")
}

// routeRequest routes the request to the corresponding handler based on the request URL
func routeRequest(conn net.Conn, req []byte) (int, error) {
	var response []byte

	reqURL := getRequestURL(req)
	if isRootEndpoint(reqURL) {
		fmt.Println("Request URL: /, sending 200 OK")
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if isEchoEndpoint(reqURL) {
		fmt.Println("Request URL: /echo, sending 200 OK, and echoing back msg")

		msg := getEchoMsg(reqURL)
		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(msg), msg))
	} else if isUserAgentEndpoint(reqURL) {
		userAgent := getHeader(req, []byte("User-Agent"))
		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent))
	} else if isFileEndPoint(reqURL) {
		filename := getRequestFilename(reqURL)
		content, err := readFileContent(filename)
		if err != nil {
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		} else {
			response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), string(content)))
		}
	} else {
		fmt.Println("Request URL: ", reqURL, " sending 404 Not Found")
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	return conn.Write(response)
}

// isRootEndpoint checks if the request URL is to the root endpoint
func isRootEndpoint(url []byte) bool {
	return bytes.Equal(url, []byte("/"))
}

// isEchoEndpoint checks if the request URL is to the echo endpoint
func isEchoEndpoint(url []byte) bool {
	return bytes.HasPrefix(url, []byte("/echo/"))
}

// isUserAgentEndpoint checks if the request URL is to the user-agent endpoint
func isUserAgentEndpoint(url []byte) bool {
	return bytes.HasPrefix(url, []byte("/user-agent"))
}

// isFileEndPoint checks if the request URL is to the file endpoint
func isFileEndPoint(url []byte) bool {
	return bytes.HasPrefix(url, []byte("/files"))
}

// getEchoMsg returns the message to be echoed back for the echo endpoint
func getEchoMsg(url []byte) string {
	prefixLen := len("/echo/")
	return string(url[prefixLen:])
}

// getHeader returns the value of the header with the given key
func getHeader(req []byte, headerKey []byte) string {
	headerKey = append(headerKey, []byte(": ")...)
	headers := bytes.Split(req, []byte("\r\n"))
	// Skip the request line and the request body
	for i := 1; i < len(headers) && !bytes.Equal(headers[i], []byte("")); i++ {
		hk := headers[i]
		if bytes.HasPrefix(hk, headerKey) {
			return string(bytes.Split(hk, []byte(": "))[1])
		}
	}

	return ""
}

// getRequestFilename returns the filename from the request URL
func getRequestFilename(url []byte) string {
	prefixLen := len("/files/")
	return string(url[prefixLen:])
}

// readFileContent reads the content of the file with the given filename in /tmp/
func readFileContent(filename string) ([]byte, error) {
	basePath := os.Args[2]
	filePath := basePath + filename
	fmt.Println("Reading file: ", filePath)
	return os.ReadFile(filePath)
}

// getRequestURL parses the request URL from the request
func getRequestURL(req []byte) []byte {
	// request line \r\n headers \r\n request body \r\n
	requestLine := bytes.Split(req, []byte("\r\n"))[0]
	reqURLByte := bytes.Split(requestLine, []byte(" "))[1]
	fmt.Println("Request URL: ", string(reqURLByte), " len: ", len(reqURLByte))
	return reqURLByte
}
