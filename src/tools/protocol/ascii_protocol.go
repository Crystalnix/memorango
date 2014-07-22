package protocol

import (
	"tools"
	"strings"
	"errors"
)

// Templates for responses to a client.
const (
	ERROR_TEMP = "ERROR\r\n"
	CLIENT_ERROR_TEMP = "CLIENT_ERROR %s\r\n"
	SERVER_ERROR_TEMP = "SERVER_ERROR %s\r\n"
	NOT_FOUND = "NOT_FOUND\r\n"
)

// Specified groups of commands, which are helpful for destination handling of request.
var storage_commands = []string{"set", "add", "replace", "append", "prepend", "cas",}
var retrieve_commands = []string{"get", "gets",}
var other_commands = []string{"delete", "touch", "flush_all", "version", "quit",}

// Enumeration of protocol headers.
type Ascii_protocol_enum struct {
	command string		// the main action of the passed request.
	key []string		// key or keys for requested items.
	flags int 			// 32 or 16 bit int that server stores along with the data and sends back when the item is retrieved.
	exptime int64 		// UNIX timestamp, or time from now, which guaranties, that data won't be retrieved after this time.
	bytes int 			// the number of bytes in the data block to follow, NOT including the delimiting \r\n.
	cas_unique int64	// unique 64-bit value of an existing entry.
	noreply bool		// optional parameter instructs the server to not send the reply.
	data_string []byte	// chunk of arbitrary 8-bit data of length <bytes>
	error string		// error, which appears when something goes wrong, normally is empty string ""
}

// Public function, which parse string of input data to tokens of protocol's header and join them into one enumeration.
// Function returns pointer to Ascii_protocol_enum struct with nil value of error field if parsing succeeded.
// Otherwise error field consists information about occurred error and other fields are empty.
func ParseProtocolHeader(header string) *Ascii_protocol_enum{
	if len(header) > 0 {
		command_line := strings.Split(header, " ")
		command := command_line[0]
		switch true {
		case tools.In(command, storage_commands):
			return parseStorageCommands(command_line)
		case tools.In(command, retrieve_commands):
			return parseRetrieveCommands(command_line)
		case tools.In(command, other_commands):
			return parseOtherCommands(command_line)
		default:
			return &Ascii_protocol_enum{error: ERROR_TEMP}
		}
	} else {
		return &Ascii_protocol_enum{error: strings.Replace(CLIENT_ERROR_TEMP, "%s", "Input command line is empty", 1)}
	}
}

// Function for parsing of storage group of commands.
// Receives array of string tokens.
// Returns pointer to Ascii_protocol_enum with bound fields.
func parseStorageCommands(args []string/*, data_block string*/) *Ascii_protocol_enum{
	protocol := new(Ascii_protocol_enum)
	if len(args) < 5 || /*len(data_block) == 0 ||*/ tools.In("", args) {
		return &Ascii_protocol_enum{error: ERROR_TEMP}
	}
	var err error
	//protocol.data_string = []byte(data_block)
	protocol.command = args[0]
	protocol.key = []string{args[1],}
	protocol.flags, err = tools.StringToInt32(args[2])
	protocol.exptime, err = tools.StringToInt64(args[3])
	protocol.exptime = tools.ToTimeStampFromNow(protocol.exptime)
	protocol.bytes, err = tools.StringToInt32(args[4])
	if args[0] == "cas" {
		if len(args) < 6 {
			err = errors.New("invalid arguments number")
		}
		protocol.cas_unique, err = tools.StringToInt64(args[5])
		if len(args) == 7 {
			protocol.noreply = (args[6] == "noreply")
		}
	} else {
		if len(args) == 6 {
			protocol.noreply = (args[5] == "noreply")
		}
	}
	if err != nil {
		return &Ascii_protocol_enum{error: ERROR_TEMP}
	}
	return protocol
}

// Function for parsing of retrieving group of commands.
// Receives array of string tokens.
// Returns pointer to Ascii_protocol_enum with bound fields.
func parseRetrieveCommands(args []string) *Ascii_protocol_enum {
	protocol := new(Ascii_protocol_enum)
	if len(args) < 2 || tools.In("", args) {
		return &Ascii_protocol_enum{error: ERROR_TEMP}
	}
	protocol.command = args[0]
	protocol.noreply = false
	protocol.key = args[1 : ]
	return protocol
}

// Function for parsing of other group of commands.
// Receives array of string tokens.
// Returns pointer to Ascii_protocol_enum with bound fields.
func parseOtherCommands(args []string) *Ascii_protocol_enum {
	protocol := new(Ascii_protocol_enum)
	var err error
	err = nil
	if tools.In("", args) {
		err = errors.New("invalid arguments")
	}
	protocol.command = args[0]
	protocol.noreply = (args[len(args) - 1] == "noreply")
	switch args[0]{
	case "delete":
		if len(args) < 2 {
			err = errors.New("invalid arguments number")
		} else {
			protocol.key = []string{args[1], }
		}
	case "touch":
		if len(args) < 3{
			err = errors.New("invalid arguments number")
		} else {
			protocol.key = []string{args[1], }
			protocol.exptime, err = tools.StringToInt64(args[2])
		}
	case "flush_all":
		if len(args) >= 2 {
			protocol.exptime, err = tools.StringToInt64(args[1])
		}
	}
	protocol.exptime = tools.ToTimeStampFromNow(protocol.exptime)
	if err != nil {
		return &Ascii_protocol_enum{error: ERROR_TEMP}
	}
	return protocol
}
