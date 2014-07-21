package client

/*
	This functionality is temporary
	Client should be replaced to something more powerful, than now.
	At the moment the only action it does - send message to a localhost:port server.
*/


import (
	"net"
	"fmt"
	//"encoding/gob"
)

func Client(port string, message []byte) {

	// connect to the server
	c, err := net.Dial("tcp", "127.0.0.1:" + port)
	defer c.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	// send the message
	fmt.Println("Client has sent a message: ", message)
	_, err = c.Write(message)

	if err != nil {
		fmt.Println("Error was occured during data transmission:", err)
	}
	var received_message = make([]byte, 255)
	_, err = c.Read(received_message[0: ])
	c.Close()
	if err != nil {
		fmt.Println("Error was occured during receiving of data:", err)
	} else {
		fmt.Println("Server response: ", string(received_message))
	}
}
