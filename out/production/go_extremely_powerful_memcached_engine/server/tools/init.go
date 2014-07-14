package tools

//type for implementation of Cacheable interface
type RawByteString struct {
	value []byte
	key string
}

// Following function implements interface Key() and returns key of the value
func (container *RawByteString) Key() string {
	return container.key
}

// Following function implements interface Size() and returns amount of bytes
func (container *RawByteString) Size() int {
	return len(container.value)
}
