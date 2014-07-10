package core

import (
	"encoding/gob"
	"fmt"
	"net"
)

func server(port string) {
	// listen on a port
	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Println(err)
		return }
	for {
		// accept a connection
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue }
		// handle the connection
		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {
	// receive the message
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
	fmt.Println(err)
	} else {
	fmt.Println("Received", msg)
	}
	c.Close()
}

