package tools

import "strconv"

//type for implementation of Cacheable interface
type StoredData struct {
	value []byte
	key string
}

// Following function implements interface Key() and returns key of the value
func (container *StoredData) Key() string {
	return container.key
}

// Following function implements interface Size() and returns amount of bytes
func (container *StoredData) Size() int {
	return len(container.value)
}

//interface for different types of protocols
type Protocol interface {
	// to think about implementation
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
