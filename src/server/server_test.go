package server

import (
	"testing"
	"net"
	"fmt"
	"time"
	"bytes"
	"os"
	"bufio"
	"log"
	"tools/protocol"
)

var test_port = "60000"
var test_address = "127.0.0.1:" + test_port

func TestServerRunAndStop(t *testing.T){
	fmt.Println("TestServerRunAndStop")
//	var test_port = "60000"
//	var test_address = "127.0.0.1:" + test_port
	srv := NewServer(test_port, "", "", 1024, false, false, 2, 1024)
	srv.RunServer()
	defer srv.StopServer() // if test will be broken down before server will be closed.
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
	fmt.Println("TestServerConsistenceAndConnections")
//	var test_port = "60001"
//	var test_address = "127.0.0.1:" + test_port
	srv := NewServer(test_port, "", "", 1024, false, false, 2, 1024)
	srv.RunServer()
	time.Sleep(time.Millisecond * time.Duration(10))
	defer srv.StopServer()
	connection, err := net.Dial("tcp", test_address)
	// notice, that connection can be opened in case of failure.
	if err != nil {
		t.Fatalf("Server wasn't run: %s", err)
	}
	if srv.tcp_port != test_port || srv.sockets == nil || len(srv.sockets) != 1 || srv.storage == nil {
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
	connection.Close()
}

func TestServerResponseAndConnections(t *testing.T){
	fmt.Println("TestServerResponseAndConnections")
//	var test_port = "60002"
//	var test_address = "127.0.0.1:" + test_port
	srv := NewServer(test_port, "", "", 1024, false, false, 2, 1024)
	srv.RunServer()
	defer srv.StopServer()
	time.Sleep(time.Millisecond * time.Duration(10))
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
	connection.Close()

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


//TODO: To figure out how to check consistence of std outputs, or at least last n symbols of it.
func testReadOsOutput(output *os.File, offset int) []byte{
	var buf []byte
	log.Println(output.Seek(int64(offset), 2))
	log.Println(output.Read(buf))
	return buf
}

func TestServerLoggerTestSuite1(t *testing.T){
	logger0 := NewServerLogger(0)
	if logger0.error.Flags() != log.Ldate | log.Ltime | log.Lshortfile ||
	   logger0.warning.Flags() != 0 || logger0.info.Flags() != 0 ||
	   logger0.syslogger.Flags() != log.Ldate | log.Ltime | log.Lshortfile {
		t.Fatalf("Wrong flags of logger components:\nerr, warn, inf, sys\n%d, %d, %d, %d\n%d, %d, %d, %d\n",
			     logger0.error.Flags(), logger0.warning.Flags(), logger0.info.Flags(), logger0.syslogger.Flags(),
				 log.Ldate | log.Ltime | log.Lshortfile, 0, 0, log.Ldate | log.Ltime | log.Lshortfile)
	}
//	var buf []byte
//	logger0.Error("TestError")
//	buf = testReadOsOutput(os.Stderr, 9)
//	if string(buf) != "TestError" {
//		t.Fatalf("Unexpected message: %s ;", string(buf))
//	}
//	logger0.Warning("TestWarning")
//	buf = testReadOsOutput(os.Stdout, 11)
//	if string(buf) == "TestWarning" {
//		t.Fatalf("Unexpected logger behavior: depth permission corrupted.")
//	}
//	logger0.Info("TestInfo")
//	buf = testReadOsOutput(os.Stdout, 8)
//	if string(buf) == "TestInfo" {
//		t.Fatalf("Unexpected logger behavior: depth permission corrupted.")
//	}
}

func TestServerLoggerTestSuite2(t *testing.T){
	logger0 := NewServerLogger(1)
	if logger0.error.Flags() != log.Ldate | log.Ltime | log.Lshortfile ||
		logger0.warning.Flags() != log.Ldate | log.Ltime | log.Lshortfile || logger0.info.Flags() != 0 ||
		logger0.syslogger.Flags() != log.Ldate | log.Ltime | log.Lshortfile {
		t.Fatalf("Wrong flags of logger components:\nerr, warn, inf, sys\n%d, %d, %d, %d\n%d, %d, %d, %d\n",
			logger0.error.Flags(), logger0.warning.Flags(), logger0.info.Flags(), logger0.syslogger.Flags(),
					log.Ldate | log.Ltime | log.Lshortfile,
					log.Ldate | log.Ltime | log.Lshortfile, 0, log.Ldate | log.Ltime | log.Lshortfile)
	}
//	var buf []byte
//	logger0.Error("TestError")
//	buf = testReadOsOutput(os.Stderr, 9)
//	if string(buf) != "TestError" {
//		t.Fatalf("Unexpected message: %s ;", string(buf))
//	}
//	logger0.Warning("TestWarning")
//	buf = testReadOsOutput(os.Stdout, 11)
//	if string(buf) != "TestWarning" {
//		t.Fatalf("Unexpected logger behavior: depth permission corrupted.")
//	}
//	logger0.Info("TestInfo")
//	buf = testReadOsOutput(os.Stdout, 8)
//	if string(buf) == "TestInfo" {
//		t.Fatalf("Unexpected logger behavior: depth permission corrupted.")
//	}
}

func TestServerLoggerTestSuite3(t *testing.T){
	logger0 := NewServerLogger(2)
	if logger0.error.Flags() != log.Ldate | log.Ltime | log.Lshortfile ||
		logger0.warning.Flags() != log.Ldate | log.Ltime | log.Lshortfile ||
		logger0.info.Flags() != log.Ldate | log.Ltime ||
		logger0.syslogger.Flags() != log.Ldate | log.Ltime | log.Lshortfile {
		t.Fatalf("Wrong flags of logger components:\nerr, warn, inf, sys\n%d, %d, %d, %d\n%d, %d, %d, %d\n",
			logger0.error.Flags(), logger0.warning.Flags(), logger0.info.Flags(), logger0.syslogger.Flags(),
					log.Ldate | log.Ltime | log.Lshortfile,
					log.Ldate | log.Ltime | log.Lshortfile, log.Ldate | log.Ltime, log.Ldate | log.Ltime | log.Lshortfile)
	}
//	var buf []byte
//	logger0.Error("TestError")
//	buf = testReadOsOutput(os.Stderr, 9)
//	if string(buf) != "TestError" {
//		t.Fatalf("Unexpected message: %s ;", string(buf))
//	}
//	logger0.Warning("TestWarning")
//	buf = testReadOsOutput(os.Stdout, 11)
//	if string(buf) != "TestWarning" {
//		t.Fatalf("Unexpected logger behavior: depth permission corrupted.")
//	}
//	logger0.Info("TestInfo")
//	buf = testReadOsOutput(os.Stdout, 8)
//	if string(buf) != "TestInfo" {
//		t.Fatalf("Unexpected logger behavior: depth permission corrupted.")
//	}
}

func TestServerConnectionsLim(t *testing.T){
	fmt.Println("TestServerConnectionsLim")
	srv := NewServer(test_port, "", "", 2, false, false, 2, 1024)
	srv.RunServer()
	defer srv.StopServer()
	time.Sleep(time.Millisecond * time.Duration(10)) // Let's wait a bit while goroutines will start
	connection1, err := net.Dial("tcp", test_address)
	//defer connection1.Close()
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	_, err = connection1.Write([]byte("TEST"))
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	connection2, err := net.Dial("tcp", test_address)
	//defer connection2.Close()
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	_, err = connection2.Write([]byte("TEST"))
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	excess_conn, err := net.Dial("tcp", test_address)
	//defer excess_conn.Close()
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	_, err = excess_conn.Write([]byte("TEST"))
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	test_response := make([]byte, 10)
	n, err := excess_conn.Read(test_response[ : ])
	if err == nil || n != 0 {
		t.Fatalf("Unexpected behavior: connection shouldn't be handled.")
	}
	connection1.Close()
	connection2.Close()
	excess_conn.Close()
}

func TestServerConnectionListener(t *testing.T){
	fmt.Println("TestServerConnectionListener")
	srv := NewServer(test_port, "", "123.45.67.89", 1024, false, false, 2, 1024)
	srv.RunServer()
	defer srv.StopServer()
	time.Sleep(time.Millisecond * time.Duration(10)) // Let's wait a bit while goroutines will start
	conn, err := net.Dial("tcp", test_address)
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	_, err = conn.Write([]byte("You'll never see me"))
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	test_response := make([]byte, 10)
	n, err := conn.Read(test_response[ : ])
	if err == nil || n != 0 || string(test_response) == "ERROR\r\n"{
		t.Fatalf("Unexpected behavior: connection shouldn't be handled.")
	}
	conn.Close()
}

func TestServerForbiddenCmds(t *testing.T){
	fmt.Println("TestServerForbiddenCmds")
	srv := NewServer(test_port, "", "", 1024, true, true, 2, 1024)
	srv.RunServer()
	defer srv.StopServer()
	time.Sleep(time.Millisecond * time.Duration(10)) // Let's wait a bit while goroutines will start
	var test_response = make([]byte, 42)
	connection, err := net.Dial("tcp", test_address)
	//defer connection.Close()
	if err != nil{
		t.Fatalf("Unexpected server behavior.", err)
	}
	_, err = connection.Write([]byte("cas key 0 0 4 424242\r\nTEST\r\n"))
	if err != nil {
		t.Fatalf("Unexpected server behavior.", err)
	}
	n, err := connection.Read(test_response[ : ])
	if err != nil || n == 0 {
		t.Fatalf("Unexpected server behavior.", err)
	}
	if string(test_response) == protocol.NOT_FOUND {
		t.Fatalf("Unexpected server behavior: forbidden command CAS was handled.")
	}
	_, err = connection.Write([]byte("flush_all\r\n"))
	if err != nil {
		t.Fatalf("Unexpected server behavior.", err)
	}
	n, err = connection.Read(test_response[ : ])
	if err != nil || n == 0 {
		t.Fatalf("Unexpected server behavior.", err)
	}
	if string(test_response) == "OK\r\n" {
		t.Fatalf("Unexpected server behavior: forbidden command FLUSH_ALL was handled.")
	}

	_, err = connection.Write([]byte("gets key\r\n"))
	if err != nil {
		t.Fatalf("Unexpected server behavior.", err)
	}
	n, err = connection.Read(test_response[ : ])
	if err != nil || n == 0 {
		t.Fatalf("Unexpected server behavior.", err)
	}
	if string(test_response) == "END\r\n" {
		t.Fatalf("Unexpected server behavior: forbidden command GETS/CAS was handled.")
	}
	connection.Close()
}

func TestServerWaitSuite(t *testing.T) {
	srv := NewServer(test_port, "", "", 1024, false, false, 2, 1024)
	start := time.Now().UnixNano()
	srv.threads = 0
	srv.Wait()
	end := time.Now().UnixNano()
	if end-start > 1000 {
		t.Fatalf("Unexpected waiting behaviour: wait had to finish immidiatelly: %d.", end-start)
	}
}
