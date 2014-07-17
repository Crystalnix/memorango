package tools

import (
	"strconv"
	"reflect"
)

//type for implementation of Cacheable interface
type StoredData struct {
	value []byte
	key string
}

// Following function implements interface Key() and returns key of the value
func NewStoredData(value []byte, key string) StoredData{
	return StoredData{value: value, key: key}
}

func (container StoredData) Key() string {
	return container.key
}

// Following function implements interface Size() and returns amount of bytes
func (container StoredData) Size() int {
	return len(container.value)
}

func (container StoredData) Value() []byte {
	return container.value
}

func In(element string, collection []string) bool{
	for _, value := range collection {
		if element == value { return true }
	}
	return false
}

func StringToInt32(str string) (int, error) {
	value, err := strconv.ParseInt(str, 10, 32)
	return int(value), err
}

func IntToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func ExtractStoredData(object interface {}) []byte {
	if reflect.TypeOf(object) == reflect.TypeOf(StoredData{}){
		val, ok := object.(StoredData)
		if ok {
			return val.Value()
		}
		return nil
	}
	return nil
}
