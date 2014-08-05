package main

import (
	"fmt"
	"client"
)

func main() {
	port := "9999"
	var input string
	fmt.Scanln(&input)
	client.Client(port, []byte(input))
	fmt.Scanln(&input)
}
