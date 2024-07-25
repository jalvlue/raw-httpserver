package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
)

var FileBasePath string

// readFileContent reads the content of the file with the given filename in FileBasePath
func readFileContent(filename string) ([]byte, error) {
	filePath := FileBasePath + filename
	fmt.Println("Trying to read file: ", filePath)
	return os.ReadFile(filePath)
}

// writeFileContent writes the content to the file with the given filename in FileBasePath,
// creating the file if it does not exist
func writeFileContent(filename string, content []byte) error {
	filePath := FileBasePath + filename
	fmt.Println("Trying to write file: ", filePath)

	// each number of the permission is a bit, so 0644 is 110 100 100 in binary
	// the first 1 is for the file type, the next 3 bits are for the owner's permission,
	// the next 3 bits are for the group's permission, and the last 3 bits are for the other's permission
	// read permission is 100 in binary, write permission is 010 in binary, and execute permission is 001 in binary
	// 6 is 110 in binary, which means the owner can read and write the file
	// 4 is 100 in binary, which means the group can only read the file
	// so 0644 means the owner can read and write the file, and the group and the other can only read the file
	return os.WriteFile(filePath, content, 0644)
}

// gzipCompress compresses the content using gzip
func gzipCompress(content []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(content)
	if err != nil {
		return nil, err
	}
	writer.Close()

	return buf.Bytes(), nil
}
