package stat

import (
	"testing"
	"time"
	"os"
	"tools/cache"
	"runtime"
)

func TestNewServerStat(t *testing.T){
	stats := New(42)
	if stats.pid != os.Getpid() || stats.init_ts != time.Now().Unix() || stats.limit_maxbytes != 42 ||
	   stats.Current_connections != 0 || stats.Total_connections != 0 || stats.Read_bytes != 0 ||
	   stats.Written_bytes != 0 {
		t.Fatalf("Unexpected initialization of ServerStat struct: ", stats)
	}
}

func TestUptime(t *testing.T) {
	stats := New(42)
	time.Sleep(time.Second)
	uptime := stats.uptime()
	if uptime != 1 {
		t.Fatalf("Unexpected returned value: %d", uptime)
	}
}

func TestTime(t *testing.T) {
	stats := New(42)
	_time := stats.time()
	if _time != uint32(time.Now().Unix()) {
		t.Fatalf("Unexpected returned value: %d", _time)
	}
}

func TestBytesAmount(t *testing.T){
	stats := New(42)
	if stats.bytes(40) != 2 {
		t.Fatalf("Unexpected behavior")
	}
}

func TestSerializationCase1(t *testing.T){
	stats := New(42)
	storage := cache.New(42)
	if len(stats.Serialize(storage)) != 18 {
		t.Fatalf("Unexpected number of fields of returned value: ", stats)
	}
}

func TestSerializationCase2(t *testing.T){
	stats := New(42)
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
	if len(stats.Serialize(storage)) != 35 {
		t.Fatalf("Unexpected number of fields of returned value: ", stats)
	}
}

func TestNumberOfThreads(t *testing.T){
	if New(42).threads() != runtime.NumGoroutine() {
		t.Fatalf("Unexpected number of threads")
	}
}

func TestRusage(t *testing.T){
	var x ServerStat
	a, b, c, d := x.rusage()
	if a == 0 && b == 0 || b >= 1000000 || c == 0 && d == 0 || d >= 1000000 {
		t.Fatalf("Unexpected returned values")
	}
}
