/*
Package implements memcached plain text protocol (TODO: binary).
See more information at https://github.com/memcached/memcached/blob/master/doc/protocol.txt

ascii_protocol.go - describes rules of parsing and keeps data structures for plain text ascii protocol.
(TODO: binary_protocol.go, which is satisfied to memcached binary protocol + need to make an interface for protocols; However it has same parts;)
handling.go - describes rules of handling requests and making responses.
*/
package protocol

