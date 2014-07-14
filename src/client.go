package main

import (
	"fmt"
	"client"
)

func main() {

	port := "9999"

	client.Client(port, "Test")

	var input string
	fmt.Scanln(&input)
}
