package protocol

import (
	"server/tools/cache"
	"server/tools"
	"strings"
	"errors"
	"fmt"
)

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
	case "get":
		result, err = enum.get(storage)
	}
	return []byte(result), err
}

func (enum *Ascii_protocol_enum) set(storage *cache.LRUCache) (string, error){
	ind := storage.Set(tools.NewStoredData(enum.data_string, enum.key[0]))
	if ind {
		return "STORED\r\n", nil
	} else {
		return strings.Replace(SERVER_ERROR_TEMP, "%s", "Not enough memory", 1), errors.New("SERVER_ERROR")
	}
}

func (enum *Ascii_protocol_enum) add(storage *cache.LRUCache) string{
	return ""
}

func (enum *Ascii_protocol_enum) prepend(storage *cache.LRUCache) string{
	return ""
}

func (enum *Ascii_protocol_enum) append(storage *cache.LRUCache) string{
	return ""
}

func (enum *Ascii_protocol_enum) replace(storage *cache.LRUCache) string{
	return ""
}

func (enum *Ascii_protocol_enum) cas(storage *cache.LRUCache) string{
	return ""
}
// ######
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

func (enum *Ascii_protocol_enum) Reply() bool {
	return !enum.noreply
}
