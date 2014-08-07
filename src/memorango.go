package main

import (
	"fmt"
	"flag"
	//"sync"
	"server"
	"tools"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	tcp_port := flag.String("p", "11211", "TCP Port to listen (non required - default port is 11211)")
	memory_amount_mb := flag.Int("m", 0, "Amount of memory to allocate (MiB)")
	daemonize := flag.Bool("d", false, "Run process as background")
	// unix_socket := flag.String("s", "", "Unix socket path to listen on (disables network support)")
	// unix_perms := flag.String("a", "", "Permissions (in octal format) for Unix socket created with -s option")
	listen_ip := flag.String("l", "", "Listen on specified ip addr only; default to any address.")
	max_connections := flag.Int("c", 1024, "Use max simultaneous connections;")
	udp_port := flag.String("U", "", "UDP Port to listen (default is empty string - which means it is turned off)")
	disable_cas := flag.Bool("C", false, "Disabling of cas command support.")
	disable_flush := flag.Bool("F", false, "Disabling of flush_all command support.")
	help := flag.Bool("h", false, "Show usage manual and list of options.")
	verbose := flag.Bool("v", false, "Turning verbosity on. This option includes errors and warnings only.")
	deep_verbose := flag.Bool("vv", false, "Turning deep verbosity on. This option includes requests, responses and same output as simple verbosity.")
	flag.Parse()

	if *help {
		// TODO: It should be spread in future.
		fmt.Println("MemoranGo - memory caching service.\nusage:\nmemorango -m <memory_to_alloc> [-CvhFvvd]\n"+
				"\t[-l <listen_ip>] [-c <limit_connections>] [-p <tcp_port>] [-U <udp_port>]")
		return
	}

	if *memory_amount_mb <= 0 {
		fmt.Println("Impossible to run server with incorrect specified amount of available data.")
		return
	}

	var verbosity = 0
	if *deep_verbose {
		verbosity = 2
	} else if *verbose {
		verbosity = 1
	}


	if *daemonize {
		var transacted_options = []string{}
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		transacted_options = append(transacted_options, filepath.Join(dir, os.Args[0]), "-p", *tcp_port, "-m",
									tools.IntToString(int64(*memory_amount_mb)),
									"-c", tools.IntToString(int64(*max_connections)))
		if len(*listen_ip) > 0 {
			transacted_options = append(transacted_options, "-l", *listen_ip)
		}
		if len(*udp_port) > 0 {
			transacted_options = append(transacted_options, "-U", *udp_port)
		}
		if *disable_cas {
			transacted_options = append(transacted_options, "-C")
		}
		if *disable_flush {
			transacted_options = append(transacted_options, "-F")
		}
		if verbosity == 1 {
			transacted_options = append(transacted_options, "-v")
		} else if verbosity == 2 {
			transacted_options = append(transacted_options, "-vv")
		}
		transacted_options = append(transacted_options, "&")
		fmt.Printf("Run MemoranGo daemon at 127.0.0.1:%s with %d MiB allowed memory.\n", *tcp_port, *memory_amount_mb)
		cmd := exec.Command("/usr/bin/nohup", transacted_options...)
		start_err := cmd.Start()
		if start_err != nil {
			fmt.Println("Status: ", start_err)
		}
	} else {
		fmt.Printf("%d Run MemoranGo on 127.0.0.1:%s with %d MiB allowed memory.\n",
			os.Getpid(), *tcp_port, *memory_amount_mb)
		_server := server.NewServer(*tcp_port, *udp_port, *listen_ip, *max_connections, *disable_cas, *disable_flush,
									verbosity, int64(*memory_amount_mb)*1024*1024 /* let's convert to bytes */)
		_server.RunServer()
		defer _server.StopServer()
//		var w sync.WaitGroup
//		w.Add(1)
//		w.Wait()
		_server.ThreadSync.Wait()
	}
}
