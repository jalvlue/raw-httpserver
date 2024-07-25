package main

import (
	"fmt"
	"os"
)

// handleRootEndpoint handles the root endpoint
func handleRootEndpoint(resp *HTTPResponse) {
	resp.setCode(HTTP_CODE_OK)
}

// handleEchoEndpoint handles the echo endpoint
func handleEchoEndpoint(req *HTTPRequest, resp *HTTPResponse) {
	msg := getEchoMsg(req.URL)
	resp.setCode(HTTP_CODE_OK)
	resp.addHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprintf("%d", len(msg)),
	})
	resp.setContent(msg)
}

// handleUserAgentEndpoint handles the user-agent endpoint
func handleUserAgentEndpoint(req *HTTPRequest, resp *HTTPResponse) {
	userAgent := req.getHeader("User-Agent")
	resp.setCode(HTTP_CODE_OK)
	resp.addHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprintf("%d", len(userAgent)),
	})
	resp.setContent([]byte(userAgent))
}

// handleFileEndpoint handles the file endpoint
func handleFileEndpoint(req *HTTPRequest, resp *HTTPResponse) {
	m := string(req.Method)
	switch m {
	case "GET":
		handleFileEndpointGET(req, resp)
	case "POST":
		handleFileEndpointPOST(req, resp)
	default:
		handleFileEndpointSupport(resp)
	}
}

// handleNotFoundEndpoint handles the not found endpoint
func handleNotFoundEndpoint(resp *HTTPResponse) {
	resp.setCode(HTTP_CODE_NOT_FOUND)
}

// handleFileEndpointGET handles the GET method for the file endpoint
func handleFileEndpointGET(req *HTTPRequest, resp *HTTPResponse) {
	filename := getRequestFilename(req.URL)
	content, err := readFileContent(filename)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			fmt.Printf("Error file not found: %v\n", pathErr)
			resp.setCode(HTTP_CODE_NOT_FOUND)
		} else {
			fmt.Printf("Error reading file : %v\n", err)
			resp.setCode(HTTP_CODE_INTERNAL_SERVER_ERROR)
		}
		return
	}

	resp.setCode(HTTP_CODE_OK)
	resp.addHeaders(map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": fmt.Sprintf("%d", len(content)),
	})
	resp.setContent(content)
}

// handleFileEndpointPOST handles the POST method for the file endpoint
func handleFileEndpointPOST(req *HTTPRequest, resp *HTTPResponse) {
	filename := getRequestFilename(req.URL)

	err := writeFileContent(filename, req.Body)
	if err != nil {
		fmt.Printf("Error writing or creating file: %v\n", err)
		resp.setCode(HTTP_CODE_INTERNAL_SERVER_ERROR)
	} else {
		resp.setCode(HTTP_CODE_CREATED)
	}
}

// handleFileEndpointSupport handles the unknown method for the file endpoint
func handleFileEndpointSupport(resp *HTTPResponse) {
	resp.setCode(HTTP_CODE_METHOD_NOT_ALLOWED)
}

// getEchoMsg returns the message to be echoed back for the echo endpoint
func getEchoMsg(url []byte) []byte {
	prefixLen := len("/echo/")
	return url[prefixLen:]
}

// getRequestFilename returns the filename from the request URL
func getRequestFilename(url []byte) string {
	prefixLen := len("/files/")
	return string(url[prefixLen:])
}
