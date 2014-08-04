/*
DEPENDENCIES

- External:

fmt - used for debug and displaying information about server.

net - used for listening of TCP socket, receiving requests and sending responses.

log - used for logging of actions errors and warnings.

io - used for defining EOF

bufio - used for parsing request per byte.

errors - used for creation of custom error.

math/rand and time - used for seeding the PRNG

- Internal:

tools - contains several utilities for convenient usage facilities.

tools/cache - contains structs and methods for data storage

tools/protocol - contains methods and rules for handling requests.

tools/stat - contains struct for keeping information about server, its actions and its condition.

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
	"math/rand"
	"time"
	statistic "tools/stat"
	"sync"
	"strings"
)

const (
	//defines the maximal length of receiving request.
	MAX_KEY_LENGTH = 250
)


// The private server structure keeps information about server's port, active connections, listened socket,
// and pointer to LRUCache structure, which consists methods allowed to retrieve and store data.
type server struct {
	tcp_port string
	udp_port string
	listen_address string
	cas_disabled bool
	flush_disabled bool
	connection_limit int
	verbosity_lvl int
	sockets map[string] net.Listener
	connections map[string] net.Conn
	storage *cache.LRUCache
	Stat *statistic.ServerStat
	threadSync sync.WaitGroup
}

// Private method of server structure, which starts to listen connection, cache it and delegate it to dispatcher.
func (server *server) run(connection_type string) {
	server.threadSync.Add(1)
	defer server.threadSync.Done()
	var port string
	if connection_type == "tcp" {
		port = server.tcp_port
	} else if connection_type == "udp" {
		port = server.udp_port
	} else {
		//error log.
		return
	}
	rand.Seed(time.Now().Unix())
	listener, err := net.Listen(connection_type, ":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
	//var received_message []byte
	server.sockets[connection_type] = listener
	for {
		if server.sockets == nil {
			break
		}
		// Accept waits for incoming data and returns the next connection to the listener.
		connection, err := server.sockets[connection_type].Accept()
		if err != nil {
			//TODO: log this
			continue
		} else {
			// handle the connection
			if len(server.listen_address) > 0 {
				if strings.Split(connection.RemoteAddr().String(), ":")[0] != server.listen_address {
					//TODO: log this
					continue
				}
			}
			if len(server.connections) >= server.connection_limit{
				//TODO: log this
				continue
			}
			server.connections[connection.RemoteAddr().String()] = connection
			server.Stat.Current_connections ++
			server.Stat.Total_connections ++
			go server.dispatch(connection.RemoteAddr().String())
		}
	}
}

// Private method of server struct, which closes socket listener and stops serving.
func (server *server) stop() {
	for address, connection := range server.connections {
		if server.breakConnection(connection) {
			fmt.Println("Closed connection at", address)
		} else {
			fmt.Println("Can't close connection at", address)
		}
	}
	if server.sockets != nil {
		for conn_type, socket := range server.sockets {
			err := socket.Close()
			if err != nil {
				fmt.Println("Error occured during closing " + conn_type + " socket:", err)
			}
		}
		server.sockets = nil
	} else {
		log.Fatal("Server can't be stoped, because socket is undefined.")
	}
	fmt.Println("Waiting for goroutines...")
	server.threadSync.Wait() // Waiting for all goroutines while done their jobs.
	server.storage.FlushAll()
	server.storage = nil
	server.connections = nil
	server.Stat = nil
}

// Private method of server, which dispatches active incoming connection.
// Function receives address string and uses it as key to retrieve cached connection.
// Fetched connection is getting read by bufio.Reader, parsed to header and data string if it's size was pointed in header.
// Next, the parsed data handles by protocol and writes a response message.
// The process turns in loop until whether input stream will get an EOF or an error will be occurred.
// In the last case it will be return some error message to a client.
// Anyway, at the end connection will be broken up.
func (server *server) dispatch(address string) {
	server.threadSync.Add(1)
	defer server.threadSync.Done()

	//fmt.Println("Retrieving connection's data from ", address)
	//TODO: log this

	connection := server.connections[address]
	connectionReader := bufio.NewReader(connection)
	// let's loop the process for open connection, until it will get closed.
	for {
		// let's read a header first
		received_message, n, err := readRequest(connectionReader, -1)
		if err != nil {
			if err == io.EOF {
				//TODO: log this
				server.breakConnection(connection)
				break
			}
			fmt.Println("Dispatching error: ", err, " Message: ", received_message)
			//TODO: log this
			if !server.makeResponse(connection, []byte("ERROR\r\n"), 5){
				break
			}
		} else {
			// Here the message should be handled
			server.Stat.Read_bytes += uint64(n)
			parsed_request := protocol.ParseProtocolHeader(string(received_message[ : n - 2]))
			fmt.Println("Header: ", parsed_request)

			if (parsed_request.Command() == "cas" || parsed_request.Command() == "gets") && server.cas_disabled {
				//TODO: log this
				continue
			} else if parsed_request.Command() == "flush_all" && server.flush_disabled {
				//TODO: log this
				continue
			}

			if parsed_request.DataLen() > 0 {
				received_message, n, err := readRequest(connectionReader, parsed_request.DataLen())
				fmt.Println("Data length / read: ", len(received_message), n)
				//TODO: log this
				if err != nil {
					server.breakConnection(connection)
					break
				}
				parsed_request.SetData(received_message[0 : ])
			}
			response_message, err := parsed_request.HandleRequest(server.storage, server.Stat)
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
// Function receives pointer to bufio.Reader, which contains a stream and length, which of required data.
// If the length param equals -1, it means, that data's length is undefined, and it will be read until first \r\n seq.
// Function returns byte-string, it's length and (optionally) error.
// If process succeeded, instead of error will be returned a nil.
// Otherwise, will be returned occurred error, and also read data and it's length.
func readRequest(reader *bufio.Reader, length int) ([]byte, int, error){
	buffer := []byte("")
	var prev_symbol byte
	var counter = 0
	var token_counter = 0
	if length == 0 { return buffer, 0, nil }
	for {
		read, err := reader.ReadByte()
		if err != nil {
			//fmt.Println("Num: ", counter," read: ", read, " Err: ", err)
			return buffer, counter, err
		}
		buffer = append(buffer, read)
		counter ++
		if length == -1 || counter - 2 == length {
			if read == '\n' && prev_symbol == '\r' {
				return buffer[ : len(buffer) - 2], counter, nil
			} else {
				if length != -1 {
					return buffer, counter, errors.New("Length was achieved, but terminator wasn't met.")
				}
			}
		}
		if read != ' ' && length == -1 /* in case of header of unknown length */{
			token_counter ++
			if token_counter > MAX_KEY_LENGTH {
				return buffer, counter, errors.New("Maximal key length is exceeded.")
			}
		} else {
			token_counter = 0
		}
		prev_symbol = read
	}
}


// Private method break up the connection, closes it and removes it from cached server's connections.
func (server *server) breakConnection(connection net.Conn) bool {
	if server.sockets == nil{
		return false
	}
	delete(server.connections, connection.RemoteAddr().String())
	err := connection.Close()
	if err != nil {
		//TODO: log this
		return false
	}
	server.Stat.Current_connections --
	return true
}

// Private method writes a byte-string to connection output stream.
// Function receives connection, message and length of this message.
func (server *server) makeResponse(connection net.Conn, response_message []byte, length int) bool {
	length, err := connection.Write(response_message[0 : length])
	if err != nil {
		//TODO: log this
		return server.breakConnection(connection)
	}
	server.Stat.Written_bytes += uint64(length)
	return true
}

// This public function raises up the server.
// Function receives following params:
// tcp_port string, which uses to open tcp socket at pointed port,
// udp_port string, which uses to open tcp socket at pointed port,
// address, which specified an only ip address which server will listen to,
// max_connections, sets a limit of maximal number of active connections,
// cas, flush - flags which forbid of usage such commands if value = true,
// verbosity - defines the dept of verbosity for server
// and bytes_of_memory, which uses for limiting allocated memory.
// Function returns a pointer to a server structure with filled and prepared to usage fields.
func NewServer(tcp_port string, udp_port string, address string, max_connections int, cas bool, flush bool,
	           verbosity int, bytes_of_memory int64) *server {
	server := new(server)
	server.sockets = nil
	server.tcp_port = tcp_port
	server.udp_port = udp_port
	server.cas_disabled = cas
	server.flush_disabled = flush
	server.connection_limit = max_connections
	server.listen_address = address
	server.verbosity_lvl = verbosity
	server.storage = cache.New(bytes_of_memory)
	server.connections = make(map[string] net.Conn)
	server.Stat = statistic.New(bytes_of_memory)
	return server
}

func (server *server) RunServer() {
	server.sockets = make(map[string] net.Listener)
	go server.run("tcp")
	if len(server.udp_port) > 0 {
		go server.run("udp")
	}
}

// Public function receives the pointer to server structure, stops the server and inform about it.
func (server *server) StopServer() {
	server.stop()
	//TODO: log this
}
