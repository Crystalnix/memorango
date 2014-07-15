package main

import (
	"fmt"
	"client"
)

func main() {

	port := "9999"

	client.Client(port, []byte("get key 123 123 123 9292929292 noreply\r\n"))

	var input string
	fmt.Scanln(&input)
}
