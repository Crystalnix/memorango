package protocol

import (
	"tools/cache"
	"tools"
	"strings"
	"errors"
	"fmt"
)

// Public method of Ascii_protocol_enum operates with received storage: retrieves, discards, sets or updates items,
// related to own containment.
// Returns response to client as byte-string and error/nil.
// If process was successful, there will be returned nil instead of error, otherwise it will be returned specified error.
func (enum *Ascii_protocol_enum) HandleRequest(storage *cache.LRUCache) ([]byte, error) {
	fmt.Println("Start handle request: ", enum)
	var err error
	if len(enum.error) > 0 {
		return []byte(enum.error), nil
	}
	var result string
	switch enum.command {
	case "set":
		result, err = enum.set(storage)
	case "add":
		result, err = enum.add(storage)
	case "append":
		result, err = enum.append(storage)
	case "prepend":
		result, err = enum.prepend(storage)
	case "get":
		result, err = enum.get(storage)
	case "touch":
		result, err = enum.touch(storage)
	case "delete":
		result, err = enum.delete(storage)
	case "flush_all":
		result, err = enum.flush_all(storage)
	case "version":
		return []byte("VERSION " + tools.VERSION + "\r\n"), nil
	case "quit":
		return nil, errors.New("It is not a error")
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

func (enum *Ascii_protocol_enum) add(storage *cache.LRUCache) (string, error) {
	if storage.Get(enum.key[0]) != nil {
		return NOT_STORED, nil
	}
	return enum.set(storage)
}

func (enum *Ascii_protocol_enum) pending(storage *cache.LRUCache, existed_item *cache.LRUCacheItem, pending_data []byte) (string, error) {
	enum.SetData(pending_data, len(pending_data))
	enum.exptime = existed_item.Exptime
	enum.cas_unique = existed_item.Cas_unique
	enum.flags = existed_item.Flags
	return enum.set(storage)
}

func (enum *Ascii_protocol_enum) prepend(storage *cache.LRUCache) (string, error) {
	existed_item := storage.Get(enum.key[0])
	if existed_item == nil {
		return NOT_STORED, nil
	}
	existed_data := tools.ExtractStoredData(existed_item.Cacheable)
	if existed_data == nil {
		return NOT_STORED, nil
	}
	return enum.pending(storage, existed_item, append(enum.data_string, existed_data))
}

func (enum *Ascii_protocol_enum) append(storage *cache.LRUCache) (string, error) {
	existed_item := storage.Get(enum.key[0])
	if existed_item == nil {
		return NOT_STORED, nil
	}
	existed_data := tools.ExtractStoredData(existed_item.Cacheable)
	if existed_data == nil {
		return NOT_STORED, nil
	}
	return enum.pending(storage, existed_item, append(existed_data, enum.data_string))
}

func (enum *Ascii_protocol_enum) replace(storage *cache.LRUCache) (string, error) {
	if storage.Get(enum.key[0]) == nil {
		return NOT_STORED, nil
	}
	return enum.set(storage)
}

func (enum *Ascii_protocol_enum) cas(storage *cache.LRUCache) (string, error) {
	return "", nil
}

// Retrieving commands

// Implements get method
func (enum *Ascii_protocol_enum) get(storage *cache.LRUCache) (string, error) {
	var result = ""
	for _, value := range enum.key{
		item := storage.Get(value)
		if item != nil {
			data := tools.ExtractStoredData(item.Cacheable)
			if data == nil {
				continue
			}
			result += "VALUE " + value + " " + tools.IntToString(int64(item.Flags)) + " " + tools.IntToString(int64(len(data)))
			if item.Cas_unique != 0 {
				result += " " + tools.IntToString(item.Cas_unique)
			}
			result += "\r\n"
			result += string(data) + "\r\n"
		}
	}
	return result + "END\r\n", nil
}

func (enum *Ascii_protocol_enum) gets(storage *cache.LRUCache) (string, error) {
	return "", nil
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
func (enum *Ascii_protocol_enum) SetData(data []byte, length int) bool {
	if enum.bytes == length {
		enum.data_string = data[0 : length]
		return true
	}
	return false
}
