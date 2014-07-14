package core

import (
	"fmt"
	"net"
	"encoding/gob"
	"log"
	"server/tools/cache"
)

type server struct {
	port string
	socket net.Listener
	connections map[string] net.Conn
	storage *cache.LRUCache
}

func (server *server) run() {
	listener, err := net.Listen("tcp", ":" + server.port)
	if err != nil {
		fmt.Println(err)
		return
	}
	server.socket = listener
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

func (server *server) stop() {
	if server.socket != nil {
  		err := server.socket.Close()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		log.Fatal("Server can't be stoped, because socket is undefined.")
	}
}

func (server *server) dispatch(connection net.Conn) {
	defer server.breakConnection(connection)
	var received_message string
	err := gob.NewDecoder(connection).Decode(&received_message)
	if err != nil {
		fmt.Println(err)
	} else {
		// here message should be dispatched
		fmt.Println("Server has received a message: ", received_message)
		server.makeResponse(connection, "Your message '" + received_message + "' was received.")
	}
}

func (server *server) breakConnection(connection net.Conn) {
	delete(server.connections, connection.LocalAddr().String())
	err := connection.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (server *server) makeResponse(connection net.Conn, response_message string) {
	err := gob.NewEncoder(connection).Encode(response_message)
	if err != nil {
		fmt.Println(err)
	}
}

func RunServer(port string, memory int64) *server {
	_server := new(server)
	_server.socket = nil
	_server.port = port
	_server.storage = cache.New(memory)
	_server.connections = make(map[string] net.Conn)
	go _server.run()
	return _server
}

func StopServer(server *server) {
	server.stop()
	fmt.Println("Server is stoped.")
}
