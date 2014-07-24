/*
Package for utilities and some helpful functions for handling request of server and keeping data.
*/
package tools

import (
	"strconv"
	"reflect"
	"time"
	"crypto/sha1"
	"bytes"
	"bufio"
	"encoding/binary"
)

// Current version string.
const VERSION = "Go memcached implementation v1.0"

// The realization of Cacheable interface.
type StoredData struct {
	value []byte
	key string
}

// The public method for generalization of interface.
// Returns key of item.
func (container StoredData) Key() string {
	return container.key
}

// The public method for generalization of interface.
// Returns amount of value's bytes.
func (container StoredData) Size() int {
	return len(container.value)
}

// The public method of StoredData, which return value itself.
func (container StoredData) Value() []byte {
	return container.value
}

// Function creates instance of StoredData from received byte-string and key.
func NewStoredData(value []byte, key string) StoredData{
	return StoredData{value: value, key: key}
}

// Function matches entry of element in collection.
func In(element string, collection []string) bool{
	for _, value := range collection {
		if element == value { return true }
	}
	return false
}

// Function convert string element to 32-bit decimal integer.
// If it is impossible, there will be return error.
func StringToInt32(str string) (int, error) {
	value, err := strconv.ParseInt(str, 10, 32)
	return int(value), err
}

// Function convert string element to 32-bit decimal integer.
// If it is impossible, there will be return error.
func StringToInt64(str string) (int64, error) {
	value, err := strconv.ParseInt(str, 10, 64)
	return value, err
}

// Function convert 64-bit integer element to string.
// If it is impossible, there will be returned a error.
func IntToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// Function convert 64-bit unsigned integer element to string.
// If it is impossible, there will be returned a error.
func UIntToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// Function is supposed to convert data from (firstly) Cacheable interface or any other interface, which was generalized
// by StoredData type, back to StoredData and return byte-string value of it.
// If it is impossible there will be returned a nil.
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

// Function matches is the passed param Unix time stamp.
// The criteria of this match is set in memcached protocol:
// https://github.com/memcached/memcached/blob/master/doc/protocol.txt
// Thus parameter is timestamp, when it greater than 30 days.
func isUnixTimeStamp(ts int64) bool {
	return ts > 60 * 60 * 24 * 30 // if more than 30 days
}

// Function converts passed number into a timestamp offset from current moment,
// if such number not a zero and not a timestamp already.
// Otherwise, it just return passed value.
func ToTimeStampFromNow(ts int64) int64 {
	if !isUnixTimeStamp(ts) && ts != 0 {
		ts = time.Now().Add(time.Second * time.Duration(ts)).Unix()
	}
	return ts
}

// Function converts passed byte-string to unique uint64, thus creates Cas Unique
func GenerateCasId(buf []byte) int64 {
	var hashSum = sha1.Sum( append(buf, IntToString(time.Now().Unix())...) )
	var byteBuf = bytes.NewBuffer(hashSum[0 : ])
	reader := bufio.NewReader(byteBuf)
	num, err := binary.ReadVarint(reader)
	if err != nil {
		return 0
	}
	return num
}
