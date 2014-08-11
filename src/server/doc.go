/*
Server package implements core of memory caching server, which can listen TCP (TODO: UDP and Unix sockets)
connections, and handle them according to ascii (TODO: binary)
protocol.

Server can be initialized by NewServer function (see below) and thus can be run with RunServer method of type Server.

If it is going to be necessary, server can be stopped with StopServer method,
but by design MemoranGo doesn't required to stop server manually, as it finish its job at the same time as process dies.

Server is available to await, for all goroutines will finish their jobs, it is provided with Wait method.   CAUTION:
internal consistence of goroutines was built such as they will finish their jobs ONLY when socket is undefined. Thus random usage of .Wait() may lock your process.

Server also supported with a logger (ServerLogger) and statistics (ServerStat).
Logger has few levels and depends from verbosity. It also has system logger inside, which doesn't look at verbosity, but logs only errors.
Statistics keeps all actions of server, described into Memcached specification.

Since MemoranGo is Go lang reimplemetation of memcached, Stat hasn't all fields from the Memcached, but also has its own ones.

*/
package server
