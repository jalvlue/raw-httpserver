package main

import (
	"fmt"
	"net"
	"os"
)

func init() {
	if len(os.Args) == 3 {
		FileBasePath = os.Args[2]
	} else {
		FileBasePath = "/tmp/files/"
	}
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
// routes the request to corresponding handler,
// and writes the response back to the connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	rawRequest, err := readRequest(conn)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		return
	}

	fmt.Println("Recieve request, read ", len(rawRequest), " bytes")
	fmt.Println(string(rawRequest))

	req := parseRequest(rawRequest)
	resp := routeAndHandleRequest(req)
	n, err := conn.Write(resp.toBytes())
	if err != nil {
		fmt.Println("Error writing connection: ", err.Error())
		return
	}
	fmt.Println("Send response back, wrote ", n, " bytes")
}
