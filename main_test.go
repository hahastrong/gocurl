package main

import (
	"os"
	"testing"
)

func Test_Request(m *testing.T) {
	os.Args = []string{"cmd", "-X", "POST", "-H", "Content-Type: application/json", "-H", "Custom-Header: custom-value", "-d", `{"data": 1}`, "-p", "tlskey.txt", "-w", "1", "https://www.baidu.com"}

	main()

	// test for GET
	os.Args = []string{"cmd", "-X", "GET", "-H", "Content-Type: application/json", "-H", "Custom-Header: custom-value", "-d", `{"data": 1}`, "-p", "tlskey.txt", "-w", "1", "https://www.baidu.com"}
	main()
}
