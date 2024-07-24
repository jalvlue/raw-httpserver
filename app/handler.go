package main

import (
	"fmt"
	"os"
)

// handleRootEndpoint handles the root endpoint
func handleRootEndpoint(resp *HTTPResponse) (err error) {
	resp.setCode(HTTP_CODE_OK)

	return
}

// handleEchoEndpoint handles the echo endpoint
func handleEchoEndpoint(req *HTTPRequest, resp *HTTPResponse) (err error) {
	msg := getEchoMsg(req.URL)
	resp.setCode(HTTP_CODE_OK)
	resp.addHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprintf("%d", len(msg)),
	})
	resp.setContent(msg)

	return
}

// handleUserAgentEndpoint handles the user-agent endpoint
func handleUserAgentEndpoint(req *HTTPRequest, resp *HTTPResponse) (err error) {
	userAgent := req.getHeader("User-Agent")
	resp.setCode(HTTP_CODE_OK)
	resp.addHeaders(map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprintf("%d", len(userAgent)),
	})
	resp.setContent([]byte(userAgent))

	return
}

// handleFileEndpoint handles the file endpoint
func handleFileEndpoint(req *HTTPRequest, resp *HTTPResponse) error {
	m := string(req.Method)
	switch m {
	case "GET":
		return handleFileEndpointGET(req, resp)
	case "POST":
		return handleFileEndpointPOST(req, resp)
	default:
		return handleFileEndpointSupport(resp)
	}
}

// handleNotFoundEndpoint handles the not found endpoint
func handleNotFoundEndpoint(resp *HTTPResponse) (err error) {
	resp.setCode(HTTP_CODE_NOT_FOUND)

	return
}

// handleFileEndpointGET handles the GET method for the file endpoint
func handleFileEndpointGET(req *HTTPRequest, resp *HTTPResponse) (err error) {
	filename := getRequestFilename(req.URL)
	content, err := readFileContent(filename)
	if err != nil {
		if err == os.ErrNotExist {
			fmt.Printf("Error file not found: %v\n", err)
			resp.setCode(HTTP_CODE_NOT_FOUND)
		}
		fmt.Printf("Error reading file: %v\n", err)
		resp.setCode(HTTP_CODE_INTERNAL_SERVER_ERROR)
	} else {
		resp.setCode(HTTP_CODE_OK)
		resp.addHeaders(map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": fmt.Sprintf("%d", len(content)),
		})
		resp.setContent(content)
	}

	return
}

// handleFileEndpointPOST handles the POST method for the file endpoint
func handleFileEndpointPOST(req *HTTPRequest, resp *HTTPResponse) (err error) {
	filename := getRequestFilename(req.URL)

	err = writeFileContent(filename, req.Body)
	if err != nil {
		fmt.Printf("Error writing or creating file: %v\n", err)
		resp.setCode(HTTP_CODE_INTERNAL_SERVER_ERROR)
	} else {
		resp.setCode(HTTP_CODE_CREATED)
	}

	return
}

// handleFileEndpointSupport handles the unknown method for the file endpoint
func handleFileEndpointSupport(resp *HTTPResponse) (err error) {
	resp.setCode(HTTP_CODE_METHOD_NOT_ALLOWED)

	return
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
