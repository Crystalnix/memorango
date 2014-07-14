package protocol

import (
	"server/tools/cache"
	"server/tools"
	"strings"
	"strconv"
)

const (
	error_temp = "ERROR\r\n"
	client_error_temp = "CLIENT_ERROR %s\r\n"
	server_error_temp = "SERVER_ERROR %s\r\n"
)

var storage_commands = []string{"set", "add", "replace", "append", "prepend", "cas",}
var retrieve_commands = []string{"get", "gets",}
var other_commands = []string{"delete", "touch", "flush_all", "version", "quite",}

type text_protocol struct {
	command string		// the main action of the passed request.
	key []string		// key or keys for requested items.
	flags int 			// 32 or 16 bit int that server stores along with the data and sends back when the item is retrieved.
	exptime int 		// UNIX timestamp, which guaranties, that data won't be retrieved after this time.
	bytes int 			// the number of bytes in the data block to follow, NOT including the delimiting \r\n.
	cas_unique int64	// unique 64-bit value of an existing entry.
	noreply bool		// optional parameter instructs the server to not send the reply.
	data_string []byte	// chunk of arbitrary 8-bit data of length <bytes>
	error string		// error, which appears when something goes wrong, normally is empty string ""
}

func (p *text_protocol) Set(instance_ptr *cache.LRUCache) string {
	//...
	return ""
}

func InitProtocol(request string) *text_protocol{
	parsed_req := strings.Split(request, "\r\n")
	if parsed_req[0] {
		command_line := strings.Split(parsed_req[0], " ")
		command := command_line[0]
		switch true {
		case tools.In(command, storage_commands):
			return parseStorageCommands(command_line, parsed_req[1])
		case tools.In(command, retrieve_commands):
			return parseRetrieveCommands(command_line)
		case tools.In(command, other_commands):
			return parseOtherCommands(command_line)
		default:
			return &text_protocol{error: error_temp}
		}
	} else {
		return &text_protocol{error: strings.Replace(client_error_temp, "%s", "Input command line is empty", 1)}
	}
}

func parseStorageCommands(args []string, data_block string) *text_protocol{
	protocol := new(text_protocol)
	if len(args) < 5 || len(data_block) == 0 || tools.In("", args) {
		return &text_protocol{error: error_temp}
	}
	protocol.command = args[0]
	protocol.key = []string{args[1],}
	protocol.flags = args[2]
	protocol.exptime = args[3]
	protocol.bytes = args[4]
	if args[0] == "cas" {
		if len(args) < 6 {
			return &text_protocol{error: error_temp}
		}
		protocol.cas_unique = args[5]
		if len(args) == 7 {
			protocol.noreply = (args[6] == "noreply")
		}
	} else {
		if len(args) == 6 {
			protocol.noreply = (args[5] == "noreply")
		}
	}
	return protocol
}

func parseRetrieveCommands(args []string) *text_protocol {
	protocol := new(text_protocol)
	if len(args) < 2 || tools.In("", args) {
		return &text_protocol{error: error_temp}
	}
	protocol.command = args[0]
	if args[len(args) - 1] != "noreply" {
		protocol.key = args[1 : ]
	} else {
		protocol.key = args[1 : len(args) - 2]
		protocol.noreply = true
	}
	return protocol
}

func parseOtherCommands(args []string) *text_protocol {
	protocol := new(text_protocol)
	if tools.In("", args) {
		return &text_protocol{error: error_temp}
	}
	protocol.command = args[0]
	switch args[0]{
	case "delete":
		if len(args) < 2{
			return &text_protocol{error: error_temp}
		}
		protocol.key = []string{args[1], }
	case "touch":
		if len(args) < 3{
			return &text_protocol{error: error_temp}
		}
		protocol.key = []string{args[1], }
		protocol.exptime = args[2]
	case "flush_all":
		if len(args) >= 2 {
			exp_time, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
				return &text_protocol{error: error_temp}
			}
			protocol.exptime(int(exp_time))
		}
	}
	return protocol
}
