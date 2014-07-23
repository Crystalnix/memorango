package server

import (
	"testing"
	"net"
	"fmt"
	"time"
	"bytes"
	"bufio"
)

const (
	test_port = "60001"
	test_address = "127.0.0.1:" + test_port
)

func TestServerRunAndStop(t *testing.T){
	srv := RunServer(test_port, 1024)
	time.Sleep(time.Millisecond * time.Duration(5))
	connection, err := net.Dial("tcp", test_address)
	if err != nil {
		t.Fatalf("Server wasn't run: %s", err)
	}
	var test_msg = []byte("Test1\r\n")
	_, err = connection.Write(test_msg)
	if err != nil {
		t.Fatalf("Stream is unavailable to transmit data: ", err)
	}
	var response = make([]byte, 255)
	fmt.Println("Trying to read response...")
	_, err = connection.Read(response[0: ])
	if err != nil {
		t.Fatalf("Stream is unavailable to transmit data: ", err)
	}
	if len(response) == 0 {
		t.Fatalf("Server doesn't response.")
	}
	fmt.Println("Response: ", string(response))
	StopServer(srv)
	connection.Close()
	connection, err = net.Dial("tcp", test_address)
	if err == nil {
		t.Fatalf("Server is still running at %s", test_address)
	}
}

func TestServerConsistenceAndConnections(t *testing.T){
	srv := RunServer(test_port, 1024)
	time.Sleep(time.Millisecond * time.Duration(10))
	defer StopServer(srv)
	connection, err := net.Dial("tcp", test_address)
	// notice, that connection can be opened in case of failure.
	if err != nil {
		t.Fatalf("Server wasn't run: %s", err)
	}
	if srv.port != test_port || srv.socket == nil || srv.storage == nil {
		t.Fatalf("Unexpected consistence: %s, %s, %s", srv.port, srv.socket, srv.storage)
	}
	var test_msg = []byte("Test1\r\n")
	_, err = connection.Write(test_msg)
	if err != nil {
		t.Fatalf("Stream is unavailable to transmit data: ", err)
	}
	remote_connection := srv.connections[connection.LocalAddr().String()]
	if len(srv.connections) != 1 || remote_connection == nil {
		t.Fatalf("Connection wasn't cached: ", len(srv.connections))
	}

	if remote_connection.LocalAddr().String() != connection.RemoteAddr().String() ||
	   remote_connection.RemoteAddr().String() != connection.LocalAddr().String() {
		t.Fatalf("Connections mismatch")
	}

	srv.breakConnection(remote_connection)
	if len(srv.connections) != 0 {
		t.Fatalf("Connection is still acitve.")
	}
}

func TestServerReader(t *testing.T){
	var test_msg = []byte("TEST\r\nwith-\r\n-terminators\r\n")
	var byteBuf = bytes.NewBuffer(test_msg)
	reader := bufio.NewReader(byteBuf)
	res, n, err := readRequest(reader, -1)
	if err != nil {
		t.Fatalf("Unexpected behaviour: ", err)
	}
	if string(res[0 : n - 2]) != "TEST" {
		t.Fatalf("Unexpected response: %s", string(res))
	}
	res, n, err = readRequest(reader, 19)
	if err != nil {
		t.Fatalf("Unexpected behaviour: ", err)
	}
	if string(res[0 : n - 2]) != "with-\r\n-terminators" {
		t.Fatalf("Unexpected response: %s", string(res))
	}
	fmt.Println("Success")
}
