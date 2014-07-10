package main

import (
	"fmt"
	"server/core"
	"client"
)

func main() {

	port := "11211"

	fmt.Printf("127.0.0.1:" + string(port) + "\n")

	go core.RunServer(port)

	go client.Client(port, "Test")

	var input string
	fmt.Scanf(&input)
}
