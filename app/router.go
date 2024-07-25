package main

import (
	"bytes"
)

// routeAndHandleRequest routes the request to the corresponding handler based on the request URL
func routeAndHandleRequest(req *HTTPRequest) *HTTPResponse {
	response := &HTTPResponse{
		Status:  HTTP_CODE_OK,
		Headers: make(map[string]string),
		Body:    make([]byte, 0),
	}

	reqURL := req.URL
	if isRootEndpoint(reqURL) {
		handleRootEndpoint(response)
	} else if isEchoEndpoint(reqURL) {
		handleEchoEndpoint(req, response)
	} else if isUserAgentEndpoint(reqURL) {
		handleUserAgentEndpoint(req, response)
	} else if isFileEndPoint(reqURL) {
		handleFileEndpoint(req, response)
	} else {
		handleNotFoundEndpoint(response)
	}

	return response
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
