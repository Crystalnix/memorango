package stat

import (
	"os"
	"time"
	"syscall"
	"runtime"
	"tools/cache"
	"tools"
)

const (
	pointer_size = 32
)
// Structure contains all information provided by server.
// For extended info look at https://github.com/memcached/memcached/blob/master/doc/protocol.txt
// General-purpose statistic
type ServerStat struct {
	pid int
	init_ts int64
	limit_maxbytes int64
	Current_connections uint32
	Total_connections uint32
	Read_bytes uint64
	Written_bytes uint64
	Commands map[string] uint64
}

// Function initialize ServerStat structure by required start parameters.
func New(memory_amount int64) *ServerStat {
	return &ServerStat{
		pid: os.Getpid(),
		init_ts: time.Now().Unix(),
		limit_maxbytes: memory_amount,
		Current_connections: 0,
		Total_connections: 0,
		Read_bytes: 0,
		Written_bytes: 0,
		Commands: make(map[string] uint64),
	}
}

// Function returns number of seconds since server has been run.
func (s *ServerStat) uptime() uint32 {
	return uint32(time.Now().Unix() - s.init_ts)
}

// Function returns current UNIX time according to the server.
func (s *ServerStat) time() uint32 {
	return uint32(time.Now().Unix())
}

// Function returns accumulated time for user and system for this process:
// (users_seconds, users_microseconds, systems_seconds, systems_microseconds)
// Function returns four zeroes if syscall returned an error.
func (s *ServerStat) rusage() (int64, int64, int64, int64) {
	//var container_string string
	rusage_struct := new(syscall.Rusage)
	if syscall.Getrusage(0, rusage_struct) != nil{
		return 0, 0, 0, 0
	}
	sec_u, nsec_u := rusage_struct.Utime.Unix()
	sec_s, nsec_s := rusage_struct.Stime.Unix()

	return sec_u, nsec_u / 1000, sec_s, nsec_s / 1000
}

// Function returns amount of used bytes to store items.
func (s *ServerStat) bytes(capacity int64) int64 {
	return int64(s.limit_maxbytes) - capacity
}

// Function returns number of active goroutines for current process
func (s *ServerStat) threads() int {
	return runtime.NumGoroutine()
}

// Function serialize statistic of server and storage and returns it as map of strings
func (s *ServerStat) Serialize(storage *cache.LRUCache) map[string] string {
	dict := make(map[string] string)
	dict["pid"] = tools.IntToString(int64(s.pid))
	dict["uptime"] = tools.IntToString(int64(s.uptime()))
	dict["time"] = tools.IntToString(int64(s.time()))
	dict["version"] = tools.VERSION
	dict["pointer_size"] = tools.IntToString(int64(pointer_size))
	secu, mcsecu, secs, mcsecs := s.rusage()
	dict["rusage_user"] = tools.IntToString(secu) + "." + tools.IntToString(mcsecu)
	dict["rusage_system"] = tools.IntToString(secs) + "." + tools.IntToString(mcsecs)
	dict["curr_items"] = tools.IntToString(int64(storage.Stats.Current_items))
	dict["total_items"] = tools.IntToString(int64(storage.Stats.Total_items))
	dict["bytes"] = tools.IntToString(s.bytes(storage.Capacity()))
	dict["curr_connections"] = tools.IntToString(int64(s.Current_connections))
	dict["total_connections"] = tools.IntToString(int64(s.Total_connections))
	dict["evictions"] = tools.IntToString(int64(storage.Stats.Evictions))
	dict["expired_unfetched"] = tools.IntToString(int64(storage.Stats.Expired_unfetched))
	dict["evicted_unfetched"] = tools.IntToString(int64(storage.Stats.Evicted_unfetched))
	dict["bytes_read"] = tools.IntToString(int64(s.Read_bytes))
	dict["bytes_written"] = tools.IntToString(int64(s.Written_bytes))
	dict["threads"] = tools.IntToString(int64(s.threads()))
	for key, value := range s.Commands {
		dict[key] = tools.IntToString(int64(value))
	}
	return dict
}
