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
	verbosity int
	tcp string
	udp string
	cas_disabled bool
	flush_disabled bool
	Connections map[string] ConnectionStat
	Current_connections uint32
	Total_connections uint32
	Connections_limit int
	Read_bytes uint64
	Written_bytes uint64
	Commands map[string] uint64
}

// Structure for logging statistics for connections bound with server.
type ConnectionStat struct {
	State string
	//The number of seconds since the most recently
	//issued command on the connection. This measures
	//the time since the start of the command, so if
	//"state" indicates a command is currently executing,
	//this will be the number of seconds the current
	//command has been running.
	//state string
	Cmd_hit_ts int64
	//The current state of the connection.
	Addr string
	//The address of the remote side. For listening
	//sockets this is the listen address. Note that some
	//socket types (such as UNIX-domain) don't have
	//meaningful remote addresses.
}


// Function initialize ServerStat structure by required start parameters.
func New(memory_amount int64, tcp_port string, udp_port string,
	     conn_max int, verbosity int, cas_disabled bool, flush_disabled bool) *ServerStat {
	return &ServerStat {
		pid: os.Getpid(),
		init_ts: time.Now().Unix(),
		limit_maxbytes: memory_amount,
		connections: make(map[string] ConnectionStat),
		tcp: tcp_port,
		udp: udp_port,
		verbosity: verbosity,
		cas_disabled: cas_disabled,
		flush_disabled: flush_disabled,
		Connections_limit: conn_max,
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
	dict["goroutines"] = tools.IntToString(int64(runtime.NumGoroutine()))
	for key, value := range s.Commands {
		dict[key] = tools.IntToString(int64(value))
	}
	return dict
}

// Function serialize sub command of stats "settings"
func (s *ServerStat) Settings(storage *cache.LRUCache) map[string] string {
	dict := make(map[string] string)
	dict["maxbytes"] = tools.IntToString(storage.Capacity())
	dict["maxconns"] = tools.IntToString(int64(s.Connections_limit))
	dict["tcpport"] = s.tcp
	dict["udpport"] = s.udp
	dict["verbosity"] = tools.IntToString(int64(s.verbosity))
	dict["num_goroutines"] = tools.IntToString(int64(runtime.NumGoroutine()))
	dict["evictions"] = "on" //TODO: to think about apportunity of another value.
	if storage.Crawler.Enabled() {
		dict["lru_crawler"] = "true"
	} else {
		dict["lru_crawler"] = "false"
	}
	dict["lru_crawler_sleep"] = tools.IntToString(int64(storage.Crawler.Sleep()))
	dict["lru_crawler_tocrawl"] = tools.IntToString(int64(storage.Crawler.ItemsPerRun))
	if s.cas_disabled {
		dict["cas_enabled"] = "false"
	} else {
		dict["cas_enabled"] = "true"
	}
	if s.flush_disabled {
		dict["flush_all_enabled"] = "false"
	} else {
		dict["flush_all_enabled"] = "true"
	}
	return dict
}

// Function serialize sub command of stats "conns"
func (s *ServerStat) Conns() map[string] string {
	dict := make(map[string] string)
	return dict
}

// Constructor for connection statistic
func NewConnStat(connection_addr string) *ConnectionStat {
	return &ConnectionStat {
		State: "conn_waiting",
		Sec_since_last_cmd: 0,
		Addr: connection_addr,
	}
}
