package main

import (
	"bytes"
)

// routeAndHandleRequest routes the request to the corresponding handler based on the request URL
func routeAndHandleRequest(req *HTTPRequest) (response *HTTPResponse, err error) {
	response = &HTTPResponse{
		Headers: make(map[string]string),
	}

	// TODO: route the request to the corresponding handler gracefully
	reqURL := req.URL
	if isRootEndpoint(reqURL) {
		err = handleRootEndpoint(response)
		// err = handleRootEndpoint(req, response)
	} else if isEchoEndpoint(reqURL) {
		err = handleEchoEndpoint(req, response)
	} else if isUserAgentEndpoint(reqURL) {
		err = handleUserAgentEndpoint(req, response)
	} else if isFileEndPoint(reqURL) {
		err = handleFileEndpoint(req, response)
	} else {
		err = handleNotFoundEndpoint(response)
	}

	return
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
