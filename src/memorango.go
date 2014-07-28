package main

import (
	"fmt"
	"flag"
	"sync"
	"server"
	"tools"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	port := flag.String("p", "11211", "Port to listen (non required - default port is 11211)")
	memory_amount_mb := flag.Int("m", 0, "Amount of memory to allocate (Mb)")
	daemonize := flag.Bool("d", false, "Run process as background")

	flag.Parse()
	if *memory_amount_mb <= 0 {
		fmt.Println("Impossible to run server with incorrect specified amount of available data.")
		return
	}
	if *daemonize {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Printf("Run MemoranGo daemon at 127.0.0.1:%s with %d mb allocated memory.\n", *port, *memory_amount_mb)
		fmt.Println("Status: ",
					exec.Command("/usr/bin/nohup", filepath.Join(dir, os.Args[0]), "-p", *port, "-m",
								 tools.IntToString(int64(*memory_amount_mb)), "&").Start())
	} else {
		fmt.Printf("%d Run MemoranGo on 127.0.0.1:%s with %d mb allocated memory.\n",
			os.Getpid(), *port, *memory_amount_mb)
		_server := server.RunServer(*port, int64(*memory_amount_mb)*1024*1024/* let's convert to bytes */)
		defer server.StopServer(_server)
		var w sync.WaitGroup
		w.Add(1)
		w.Wait()
	}
}
