package core

import (
	"fmt"
	"net"
	"encoding/gob"
	"container/list"
)

type server struct {
	port string
	connections list.List
	storage map[string] []byte  // TODO: replace by type Storage, which will be implemented in core.storage
}

func (server *server) run() {
	link, err := net.Listen("tcp", ":" + server.port)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// accept a connection
		connection, err := link.Accept()
		if err != nil {
			fmt.Println(err) // TODO: replace by kind of traceback
			continue
		}
		// handle the connection
		server.connections.PushBack(connection)
		go server.dispatch(&connection)
	}
}

func (server *server) dispatch(connection *net.Conn){
	var received_message string
	err := gob.NewDecoder(*connection).Decode(&received_message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Received: ", received_message)
	}
	defer server.breakConnection(connection)
}

func (server *server) breakConnection(connection *net.Conn) {
	server.connections.Remove(&connection)
	err := connection.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func RunServer(port string) *server{
	_server := new(server)
	_server.port = port
	_server.storage = make(map[string] []byte) // TODO: replace by Storage initialization
	go _server.run()
	return _server
}
