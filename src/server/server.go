/*
DEPENDENCIES

- External:

fmt - used for debug and displaying information about server.

net - used for listening of TCP socket, receiving requests and sending responses.

log - used for logging of actions errors and warnings.

io - used for defining EOF

bufio - used for parsing request per byte.

errors - used for creation of custom error.

- Internal:

tools - contains several utilities for convenient usage facilities.

tools/cache - contains structs and methods for data storage

tools/protocol - contains methods and rules for handling requests.

*/
package server

import (
	"fmt"
	"net"
	"log"
	"tools/cache"
	"tools/protocol"
	"io"
	"bufio"
	"errors"
)

const (
	//defines the maximal length of receiving request.
	MAX_KEY_LENGTH = 250
)


// The private server structure keeps information about server's port, active connections, listened socket,
// and pointer to LRUCache structure, which consists methods allowed to retrieve and store data.
type server struct {
	port string
	socket net.Listener
	connections map[string] net.Conn
	storage *cache.LRUCache
}

// Private method of server structure, which starts to listen connection, cache it and delegate it to dispatcher.
func (server *server) run() {
	listener, err := net.Listen("tcp", ":" + server.port)
	if err != nil {
		fmt.Println(err)
		return
	}
	//var received_message []byte
	server.socket = listener
	for {
		// Accept waits for incoming data and returns the next connection to the listener.
		fmt.Println("Waiting for connection...")
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err, server.socket) // TODO: replace by kind of traceback
			if server.socket == nil {
				break
			} else { continue }
		} else {
			// handle the connection
			server.connections[connection.RemoteAddr().String()] = connection
			go server.dispatch(connection.RemoteAddr().String())
		}
	}
}

// Private method of server struct, which closes socket listener and stops serving.
func (server *server) stop() {
	server.storage.FlushAll()
	server.storage = nil
	server.connections = nil
	if server.socket != nil {
  		err := server.socket.Close()
		if err != nil {
			fmt.Println("Error occured during closing socket:", err)
		}
		server.socket = nil
	} else {
		log.Fatal("Server can't be stoped, because socket is undefined.")
	}

}

// Private method of server, which dispatches active incoming connection.
// Function receives address string and uses it as key to retrieve cached connection.
// Fetched connection is getting read by bufio.Reader, parsed to header and data string if it's size was pointed in header.
// Next, the parsed data handles by protocol and writes a response message.
// The process turns in loop until whether input stream will get an EOF or an error will be occurred.
// In the last case it will be return some error message to a client.
// Anyway, at the end connection will be broken up.
func (server *server) dispatch(address string) {

	fmt.Println("Retrieving connection's data from ", address)
	connection := server.connections[address]
	connectionReader := bufio.NewReader(connection)
	// let's loop the process for open connection, until it will get closed.
	for {
		// let's read a header first
		received_message, n, err := readRequest(connectionReader, -1)
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
			// Here the message should be handled
			parsed_request := protocol.ParseProtocolHeader(string(received_message[0 : n - 2]))
			fmt.Println("Header: ", parsed_request)
			received_message, n, err := readRequest(connectionReader, parsed_request.DataLen())
			fmt.Println("Data: ", received_message)
			if err != nil {
				server.breakConnection(connection)
				break
			}
			parsed_request.SetData(received_message, n - 2)
			response_message, err := parsed_request.HandleRequest(server.storage)
			fmt.Println("Server is sending response:\n", string(response_message[0 : len(response_message)]))
			// if there is no flag "noreply" in the header:
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

// This private function serves for reading an input stream per byte till the \r\n terminator
// or until the length param won't be achieved.
// Function returns byte-string, it's length and error.
// If process succeeded, instead of error will be return nil.
// Otherwise, if read data exceed of MAX_KEY_LENGTH, it will be returned var of error type, but also data,
// which was read and it's length (MAX_KEY_LENGTH)
func readRequest(reader *bufio.Reader, length int) ([]byte, int, error){
	buffer := make([]byte, MAX_KEY_LENGTH)
	var symbol byte
	var counter = 0
	if length == 0 { return buffer, 0, nil }
	for index, _ := range buffer {
		if index >= MAX_KEY_LENGTH { return buffer, counter, errors.New("Header length is exceeded.") }
		read, err := reader.ReadByte()
		if err != nil { return buffer, counter, err }
		buffer[index] = read
		counter ++
		if read == '\n' && symbol == '\r' && (length == -1 || index - 1 == length) { break }
		symbol = read
	}
	return buffer, counter, nil
}

// Private method break up the connection, closes it and removes it from cached server's connections.
func (server *server) breakConnection(connection net.Conn) {
	delete(server.connections, connection.RemoteAddr().String())
	err := connection.Close()
	if err != nil {
		fmt.Println("Can't break up connection: ", err)
	}
}

// Private method writes a byte-string to connection output stream.
// Function receives connection, message and length of this message.
func (server *server) makeResponse(connection net.Conn, response_message []byte, length int) {
	length, err := connection.Write(response_message[0 : length])
	if err != nil {
		fmt.Println("Error was occured during making response:", err)
	}
}

// This public function raises up the server.
// Function receives port string, which uses to open socket at pointed port and int64 value,
// which uses for limiting of allowed memory to use.
// Function returns a pointer to server structure with filled and prepared to usage fields.
func RunServer(port string, bytes_of_memory int64) *server {
	server := new(server)
	server.socket = nil
	server.port = port
	server.storage = cache.New(bytes_of_memory)
	server.connections = make(map[string] net.Conn)
	go server.run()
	return server
}

// Public function receives the pointer to server structure, stops the server and inform about it.
func StopServer(server *server) {
	server.stop()
	fmt.Println("Server is stopped.")
}
