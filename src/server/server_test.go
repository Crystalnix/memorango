package server

import (
	"testing"
	"net"
	"fmt"
	"time"
	"bytes"
	"os"
	"bufio"
)

var test_port = "60000"
var test_address = "127.0.0.1:" + test_port

func TestServerRunAndStop(t *testing.T){
//	var test_port = "60000"
//	var test_address = "127.0.0.1:" + test_port
	srv := NewServer(test_port, "", "", 1024, false, false, 0, 1024)
	srv.RunServer()
	time.Sleep(time.Millisecond * time.Duration(10)) // Let's wait a bit while goroutines will start
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
	srv.StopServer()
	connection.Close()
	connection, err = net.Dial("tcp", test_address)
	if err == nil {
		t.Fatalf("Server is still running at %s", test_address)
	}
}

func TestServerConsistenceAndConnections(t *testing.T){
//	var test_port = "60001"
//	var test_address = "127.0.0.1:" + test_port
	srv := NewServer(test_port, "", "", 1024, false, false, 0, 1024)
	srv.RunServer()
	time.Sleep(time.Millisecond * time.Duration(10))
	defer srv.StopServer()
	connection, err := net.Dial("tcp", test_address)
	// notice, that connection can be opened in case of failure.
	if err != nil {
		t.Fatalf("Server wasn't run: %s", err)
	}
	if srv.tcp_port != test_port || srv.sockets == nil || srv.storage == nil {
		t.Fatalf("Unexpected consistence: %s, %s, %s", srv.tcp_port, srv.sockets, srv.storage)
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
}

func TestServerResponseAndConnections(t *testing.T){
//	var test_port = "60002"
//	var test_address = "127.0.0.1:" + test_port
	srv := NewServer(test_port, "", "", 1024, false, false, 0, 1024)
	srv.RunServer()
	time.Sleep(time.Millisecond * time.Duration(10))
	defer srv.StopServer()
	connection, err := net.Dial("tcp", test_address)
	var test_msg = []byte("Test1\r\n")
	_, err = connection.Write(test_msg)
	if err != nil {
		t.Fatalf("Stream is unavailable to transmit data: ", err)
	}
	remote_connection := srv.connections[connection.LocalAddr().String()]
	if !srv.makeResponse(remote_connection, []byte("TestResponse"), 12){
		t.Fatalf("Server is unavailable to make response.")
	}
	var test_response = make([]byte, 12)
	n, err := connection.Read(test_response[0 : ])
	if n != 12 || err != nil {
		t.Fatalf("Connection stream is empty: ", n, err)
	}
	if string(test_response) != "TestResponse" {
		t.Fatalf("Unexpected response:< %s >", string(test_response))
	}
	if !srv.breakConnection(remote_connection){
		t.Fatalf("Server is unavailable to break connection at %s", remote_connection.RemoteAddr().String())
	}
	if len(srv.connections) != 0 {
		t.Fatalf("Connection is still alive: ", srv.connections)
	}
}

func TestServerReader1(t *testing.T){
	var test_msg = []byte("TEST\r\nwith-\r\n-terminators\r\n")
	var byteBuf = bytes.NewBuffer(test_msg)
	reader := bufio.NewReader(byteBuf)
	res, n, err := readRequest(reader, -1)
	if err != nil {
		t.Fatalf("Unexpected behaviour: ", err, res, n)
	}
	if string(res[0 : n - 2]) != "TEST" {
		t.Fatalf("Unexpected response: %s", string(res))
	}
	res, n, err = readRequest(reader, 19)
	if err != nil {
		t.Fatalf("Unexpected behaviour: ", err, res, n)
	}
	if string(res[0 : n - 2]) != "with-\r\n-terminators" {
		t.Fatalf("Unexpected response: %s", string(res))
	}

}

func TestServerReader2(t *testing.T){
	var test_msg = make([]byte, 300)
	var byteBuf = bytes.NewBuffer(test_msg)
	reader := bufio.NewReader(byteBuf)
	res, n, err := readRequest(reader, -1)
	if err == nil {
		t.Fatalf("Unexpected behaviour: ", res, n)
	}
	byteBuf = bytes.NewBuffer(test_msg)
	reader = bufio.NewReader(byteBuf)
	res, n, err = readRequest(reader, 42)
	if err == nil {
		t.Fatalf("Unexpected behaviour: ", res, n)
	}
	test_msg[298] = '\r'
	test_msg[299] = '\n'
	byteBuf = bytes.NewBuffer(test_msg)
	reader = bufio.NewReader(byteBuf)
	res, n, err = readRequest(reader, 298)
	if err != nil {
		t.Fatalf("Unexpected behaviour: ", err, res, n)
	}
}

//TODO: out of order yet.
func TestServerLogger(t *testing.T){
	var buf = make([]byte, 20)
	logger := NewServerLogger(0)
	logger.Error("errortest")
	logger.Warning("warningtest")
	logger.Info("infotest")
	n, _ := os.Stdout.Read(buf)
	if n != 0 {
		t.Fatalf("Unexpected logger behavior", string(buf), n)
	}
	n, _ = os.Stderr.Read(buf)
	if n == 0 || string(buf) != "errortest" {
		t.Fatalf("Unexpected logger behavior", string(buf), n)
	}

	logger = NewServerLogger(1)
	logger.Error("errortest")
	logger.Warning("warningtest")
	logger.Info("infotest")

	n, _ = os.Stdout.Read(buf)
	if n == 0 || string(buf) != "warningtest" {
		t.Fatalf("Unexpected logger behavior", string(buf), n)
	}
	n, _ = os.Stderr.Read(buf)
	if n == 0 || string(buf) != "errortest" {
		t.Fatalf("Unexpected logger behavior", string(buf), n)
	}
}
