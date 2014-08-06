package protocol

import (
	"tools/cache"
	"tools/stat"
	"tools"
	"strings"
	"errors"
)

// Public method of Ascii_protocol_enum operates with received storage: retrieves, discards, sets or updates items,
// related to own containment.
// Also, function receives stats structure, which possibly may be a nil. This structure serves for recording statistic
// of processing request.
// Returns response to client as byte-string and error/nil.
// If process was successful, there will be returned nil instead of error, otherwise it will be returned specified error.
func (enum *Ascii_protocol_enum) HandleRequest(storage *cache.LRUCache, stats *stat.ServerStat) ([]byte, error) {
	var err error
	if len(enum.error) > 0 {
		return []byte(enum.error), nil
	}
	var result string
	switch enum.command {
	case "set":
		result, err = enum.set(storage)
	case "cas":
		result, err = enum.cas(storage)
	case "add":
		result, err = enum.add(storage)
	case "replace":
		result, err = enum.replace(storage)
	case "append":
		result, err = enum.append(storage)
	case "prepend":
		result, err = enum.prepend(storage)
	case "incr":
		result, err = enum.fold(storage, 1)
	case "decr":
		result, err = enum.fold(storage, -1)
	case "get":
		result, err = enum.get(storage, false)
	case "gets":
		result, err = enum.get(storage, true)
	case "touch":
		result, err = enum.touch(storage)
	case "delete":
		result, err = enum.delete(storage)
	case "flush_all":
		result, err = enum.flush_all(storage)
	case "lru_crawler":
		return []byte(enum.lru_crawler(storage)), nil
	case "stats":
		if stats != nil {
			return []byte(enum.stat(storage, stats)), nil
		} else {
			return nil, errors.New("Statistic is not supported.")
		}
	case "version":
		return []byte("VERSION " + tools.VERSION + "\r\n"), nil
	case "quit":
		return nil, errors.New("Exit.")
	}
	if stats != nil {
		enum.RecordStats(stats, result)
	}
	return []byte(result), err
}

// Storage commands

// Implements set method
func (enum *Ascii_protocol_enum) set(storage *cache.LRUCache) (string, error){
	if storage.Set(tools.NewStoredData(enum.data_string, enum.key[0]), enum.flags, enum.exptime, 0) {
		return STORED, nil
	} else {
		return strings.Replace(SERVER_ERROR_TEMP, "%s", "Not enough memory", 1), errors.New("SERVER_ERROR")
	}
}

// Implements add method
func (enum *Ascii_protocol_enum) add(storage *cache.LRUCache) (string, error) {
	if storage.Get(enum.key[0]) != nil {
		return NOT_STORED, nil
	}
	return enum.set(storage)
}

// Utility method, for joining common parts of prepend/append methods.
// Receives additional parameters: existing_item - item retrieved from cache, uses for inheritance such params as flags,
// cas, exptime; pending_data - new concatenated data.
func (enum *Ascii_protocol_enum) pending(storage *cache.LRUCache,
										 existed_item *cache.LRUCacheItem, pending_data []byte) (string, error) {
	enum.bytes = len(pending_data)
	enum.SetData(pending_data)
	enum.exptime = existed_item.Exptime
	enum.cas_unique = existed_item.Cas_unique
	enum.flags = existed_item.Flags
	return enum.set(storage)
}

// Implements prepend method
func (enum *Ascii_protocol_enum) prepend(storage *cache.LRUCache) (string, error) {
	existed_item := storage.Get(enum.key[0])
	if existed_item == nil {
		return NOT_STORED, nil
	}
	existed_data := tools.ExtractStoredData(existed_item.Cacheable)
	if existed_data == nil {
		return NOT_STORED, nil
	}

	return enum.pending(storage, existed_item, append(enum.data_string, existed_data...)) // some kind of golang magic
}

// Implements append method
func (enum *Ascii_protocol_enum) append(storage *cache.LRUCache) (string, error) {
	existed_item := storage.Get(enum.key[0])
	if existed_item == nil {
		return NOT_STORED, nil
	}
	existed_data := tools.ExtractStoredData(existed_item.Cacheable)
	if existed_data == nil {
		return NOT_STORED, nil
	}
	return enum.pending(storage, existed_item, append(existed_data, enum.data_string...)) // ...
}

// Implements replace method
func (enum *Ascii_protocol_enum) replace(storage *cache.LRUCache) (string, error) {
	if storage.Get(enum.key[0]) == nil {
		return NOT_STORED, nil
	}
	return enum.set(storage)
}

// Implementation of Check And Set method
func (enum *Ascii_protocol_enum) cas(storage *cache.LRUCache) (string, error) {
	existed_item := storage.Get(enum.key[0])
	if existed_item != nil {
		if existed_item.Cas_unique == enum.cas_unique && existed_item.Cas_unique != 0{
			return enum.set(storage)
		}
	}
	return NOT_FOUND, nil
}

// Retrieving commands

// Implements get method
// Passed boolean param cas - defines of returning cas_unique
func (enum *Ascii_protocol_enum) get(storage *cache.LRUCache, cas bool) (string, error) {
	var result = ""
	for _, value := range enum.key{
		item := storage.Get(value)
		if item != nil {
			data := tools.ExtractStoredData(item.Cacheable)
			if data == nil {
				continue
			}
			result += "VALUE " + value + " " + tools.IntToString(int64(item.Flags)) + " " + tools.IntToString(int64(len(data)))
			if cas {
				cas_id := tools.GenerateCasId()
				storage.SetCas(value, cas_id)
				result += " " + tools.IntToString(cas_id)
			}
			result += "\r\n"
			result += string(data) + "\r\n"
		}
	}
	return result + "END\r\n", nil
}

// Other commands

// Implements touch method
func (enum *Ascii_protocol_enum) touch(storage *cache.LRUCache) (string, error) {
	if item := storage.Get(enum.key[0]); item == nil {
		return NOT_FOUND, nil
	} else {
		if !storage.Set(item.Cacheable, item.Flags, enum.exptime, item.Cas_unique){
			return strings.Replace(SERVER_ERROR_TEMP, "%s", "Not enough memory", 1), errors.New("SERVER_ERROR")
		}
		return "TOUCHED\r\n", nil
	}
}

// Implements delete method
func (enum *Ascii_protocol_enum) delete(storage *cache.LRUCache) (string, error) {
	if storage.Flush(enum.key[0]){
		return "DELETED\r\n", nil
	}
	return NOT_FOUND, nil
}

// Implements flush all method
func (enum *Ascii_protocol_enum) flush_all(storage *cache.LRUCache) (string, error) {
	storage.FlushAll()
	return "OK\r\n", nil
}

// Utility method, for joining common parts of incr/decr methods.
// Receives additional param sign, which defines operation: -1 or 1
func (enum *Ascii_protocol_enum) fold(storage *cache.LRUCache, sign int) (string, error) {
	if item := storage.Get(enum.key[0]); item != nil && (sign == 1 || sign == -1) {
		existed_data := tools.ExtractStoredData(item.Cacheable)
		if existed_data != nil {
			evaluated_data_for_existed, err_for_existed := tools.StringToInt64(string(existed_data))
			evaluated_data_for_passed, err_for_passed := tools.StringToUInt64(string(enum.data_string))
			if err_for_existed == nil && err_for_passed == nil {
				var result string
				if sign > 0 {
					result = tools.IntToString(evaluated_data_for_existed + int64(evaluated_data_for_passed))
				} else {
					result = tools.IntToString(evaluated_data_for_existed - int64(evaluated_data_for_passed))
				}
				if storage.Set(tools.NewStoredData([]byte(result), enum.key[0]), item.Flags, item.Exptime, 0) {
					return result+"\r\n", nil
				}
				return strings.Replace(SERVER_ERROR_TEMP, "%s", "Not enough memory", 1), errors.New("SERVER_ERROR")
			}
			return ERROR_TEMP, nil
		}
	}
	return NOT_FOUND, nil
}

// Implements fetching of statistic without arguments.
func (enum *Ascii_protocol_enum) stat(storage *cache.LRUCache, stats *stat.ServerStat) string {
	var result = ""
	if len(enum.key) == 0 {
		for key, value := range stats.Serialize(storage) {
			result += "STAT " + key + " " + value + "\r\n"
		}
	} else {
		switch enum.key[0] {
		case "settings":
			for key, value := range stats.Settings(storage) {
				result += "STAT " + key + " " + value + "\r\n"
			}
		case "items":
			/* ... */
		case "conns":
		default:
			return strings.Replace(CLIENT_ERROR_TEMP, "%s", "Command is not implemented.", 1)

		}
	}
	return result + "END\r\n"
}

//
func (enum *Ascii_protocol_enum) lru_crawler(storage *cache.LRUCache) string {
	switch enum.key[0]{
	case "enable":
		err := storage.EnableCrawler()
		if err != nil {
			return strings.Replace(CLIENT_ERROR_TEMP, "%s", err.Error(), 1)
		}
		return "OK\r\n"
	case "disable":
		storage.DisableCrawler()
		return "OK\r\n"
	case "tocrawl":
		if len(enum.key) < 2 {
			return strings.Replace(CLIENT_ERROR_TEMP, "%s", "Wrong parameters number.", 1)
		}
		amount, err := tools.StringToInt32(enum.key[1])
		if amount <= 0 || err != nil {
			return strings.Replace(CLIENT_ERROR_TEMP, "%s", "Invalid value of passed param.", 1)
		}
		storage.Crawler.ItemsPerRun = uint(amount)
		return "OK\r\n"
	case "sleep":
		if len(enum.key) < 2 {
			return strings.Replace(CLIENT_ERROR_TEMP, "%s", "Wrong parameters number.", 1)
		}
		amount, err := tools.StringToInt32(enum.key[1])
		if err != nil {
			return strings.Replace(CLIENT_ERROR_TEMP, "%s", "Invalid value of passed param.", 1)
		}
		err = storage.Crawler.SetSleep(amount)
		if err != nil {
			return strings.Replace(CLIENT_ERROR_TEMP, "%s", err.Error(), 1)
		}
		return "OK\r\n"
	default:
		return strings.Replace(CLIENT_ERROR_TEMP, "%s", "Command is not implemented.", 1)
	}
}

// Utilities

// Returns true if there was no "noreply" param in request.
func (enum *Ascii_protocol_enum) Reply() bool {
	return !enum.noreply
}

// Returns amount of bytes specified for data byte-string.
func (enum *Ascii_protocol_enum) DataLen() int {
	return enum.bytes
}

// Sets data byte-string of specified length to enumeration.
func (enum *Ascii_protocol_enum) SetData(data []byte) bool {
	if enum.bytes == len(data) {
		enum.data_string = data[0 : ]
		return true
	}
	return false
}

// Function checks was the passed param res successful whether not.
func IsMissed(res string) bool {
	return (res == NOT_FOUND || res == ERROR_TEMP || res == "END\r\n")
}

// Function increases fields of passed structure stats, if some of commands or passed param res were matched to required.
func (enum *Ascii_protocol_enum) RecordStats(stats *stat.ServerStat, res string) {
	if tools.In(enum.command, []string{"get", "set", "delete", "touch", }){
		stats.Commands["cmd_" + enum.command] ++
	}
	if tools.In(enum.command, []string{"get", "delete", "incr", "decr", "cas", "touch", }){
		if IsMissed(res){
			stats.Commands[enum.command + "_misses"] ++
		} else {
			stats.Commands[enum.command + "_hits"] ++
		}
	}
	if enum.command == "cas" && res == NOT_FOUND {
		stats.Commands["cas_badval"] ++
	}
}

// Function returns value of command field. //TODO: test
func (enum *Ascii_protocol_enum) Command() string {
	return enum.command
}
