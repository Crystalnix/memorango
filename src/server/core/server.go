package core

import (
	"fmt"
	"net"
	"encoding/gob"
)

type server struct {
	port string
	connections map[string] net.Conn
	storage map[string] []byte  // TODO: replace by type Storage, which will be implemented in core.storage
}

func (server *server) run() {
	listener, err := net.Listen("tcp", ":" + server.port)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// accept a connection
		// Accept waits for and returns the next connection to the listener.
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println(err) // TODO: replace by kind of traceback
			continue
		}
		// handle the connection
		server.connections[connection.LocalAddr().String()] = connection
		go server.dispatch(connection)  // may be it has sense to use pointers for optimization of process.
	}
}

func (server *server) dispatch(connection net.Conn){
	var received_message string
	err := gob.NewDecoder(connection).Decode(&received_message)
	if err != nil {
		fmt.Println(err)
	} else {
		// here message should be dispatched
		fmt.Println("Received: ", received_message)
	}
	fmt.Println(connection)
	server.breakConnection(connection)
}

func (server *server) breakConnection(connection net.Conn) {
	fmt.Println(connection)
	delete(server.connections, connection.LocalAddr().String())
	err := connection.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func RunServer(port string) *server{
	_server := new(server)
	_server.port = port
	_server.storage = make(map[string] []byte) // TODO: replace by Storage initialization
	_server.connections = make(map[string] net.Conn)
	go _server.run()
	return _server
}
