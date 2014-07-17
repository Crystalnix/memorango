package main

import (
	"fmt"
	"flag"
	server "server/core"
)

func main() {


	port := flag.String("p", "11211", "Port to listen (non required - default port is 11211)")
	memory_amount_mb := flag.Int("m", 0, "Amount of memory to allocate (Mb)")
	//daemonize := flag.Bool("d", false, "Run process as background")

	flag.Parse()
	fmt.Printf("Run memcached on 127.0.0.1:%s with %d mb allocated memory.\nType 'stop' to stop the server. \n",
		       *port, *memory_amount_mb)
	_server := server.RunServer(*port, int64(*memory_amount_mb) * 1024 /* let's convert to bytes */)
	defer server.StopServer(_server)
	for {
		var input string
		fmt.Scanln(&input)
		if input == "stop" {
			break
		}
	}
}
