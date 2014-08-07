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
	"net"
	"log"
	"log/syslog"
	"tools/cache"
	"tools/protocol"
	"io"
	"os"
	"bufio"
	"errors"
	"math/rand"
	"time"
	statistic "tools/stat"
	"sync"
	"strings"
	"io/ioutil"
	"tools/stat"
)

const (
	//defines the maximal length of receiving request.
	MAX_KEY_LENGTH = 250
)

//TODO: Tests for logger
// Structure for managing of output information during work of server.
type ServerLogger struct {
	info    	*log.Logger
	warning 	*log.Logger
	error   	*log.Logger
	syslogger  	*log.Logger
}

// Initialization of information management system for server by verbosity level:
// 0 - only errors,
// 1 - errors and warnings,
// 2 - errors, warnings and info.
func NewServerLogger(verbosity int) *ServerLogger {
	var result ServerLogger
	result.error = log.New(os.Stderr, "Error: ", log.Ldate | log.Ltime | log.Lshortfile)
	if verbosity >= 1 {
		result.warning = log.New(os.Stdout, "Warning: ", log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		result.warning = log.New(ioutil.Discard, "", 0)
	}
	if verbosity == 2 {
		result.info = log.New(os.Stdout, "Info: ", log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		result.info = log.New(ioutil.Discard, "", 0)
	}
	result.syslogger, _ = syslog.NewLogger(syslog.LOG_ERR, log.Ldate | log.Ltime | log.Lshortfile)
	result.syslogger.SetPrefix("MemoranGo ")
	return &result
}

// Display info-level
func (l *ServerLogger) Info(args ...interface{}){
	l.info.Println(args)
}

// Display error-level
func (l *ServerLogger) Error(args ...interface{}){
	l.syslogger.Println(args)
	l.error.Println(args)
}

// Display warning-level
func (l *ServerLogger) Warning(args ...interface{}){
	l.syslogger.Println(args)
	l.warning.Println(args)
}

// The private server structure keeps information about server's port, active connections, listened socket,
// and pointer to LRUCache structure, which consists methods allowed to retrieve and store data.
type server struct {
	tcp_port string
	udp_port string
	listen_address string
	cas_disabled bool
	flush_disabled bool
	connection_limit int
	sockets map[string] net.Listener
	connections map[string] net.Conn
	storage *cache.LRUCache
	Stat *statistic.ServerStat
	ThreadSync sync.WaitGroup
	Logger *ServerLogger
}

// Private method of server structure, which starts to listen connection, cache it and delegate it to dispatcher.
func (server *server) run(connection_type string) {
	defer server.ThreadSync.Done()
	var port string
	if connection_type == "tcp" {
		port = server.tcp_port
	} else if connection_type == "udp" {
		port = server.udp_port
	} else {
		server.Logger.Error("There is no such type of connection:", connection_type)
		return
	}
	rand.Seed(time.Now().Unix())
	listener, err := net.Listen(connection_type, ":" + port)
	if err != nil {
		server.Logger.Error("Couldn't establish listener:", err)
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
			server.Logger.Warning("Connection couldn't be accepted:", err)
			continue
		} else {
			// handle the connection
			if len(server.listen_address) > 0 {
				if strings.Split(connection.RemoteAddr().String(), ":")[0] != server.listen_address {
					server.Logger.Warning("Connection address", connection.RemoteAddr().String(), "doesn't match with", server.listen_address)
					connection.Close()
					continue
				}
			}
			if len(server.connections) >= server.connection_limit{
				server.Logger.Error("Impossible connect to the server. There are too much active connections right now.")
				connection.Close()
				continue
			}
			addr := connection.RemoteAddr().String()
			if server.connections[addr] == nil {
				server.connections[addr] = connection
				server.Stat.Connections[addr] = stat.NewConnStat(connection)
				server.Stat.Current_connections ++
			}
			server.Stat.Total_connections ++
			go server.dispatch(connection.RemoteAddr().String())
		}
	}
}

// Private method of server struct, which closes socket listener and stops serving.
func (server *server) stop() {
	for address, connection := range server.connections {
		if server.breakConnection(connection) {
			server.Logger.Info("Close connection at", address)
		} else {
			server.Logger.Warning("Impossible to close connection at", address)
		}
	}
	if server.sockets != nil {
		for conn_type, socket := range server.sockets {
			err := socket.Close()
			if err != nil {
				server.Logger.Error("Error occured during closing " + conn_type + " socket:", err)
			}
		}
		server.sockets = nil
	} else {
		server.Logger.Error("Server can't be stoped, because socket is undefined.")
	}
	server.Logger.Info("Waiting for ending process of goroutines...")
	server.ThreadSync.Wait() // Waiting for all goroutines while done their jobs.
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
	server.Stat.Connections[address].State = "conn_new_cmd"
	server.ThreadSync.Add(1)
	defer server.ThreadSync.Done()
	connection := server.connections[address]
	connectionReader := bufio.NewReader(connection)
	// let's loop the process for open connection, until it will get closed.
	for {
		// let's read a header first
		server.Stat.Connections[address].State = "conn_read"
		received_message, n, err := readRequest(connectionReader, -1)
		server.Stat.Connections[address].Cmd_hit_ts = time.Now().Unix()
		if err != nil {
			server.Stat.Connections[address].State = "conn_swallow"
			if err == io.EOF {
				server.Logger.Info("Input stream has got EOF, and now is being closed.")
				server.breakConnection(connection)
				break
			}
			server.Logger.Warning("Dispatching error: ", err, " Message: ", received_message)
			if !server.makeResponse(connection, []byte("ERROR\r\n"), 5){
				break
			}
		} else {
			// Here the message should be handled
			server.Stat.Read_bytes += uint64(n)
			parsed_request := protocol.ParseProtocolHeader(string(received_message[ : n - 2]))
			server.Logger.Info("Header: ", *parsed_request)

			if (parsed_request.Command() == "cas" || parsed_request.Command() == "gets") && server.cas_disabled {
				server.Logger.Warning("Cas command is forbidden.")
				continue
			} else if parsed_request.Command() == "flush_all" && server.flush_disabled {
				server.Logger.Warning("Flush all command is forbidden.")
				continue
			}

			if parsed_request.DataLen() > 0 {
				server.Stat.Connections[address].State = "conn_nread"
				received_message, _, err := readRequest(connectionReader, parsed_request.DataLen())
				if err != nil {
					server.Logger.Error("Error occurred while reading data:", err)
					server.breakConnection(connection)
					break
				}
				parsed_request.SetData(received_message[0 : ])
			}
			server.Logger.Info("Start handling request:", *parsed_request)
			response_message, err := parsed_request.HandleRequest(server.storage, server.Stat)
			server.Logger.Info("Server is sending response:\n", string(response_message[0 : len(response_message)]))
			// if there is no flag "noreply" in the header:
			if parsed_request.Reply() {
				server.Stat.Connections[address].State = "conn_write"
				server.makeResponse(connection, response_message, len(response_message))
			}
			if err != nil {
				server.Logger.Error("Impossible to send response:", err)
				server.breakConnection(connection)
				break
			}
		}
		server.Stat.Connections[address].State = "conn_waiting"
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
	address := connection.RemoteAddr().String()
	server.Stat.Connections[address].State = "conn_closing"
	defer delete(server.Stat.Connections, address)

	if server.sockets == nil{
		return false
	}
	delete(server.connections, connection.RemoteAddr().String())
	err := connection.Close()
	if err != nil {
		server.Logger.Warning("Impossible to break connection:", err)
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
		server.Logger.Warning("Error occurred during writing data to output stream:", err)
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
	server.storage = cache.New(bytes_of_memory)
	server.connections = make(map[string] net.Conn)
	server.Stat = statistic.New(bytes_of_memory, tcp_port, udp_port, max_connections, verbosity, cas, flush)
	server.Logger = NewServerLogger(verbosity)
	return server
}

// Public function runs loops with all available protocols
func (server *server) RunServer() {
	server.sockets = make(map[string] net.Listener)
	server.ThreadSync.Add(1)
	go server.run("tcp")
	if len(server.udp_port) > 0 {
		server.ThreadSync.Add(1)
		go server.run("udp")
	}
}

// Public function receives the pointer to server structure, stops the server and inform about it.
func (server *server) StopServer() {
	server.stop()
	server.Logger.Info("Server is now stopped.")
}
