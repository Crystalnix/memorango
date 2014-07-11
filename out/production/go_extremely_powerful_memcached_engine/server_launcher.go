package main

import (
	"fmt"
	server "server/core"
)

func main() {
	port := "9999"

	fmt.Printf("127.0.0.1:" + string(port) + "\n")

	_server := server.RunServer(port)
	defer server.StopServer(_server)

	var input string
	fmt.Scanln(&input)
}
