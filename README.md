MemoranGo
=========

MemoranGo is reimplementation of [Memcahed project](https://github.com/memcached/memcached) into [Google Go](http://golang.org/)

Requirements
------------
The MemoranGo is supposed for UNIX-like operation systems.
So if you have any of one, feel free to use.

Installation
------------
* [Download](http://golang.org/dl/) and [install](http://golang.org/doc/install#install) Go compiler.
* Clone this project, or download as zip and unpack it.
* Open your terminal, cd to the project directory, add additional path environment for Go and build the binary file:
>> `cd /path/to/project/`   
>> `GOPATH=$GOPATH:$PWD`   
>> `go build src/memorango.go`   
* To make sure, that whole system works fine run tests:
>> `go test src/

TODO: need to remove all useless files and add test for memorango.go file.

* And build documentation:
>> `godoc -http=":6060" -goroot="src/"`   
    
   That is it. Now you are ready to run MemoranGo! 

Usage
-----
You can set path to binary file within environment PATH or simply run from current folder.    

**__Example:__**   
Run MemoranGo with specified flags `memorango -m 100 -p 10000` with 100 MiB on port 10000.   

MemoranGo can be used with following flags:   

* -p - TCP Port to listen (non required - default port is 11211)   
* -m - Amount of memory to allocate (MiB)   
* -d - Run process as background.   
* -l - Listen on specified ip addr only; default is any address.   
* -c - Use max simultaneous connections; default is 1024.   
* -U - UDP Port to listen (default is turned off)   
* -C - Disabling of cas command support.   
* -F - Disabling of flush_all command support.   
* -h - Show usage manual and list of options.   
* -v - Turning verbosity on. This option includes errors and warnings only.   
* -vv - Turning deep verbosity on. This option includes requests, responses and same output as simple verbosity.   

License
-------
License will be here

Contacts
--------
Contacts will be here
