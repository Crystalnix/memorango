package stat

import (
	//"os"
	//"os/exec"
	"time"
	"syscall"
	//"fmt"
	"tools/cache"
)

const (
	pointer_size = 32
)
// Structure contains all information provided by server.
// For extended info look at https://github.com/memcached/memcached/blob/master/doc/protocol.txt
// General-purpose statistic
type ServerStat struct {
	pid uint32
	init_ts int64
	current_items uint32
	total_items uint32
	current_connections uint32
	total_connections uint32
	read_bytes uint64
	written_bytes uint64
	limit_maxbytes uint32
	storage *cache.LRUCache
	Commands map[string] uint64
}

// TODO: Tests for everything.

// Function returns number of seconds since server has been run.
func (s *ServerStat) uptime() uint32 {
	return uint32(time.Now().Unix() - s.init_ts)
}

// Function returns current UNIX time according to the server.
func (s *ServerStat) time() uint32 {
	return uint32(time.Now().Unix())
}

// Function returns accumulated specified by owner time for this process (seconds, microseconds)
// Passed parameter owner could be 0 - user time or 1 - system time, otherwise function returns 0, 0
// Function returns 0, 0 if syscall returned an error as well.
func (s *ServerStat) rusage(owner int) (int64, int64) {
	//var container_string string
	rusage_struct := new(syscall.Rusage)
	if syscall.Getrusage(0, rusage_struct) != nil{
		return 0, 0
	}
	var (
		sec int64
		nsec int64
	)
	if owner == 0 {
		sec, nsec = rusage_struct.Utime.Unix()
	}
	if owner == 1 {
		sec, nsec = rusage_struct.Stime.Unix()
	}
	return sec, nsec / 1000
}

// Function returns amount of used bytes to store items.
func (s *ServerStat) bytes() int64 {
	return int64(s.limit_maxbytes) - s.storage.Capacity()
}

f
