package core

import (
	"fmt"
	"net"
	"log"
	"server/tools/cache"
	"server/tools/protocol"
)

const (
	MAX_KEY_LENGTH = 250
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
	//var received_message []byte
	server.socket = listener
	for {
		// Accept waits for and returns the next connection to the listener.
		fmt.Println("Waiting for connection...")
		connection, err := listener.Accept()
		fmt.Println("Accepted!")
		if err != nil {
			fmt.Println(err) // TODO: replace by kind of traceback
			continue
		}
		// handle the connection

		server.connections[connection.RemoteAddr().String()] = connection
		go server.dispatch(connection.RemoteAddr().String())

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

func (server *server) dispatch(address string) {
	received_message := make([]byte, MAX_KEY_LENGTH)
	fmt.Println("Retrieving connection's data from ", address)
	connection := server.connections[address]
	n, err := connection.Read(received_message[0 : ])
	fmt.Println("Done!")
	if err != nil {
		fmt.Println("Dispatching error: ", err, " Message: ", received_message)
		server.makeResponse(connection, []byte("ERROR\r\n"), 5)
	} else {
		// here message should be dispatched
		fmt.Println("Server has received a message: ", string(received_message[0 : n]))
		some := protocol.ParseProtocolHeaders(string(received_message))
		fmt.Println(some)
		server.makeResponse(connection, []byte("OK\r\n"), 2)
	}
}

func (server *server) breakConnection(con *net.Conn) {
	connection := *con
	delete(server.connections, connection.RemoteAddr().String())
	err := connection.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (server *server) makeResponse(connection net.Conn, response_message []byte, length int) {
	length, err := connection.Write(response_message[0 : length])
	if err != nil {
		fmt.Println("Error was occured during making response:", err)
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
	fmt.Println("Server is stopped.")
}
