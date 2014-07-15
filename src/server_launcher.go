package main

import (
	"fmt"
	"flag"
	server "server/core"
)

func main() {


	port := flag.String("p", "", "Port to listen (non required - default port is 11211)")
	memory_amount_mb := flag.Int("m", 0, "Amount of memory to allocate (Mb)")
	daemonize := flag.Bool("d", false, "Run process as background")

	flag.Parse()
	if len(*port) == 0 {
		*port = "11211"
	}


	fmt.Printf("127.0.0.1:%s -m=%d -d=%b\r\n", *port, *memory_amount_mb, *daemonize)

	if *daemonize {
		fmt.Println("Run background process...")
		go server.RunServer(*port, int64(*memory_amount_mb) * 1024 /* let's convert to bytes */)
	} else {
		_server := server.RunServer(*port, int64(*memory_amount_mb) * 1024 /* let's convert to bytes */)
		defer server.StopServer(_server)
	}
	var input string
	fmt.Scanln(&input)
}