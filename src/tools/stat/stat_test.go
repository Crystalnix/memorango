package stat

import (
	"testing"
	"time"
	"os"
	"tools/cache"
	"net"
	"fmt"
)

func TestNewServerStat(t *testing.T){
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	if stats.pid != os.Getpid() ||
		stats.init_ts != time.Now().Unix() ||
		stats.limit_maxbytes != 42 ||
	    stats.Current_connections != 0 ||
		stats.Total_connections != 0 ||
		stats.Read_bytes != 0 ||
	    stats.Written_bytes != 0 ||
		stats.tcp != "9999" ||
		stats.udp != "8888" ||
		stats.Connections_limit != 1024 ||
		stats.verbosity != 2 ||
		!stats.cas_disabled ||
		!stats.flush_disabled {
		t.Fatalf("Unexpected initialization of ServerStat struct: ", stats)
	}
}

func TestUptime(t *testing.T) {
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	time.Sleep(time.Second)
	uptime := stats.uptime()
	if uptime != 1 {
		t.Fatalf("Unexpected returned value: %d", uptime)
	}
}

func TestTime(t *testing.T) {
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	_time := stats.time()
	if _time != uint32(time.Now().Unix()) {
		t.Fatalf("Unexpected returned value: %d", _time)
	}
}

func TestBytesAmount(t *testing.T){
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	if stats.bytes(40) != 2 {
		t.Fatalf("Unexpected behavior")
	}
}

func TestSerializationCase1(t *testing.T){
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	storage := cache.New(42)
	if len(stats.Serialize(storage)) != 19 {
		t.Fatalf("Unexpected number of fields of returned value: ", stats)
	}
}

func TestSerializationCase2(t *testing.T){
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	storage := cache.New(42)
	stats.Commands["cmd_get"] ++
	stats.Commands["cmd_set"] ++
	stats.Commands["cmd_delete"] ++
	stats.Commands["cmd_touch"] ++
	stats.Commands["get_misses"] ++
	stats.Commands["get_hits"] ++
	stats.Commands["cas_misses"] ++
	stats.Commands["cas_hits"] ++
	stats.Commands["incr_misses"] ++
	stats.Commands["incr_hits"] ++
	stats.Commands["decr_misses"] ++
	stats.Commands["decr_hits"] ++
	stats.Commands["delete_misses"] ++
	stats.Commands["delete_hits"] ++
	stats.Commands["touch_misses"] ++
	stats.Commands["touch_hits"] ++
	stats.Commands["cas_badval"] ++
	if len(stats.Serialize(storage)) != 36 {
		t.Fatalf("Unexpected number of fields of returned value: ", stats)
	}
}

func TestRusage(t *testing.T){
	var x ServerStat
	a, b, c, d := x.rusage()
	if a == 0 && b == 0 || b >= 1000000 || c == 0 && d == 0 || d >= 1000000 {
		t.Fatalf("Unexpected returned values")
	}
}

func TestSettingsSerialization(t *testing.T){
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	storage := cache.New(42)
	res := stats.Settings(storage)
	if len(res) != 12 {
		t.Fatalf("Unexpected length of Settings serialization: %d, expected 12;", len(res), res)
	}
}

func TestConnectionsSerialization(t *testing.T) {
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	go func(){
		listener, _ := net.Listen("tcp", "127.0.0.1:9999")
		defer listener.Close()
		for {
			con, err := listener.Accept()
			if err != nil{
				break
			} else {
				if stats.Connections != nil && stats.Connections[con.RemoteAddr().String()] == nil {
					stats.Connections[con.RemoteAddr().String()] = NewConnStat(con)
				}
			}
		}
	}()
	time.Sleep(time.Millisecond * time.Duration(10)) // Let's wait a bit while goroutine will start
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	defer conn.Close()
	if err != nil {
		t.Fatalf("Unexpected error:", err)
	}
	fmt.Println(conn, err)
	if len(stats.Conns()) != 3 {
		t.Fatalf("Unexpected length of returned value; expected 3.")
	}
}

func TestConnectionsItems(t *testing.T) {
	stats := New(42, "9999", "8888", 1024, 2, true, true)
	storage := cache.New(42)
	if len(stats.Items(storage)) != 6 {
		t.Fatalf("Unexpected length of returned value; expected 6.")
	}
}
