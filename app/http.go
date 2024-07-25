package main

import (
	"bytes"
	"fmt"
	"net"
)

type HTTPCode int

const (
	HTTP_CODE_OK                    HTTPCode = 200
	HTTP_CODE_CREATED               HTTPCode = 201
	HTTP_CODE_NOT_FOUND             HTTPCode = 404
	HTTP_CODE_METHOD_NOT_ALLOWED    HTTPCode = 405
	HTTP_CODE_INTERNAL_SERVER_ERROR HTTPCode = 500
)

func HTTPCode2StatusText(code HTTPCode) string {
	switch code {
	case HTTP_CODE_OK:
		return "OK"
	case HTTP_CODE_CREATED:
		return "Created"
	case HTTP_CODE_NOT_FOUND:
		return "Not Found"
	case HTTP_CODE_METHOD_NOT_ALLOWED:
		return "Method Not Allowed"
	case HTTP_CODE_INTERNAL_SERVER_ERROR:
		return "Internal Server Error"
	default:
		return "Unknown HTTPCode"
	}
}

type HTTPRequest struct {
	Method  []byte
	URL     []byte
	Headers map[string]string
	Body    []byte
}

// getHeader returns the value of the header with the given key
func (req *HTTPRequest) getHeader(headerKey string) string {
	return req.Headers[headerKey]
}

type HTTPResponse struct {
	Status  HTTPCode
	Headers map[string]string
	Body    []byte
}

func NewHTTPResponse() *HTTPResponse {
	return &HTTPResponse{
		Status:  200,
		Headers: make(map[string]string),
	}
}

func (resp *HTTPResponse) setCode(code HTTPCode) {
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

func (resp *HTTPResponse) setContent(content []byte) {
	resp.Body = content
}

// toBytes converts the HTTPResponse to well-formatted slice of bytes
func (resp *HTTPResponse) toBytes() []byte {
	var buffer bytes.Buffer

	// check if the response body should be compressed
	if resp.Headers["Content-Encoding"] == "gzip" {
		compressedBodyBytes, err := gzipCompress(resp.Body)
		if err != nil {
			fmt.Printf("Error compressing content: %v\n", err)
		} else {
			resp.setContent(compressedBodyBytes)
			resp.addHeader("Content-Length", fmt.Sprintf("%d", len(compressedBodyBytes)))
		}
	}

	// write status line
	buffer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", resp.Status, HTTPCode2StatusText(resp.Status)))

	// write headers
	for key, value := range resp.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	buffer.WriteString("\r\n")

	// write body
	buffer.Write(resp.Body)

	return buffer.Bytes()
}

// readRequest reads the raw request content from the connection
func readRequest(conn net.Conn) ([]byte, error) {
	// var buffer bytes.Buffer
	// tmp := make([]byte, 256)
	// conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// for {
	// 	n, err := conn.Read(tmp)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err, ok := err.(net.Error); ok && err.Timeout() {
	// 			break
	// 		}
	// 		return nil, err
	// 	}
	// 	buffer.Write(tmp[:n])
	// }

	// rawReq := buffer.Bytes()
	// return rawReq, nil

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading connection: %v\n", err)
		return nil, err
	}
	return buf[:n], nil
}

// parseRequest parses the raw request into an HTTPRequest struct
func parseRequest(rawRequest []byte) *HTTPRequest {
	// Split the raw request into request line, headers, and body
	requestLine, headers, body := splitRequest(rawRequest)

	// Parse the request line
	method, url := parseRequestLine(requestLine)

	// Parse the headers
	headerMap := parseHeaders(headers)

	return &HTTPRequest{
		Method:  method,
		URL:     url,
		Headers: headerMap,
		Body:    body,
	}
}

// splitRequest splits the raw request into request line, headers, and body
func splitRequest(rawRequest []byte) (requestLine []byte, requestHeaders [][]byte, requestBody []byte) {
	requestComponents := bytes.Split(rawRequest, []byte("\r\n\r\n"))
	requestBody = requestComponents[1]

	requestLineAndHeaders := bytes.Split(requestComponents[0], []byte("\r\n"))
	requestLine = requestLineAndHeaders[0]
	requestHeaders = requestLineAndHeaders[1:]

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
