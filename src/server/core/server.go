package core

import (
	"fmt"
	"net"
	"log"
	"server/tools/cache"
	"server/tools/protocol"
	"io"
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

	fmt.Println("Retrieving connection's data from ", address)
	connection := server.connections[address]
	for { // let's loop the process for open connection, until it will get closed.
		received_message := make([]byte, MAX_KEY_LENGTH)
		n, err := connection.Read(received_message[0 : ])
		fmt.Println("Connection stream was read.")
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection is closed.")
				server.breakConnection(connection)
				break

			}
			fmt.Println("Dispatching error: ", err, " Message: ", received_message)
			server.makeResponse(connection, []byte("ERROR\r\n"), 5)
		} else {
			// here message should be dispatched
			fmt.Println("Server has received a message: ", string(received_message[0 : n]))
			parsed_request := protocol.ParseProtocolHeaders(string(received_message))
			// handle some
			response_message, err := parsed_request.HandleRequest(server.storage)
			fmt.Println("Server is sending response:\n", string(response_message[0 : len(response_message)]))
			if parsed_request.Reply() {
				server.makeResponse(connection, response_message, len(response_message))
			}
			if err != nil {
				server.breakConnection(connection)
				break
			}
		}
	}
}

func (server *server) breakConnection(connection net.Conn) {
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
	return
}

func RunServer(port string, memory int64) *server {
	server := new(server)
	server.socket = nil
	server.port = port
	server.storage = cache.New(memory)
	server.connections = make(map[string] net.Conn)
	go server.run()
	return server
}

func StopServer(server *server) {
	server.stop()
	fmt.Println("Server is stopped.")
}
