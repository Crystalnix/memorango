package client

/*
	This functionality is temporary
	Client should be replaced to something more powerful, than now.
	At the moment the only action it does - send message to a localhost:port server.
*/


import (
	"net"
	"fmt"
	"encoding/gob"
)

func Client(port string, message string) {
	// connect to the server
	c, err := net.Dial("tcp", "127.0.0.1:" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
	// send the message
	fmt.Println("Sending", message)
	err = gob.NewEncoder(c).Encode(message)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}
