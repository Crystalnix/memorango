package main

import (
	"fmt"
	"flag"
	server "server/core"
)

func main() {
	port := ""
	var memory_amount_mb int



	flag.Parse()
	flag.Var(&memory_amount_mb, "m", "Amount of memory to allocate (Mb)")
	flag.Var(&port, "p", "Port to listen (non required - default port is 11211)")
	if !port {
		port = "11211"
	}

	fmt.Printf("127.0.0.1:" + string(port) + "\n")

	_server := server.RunServer(port, int64(memory_amount_mb) * 1024)
	defer server.StopServer(_server)

	var input string
	fmt.Scanln(&input)
}
