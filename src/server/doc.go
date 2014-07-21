/*
The package implements memory cache server with same functionality as "memcached" project.

The basic usage:

import "server"

server.RunServer(port string, bytes_of_memory int64) - will run server on localhost:port,
and allow to keep data, which volume less or equal than 2nd param (bytes_of_memory).
This function will return a pointer on "server" struct, with public method StopServer, which gives you access to finish
work with server.
*/
package server
