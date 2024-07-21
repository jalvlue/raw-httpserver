package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

type HTTPRequest struct {
	Method  []byte
	URL     []byte
	Headers map[string]string
	Body    []byte
}

type HTTPResponse struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

func NewHTTPResponse() *HTTPResponse {
	return &HTTPResponse{
		Status:  200,
		Headers: make(map[string]string),
	}
}

func (resp *HTTPResponse) addCode(code int) {
	resp.Status = code
}

func (resp *HTTPResponse) addHeader(key string, value string) {
	resp.Headers[key] = value
}

func (resp *HTTPResponse) addHeaders(headers map[string]string) {
	for k, v := range headers {
		resp.addHeader(k, v)
	}
}

func (resp *HTTPResponse) addContent(content []byte) {
	resp.Body = content
}

func (resp *HTTPResponse) toBytes() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", resp.Status, httpStatusText(resp.Status)))
	for key, value := range resp.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	buffer.WriteString("\r\n")
	buffer.Write(resp.Body)

	return buffer.Bytes()
}

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

	// TODO: read request data from connection gracefully
	rawRequest := make([]byte, 1024)
	n, err := conn.Read(rawRequest)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		return
	}
	fmt.Println("Recieve request, read ", n, " bytes")
	fmt.Println(string(rawRequest))

	req := parseRequest(rawRequest)
	n, err = routeRequest(conn, req)
	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		return
	}
	fmt.Println("Send response back, wrote ", n, " bytes")
}

// routeRequest routes the request to the corresponding handler based on the request URL
func routeRequest(conn net.Conn, req HTTPRequest) (int, error) {
	var response HTTPResponse

	reqURL := req.URL
	if isRootEndpoint(reqURL) {
		response.addCode(200)
	} else if isEchoEndpoint(reqURL) {
		msg := getEchoMsg(reqURL)
		response.addHeaders(map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(msg)),
		})
		response.addContent(msg)
	} else if isUserAgentEndpoint(reqURL) {
		userAgent := getHeader(req.Headers, "User-Agent")
		response.addCode(200)
		response.addHeaders(map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(userAgent)),
		})
		response.addContent([]byte(userAgent))
	} else if isFileEndPoint(reqURL) {
		filename := getRequestFilename(reqURL)
		content, err := readFileContent(filename)
		if err != nil {
			response.addCode(404)
		} else {
			response.addCode(200)
			response.addHeaders(map[string]string{
				"Content-Type":   "application/octet-stream",
				"Content-Length": fmt.Sprintf("%d", len(content)),
			})
			response.addContent(content)
		}
	} else {
		response.addCode(404)
	}

	return conn.Write(response.toBytes())
}

// parseRequest parses the raw request into an HTTPRequest struct
func parseRequest(rawRequest []byte) HTTPRequest {
	// Split the raw request into request line, headers, and body
	requestLine, headers, body := splitRequest(rawRequest)

	// Parse the request line
	method, url := parseRequestLine(requestLine)

	// Parse the headers
	headerMap := parseHeaders(headers)

	return HTTPRequest{
		Method:  method,
		URL:     url,
		Headers: headerMap,
		Body:    body,
	}
}

// TODO: check requestHeaders
func splitRequest(rawRequest []byte) (requestLine []byte, requestHeaders [][]byte, requestBody []byte) {
	requestComponents := bytes.Split(rawRequest, []byte("\r\n"))
	requestLine = requestComponents[0]
	requestBody = requestComponents[len(requestComponents)-1]
	requestHeaders = requestComponents[1 : len(requestComponents)-2]

	return
}

// parseRequestLine parses the request line into the method and URL
func parseRequestLine(requestLine []byte) ([]byte, []byte) {
	requestLineComponents := bytes.Split(requestLine, []byte(" "))
	method := requestLineComponents[0]
	url := requestLineComponents[1]

	return method, url
}

// parseHeaders parses the headers into a map
func parseHeaders(headers [][]byte) map[string]string {
	headerMap := make(map[string]string)

	for _, header := range headers {
		headerComponents := bytes.Split(header, []byte(": "))
		headerMap[string(headerComponents[0])] = string(headerComponents[1])
	}

	return headerMap
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
func getEchoMsg(url []byte) []byte {
	prefixLen := len("/echo/")
	return url[prefixLen:]
}

// getHeader returns the value of the header with the given key
func getHeader(reqHeaders map[string]string, headerKey string) string {
	for key, value := range reqHeaders {
		if key == headerKey {
			return value
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

func httpStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	default:
		return "Internal Server Error"
	}
}
