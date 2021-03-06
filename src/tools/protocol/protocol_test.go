package protocol

import (
	"testing"
	"reflect"
	"time"
	"fmt"
	"strings"
	"tools/cache"
	"tools/stat"
	"tools"
)

func matchEnumFields(enum *Ascii_protocol_enum,
					 com string,
					 keys []string,
					 flags int,
					 ts int64,
					 bytes int,
					 cas int64,
					 data []byte,
					 noreply bool,
					 error string ) bool {
	if error == enum.error && error != "" { return true }
	if com == enum.command && flags == enum.flags && ts == enum.exptime && bytes == enum.bytes &&
	   cas == enum.cas_unique && reflect.DeepEqual(keys, enum.key) &&
			reflect.DeepEqual(data, enum.data_string) && enum.noreply == noreply {
		return true
	}
	fmt.Println(enum)
	return false
}

func TestParsingSuiteSet1(t *testing.T){
	request := "set OMFG 1 42 15"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"set",
		keys,
		1,
		time.Now().Add(time.Second * time.Duration(42)).Unix(),
		15,
		0,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteSet2(t *testing.T){
	request := "set OMFG 1 42 15 noreply"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
					"set",
					keys,
					1,
					time.Now().Add(time.Second * time.Duration(42)).Unix(),
					15,
					0,
					nil,
					true,
					"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteCas1(t *testing.T){
	request := "cas OMFG 1 42 15 424242"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"cas",
		keys,
		1,
		time.Now().Add(time.Second * time.Duration(42)).Unix(),
		15,
		424242,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteCas2(t *testing.T){
	request := "cas OMFG 1 42 15 424242 noreply"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"cas",
		keys,
		1,
		time.Now().Add(time.Second * time.Duration(42)).Unix(),
		15,
		424242,
		nil,
		true,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteGet1(t *testing.T){
	request := "get OMFG"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"get",
		keys,
		0,
		0,
		0,
		0,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteGet2(t *testing.T){
	request := "get OMFG BBQ TEST TEST1 noreply"
	var keys = []string{"OMFG", "BBQ", "TEST", "TEST1", "noreply"}
	if !matchEnumFields(ParseProtocolHeader(request),
		"get",
		keys,
		0,
		0,
		0,
		0,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteAdd1(t *testing.T){
	request := "add OMFG 0 0 4"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"add",
		keys,
		0,
		0,
		4,
		0,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteAdd2(t *testing.T){
	request := "add OMFG 0 0 4 noreply"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"add",
		keys,
		0,
		0,
		4,
		0,
		nil,
		true,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteReplace1(t *testing.T){
	request := "replace OMFG 0 0 4"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"replace",
		keys,
		0,
		0,
		4,
		0,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteReplace2(t *testing.T){
	request := "replace OMFG 0 0 4 noreply"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"replace",
		keys,
		0,
		0,
		4,
		0,
		nil,
		true,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteAppend1(t *testing.T){
	request := "append OMFG 0 0 4"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"append",
		keys,
		0,
		0,
		4,
		0,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteAppend2(t *testing.T){
	request := "append OMFG 0 0 4 noreply"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"append",
		keys,
		0,
		0,
		4,
		0,
		nil,
		true,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuitePrepend1(t *testing.T){
	request := "prepend OMFG 0 0 4"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"prepend",
		keys,
		0,
		0,
		4,
		0,
		nil,
		false,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuitePrepend2(t *testing.T){
	request := "prepend OMFG 0 0 4 noreply"
	var keys = []string{"OMFG", }
	if !matchEnumFields(ParseProtocolHeader(request),
		"prepend",
		keys,
		0,
		0,
		4,
		0,
		nil,
		true,
		"") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingSuiteOther(t *testing.T){
	if !matchEnumFields(ParseProtocolHeader("delete OMFG noreply"),
		"delete", []string{"OMFG", }, 0, 0, 0, 0, nil, true, "") {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("touch OMFG 42"),
		"touch", []string{"OMFG", }, 0, time.Now().Add(time.Second * time.Duration(42)).Unix(), 0, 0, nil, false, "") {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("flush_all"),
		"flush_all", nil, 0, 0, 0, 0, nil, false, "") {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("version"),
		"version", nil, 0, 0, 0, 0, nil, false, "") {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("quit"),
		"quit", nil, 0, 0, 0, 0, nil, false, "") {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("stats"),
		"stats", []string{}, 0, 0, 0, 0, nil, false, "") {
		t.Fatalf("The parser works incorrect.")
	}
}

func TestParsingErrors(t *testing.T){
	if !matchEnumFields(ParseProtocolHeader("such command doesn't exist 1 2 3 4"),
		"", nil, 0, 0, 0, 0, nil, false, ERROR_TEMP) {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader(""),
		"", nil, 0, 0, 0, 0, nil, false, strings.Replace(CLIENT_ERROR_TEMP, "%s", "Input command line is empty", 1)) {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("set 1"),
		"", nil, 0, 0, 0, 0, nil, false, ERROR_TEMP) {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("cas a b c d e f"),
		"", nil, 0, 0, 0, 0, nil, false, ERROR_TEMP) {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("get"),
		"", nil, 0, 0, 0, 0, nil, false, ERROR_TEMP) {
		t.Fatalf("The parser works incorrect.")
	}
	if !matchEnumFields(ParseProtocolHeader("delete"),
		"", nil, 0, 0, 0, 0, nil, false, ERROR_TEMP) {
		t.Fatalf("The parser works incorrect.")
	}

}

func TestEnumReply1(t *testing.T){
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 42, 0, true, nil, ""}
	if testEnum.Reply(){
		t.Fatalf("Wrong behavior of Reply() function.")
	}
}

func TestEnumReply2(t *testing.T){
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 42, 0, false, nil, ""}
	if !testEnum.Reply(){
		t.Fatalf("Wrong behavior of Reply() function.")
	}
}

func TestEnumDataLen(t *testing.T){
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 42, 0, true, nil, ""}
	if testEnum.DataLen() != 42 {
		t.Fatalf("Wrong behavior of DataLen() function.")
	}
}

func TestEnumSetData1(t *testing.T){
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 42, 0, true, nil, ""}
	if !testEnum.SetData(make([]byte, 42)){
		t.Fatalf("Wrong behavior of SetData() function.")
	}
}

func TestEnumSetData2(t *testing.T){
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 42, 0, true, nil, ""}
	if testEnum.SetData(make([]byte, 41)){
		t.Fatalf("Wrong behavior of SetData() function.")
	}
}

func TestHandlingSuiteSet1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteSet2(t *testing.T){
	var storage = cache.New(4)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 42, 0, false, make([]byte, 42), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err == nil || string(res) != strings.Replace(SERVER_ERROR_TEMP, "%s", "Not enough memory", 1) {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteCas1(t *testing.T){
	var storage = cache.New(42)
	if !storage.Set(tools.NewStoredData([]byte("test1"), "key"), 0, 0, 424242) {
		t.Fatalf("Unexpecting behavior ")
	}
	var testEnum = Ascii_protocol_enum{"cas", []string{"key", }, 1, 0, 5, 424242, false, []byte("TEST2"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteCas2(t *testing.T){
	var storage = cache.New(42)
	if !storage.Set(tools.NewStoredData([]byte("test1"), "key"), 0, 0, 0) {
		t.Fatalf("Unexpecting behavior ")
	}
	var testEnum = Ascii_protocol_enum{"cas", []string{"key", }, 1, 0, 42, 424242, false, make([]byte, 42), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != NOT_FOUND {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteAdd1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"add", []string{"key1", }, 0, 0, 4, 0, false, []byte("TEST"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteAdd2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 22, 0, false, make([]byte, 22), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"add", []string{"key", }, 0, 0, 5, 0, false, []byte("TEST2"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "NOT_STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteReplace1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"replace", []string{"key1", }, 0, 0, 4, 0, false, []byte("TEST"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "NOT_STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteReplace2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 42, 0, false, make([]byte, 42), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"replace", []string{"key", }, 0, 0, 5, 0, false, []byte("TEST2"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteAppend1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"append", []string{"key1", }, 0, 0, 4, 0, false, []byte("TEST"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "NOT_STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteAppend2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"append", []string{"key", }, 0, 0, 5, 0, false, []byte("TEST2"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	stored := tools.ExtractStoredData(storage.Get("key").Cacheable)
	if stored == nil || string(stored) != "TESTTEST2" {
		t.Fatalf("Stored value is invalid: ", err, stored)
	}
}

func TestHandlingSuitePrepend1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"prepend", []string{"key1", }, 0, 0, 4, 0, false, []byte("TEST"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "NOT_STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuitePrepend2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	testEnum = Ascii_protocol_enum{"prepend", []string{"key", }, 0, 0, 5, 0, false, []byte("TEST2"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "STORED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	stored := tools.ExtractStoredData(storage.Get("key").Cacheable)
	if stored == nil || string(stored) != "TEST2TEST" {
		t.Fatalf("Stored value is invalid: ", err, stored)
	}
}

func TestHandlingSuiteGet1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"get", []string{"key", }, 0, 0, 0, 0, false, nil, ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "VALUE key 1 4\r\nTEST\r\nEND\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteGet2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"get", []string{"not_key", }, 0, 0, 0, 0, false, nil, ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "END\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteGets1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"gets", []string{"key", }, 0, 0, 0, 0, false, nil, ""}
	res, err = testEnum.HandleRequest(storage, nil)
	item := storage.Get("key")
	if item == nil {
		t.Fatalf("Item wasn't stored")
	}
	cas := item.Cas_unique
	if cas == 0 {
		t.Fatalf("Cas unique wasn't set.")
	}
	match_str := "VALUE key 1 4 " + tools.IntToString(cas) + "\r\nTEST\r\nEND\r\n"
	if err != nil || string(res) != match_str {
		t.Fatalf("Unexpected returned values of handling: ", err, res, testEnum)
	}
}

func TestHandlingSuiteGets2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"gets", []string{"not_key", }, 0, 0, 0, 0, false, nil, ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "END\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteGetMultiple(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key1", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"set", []string{"key2", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"get", []string{"key1", "key2", }, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "VALUE key1 1 4\r\nTEST\r\nVALUE key2 1 4\r\nTEST\r\nEND\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteIncrDecr1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 3, 0, false, []byte("123"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"incr", []string{"key", }, 0, 0, 0, 0, false, []byte("100"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "223\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, string(res))
	}
	testEnum = Ascii_protocol_enum{"decr", []string{"key", }, 0, 0, 0, 0, false, []byte("100"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "123\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, string(res))
	}
}

func TestHandlingSuiteIncrDecr2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 6, 0, false, []byte("3.1459"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"incr", []string{"key", }, 0, 0, 0, 0, false, []byte("100"), ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "ERROR\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, string(res))
	}
	testEnum = Ascii_protocol_enum{"decr", []string{"key1", }, 0, 0, 0, 0, false, []byte("100"), ""}
	res, err = testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "NOT_FOUND\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, string(res))
	}
}

func TestHandlingSuiteTouch1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"touch", []string{"key", }, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "TOUCHED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteTouch2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"touch", []string{"not_key", }, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != NOT_FOUND {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteDelete1(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"delete", []string{"key", }, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "DELETED\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteDelete2(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 4, 0, false, []byte("TEST"), ""}
	testEnum.HandleRequest(storage, nil)
	testEnum = Ascii_protocol_enum{"delete", []string{"not_key", }, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != NOT_FOUND {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteFlushAll(t *testing.T){
	var storage = cache.New(42)
	var testEnum = Ascii_protocol_enum{"flush_all", nil, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(storage, nil)
	if err != nil || string(res) != "OK\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteVersion(t *testing.T){
	var testEnum = Ascii_protocol_enum{"version", nil, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(nil, nil)
	if err != nil || string(res) != "VERSION "+ tools.VERSION +"\r\n" {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingSuiteQuit(t *testing.T){
	var testEnum = Ascii_protocol_enum{"quit", nil, 0, 0, 0, 0, false, nil, ""}
	res, err := testEnum.HandleRequest(nil, nil)
	if err == nil || res != nil {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingStatistic(t *testing.T){
	var testEnum = Ascii_protocol_enum{"stats", nil, 0, 0, 0, 0, false, nil, ""}
	var stats = stat.New(42, "9999", "8888", 1024, 2, true, true)
	var storage = cache.New(42)

	res, err := testEnum.HandleRequest(storage, stats)
	if err != nil || res == nil {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
}

func TestHandlingStatsRecording(t *testing.T){
	var testEnum = Ascii_protocol_enum{"set", []string{"key", }, 1, 0, 2, 0, true, []byte("42"), ""}
	stats := stat.New(42, "9999", "8888", 1024, 2, true, true)
	var storage = cache.New(42)
	res, err := testEnum.HandleRequest(storage, stats)
	if err != nil || res == nil {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	if stats.Commands == nil || stats.Commands["cmd_set"] != 1 {
		t.Fatalf("Wrong stats handling: ", stats.Commands)
	}

	stats = stat.New(42, "9999", "8888", 1024, 2, true, true)
	testEnum = Ascii_protocol_enum{"get", []string{"not_key", }, 0, 0, 0, 0, false, nil, ""}
	res, err = testEnum.HandleRequest(storage, stats)
	if err != nil || res == nil {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	if stats.Commands == nil || stats.Commands["cmd_get"] != 1 || stats.Commands["get_misses"] != 1 {
		t.Fatalf("Wrong stats handling: ", stats.Commands)
	}

	stats = stat.New(42, "9999", "8888", 1024, 2, true, true)
	testEnum = Ascii_protocol_enum{"get", []string{"key", }, 0, 0, 0, 0, false, nil, ""}
	res, err = testEnum.HandleRequest(storage, stats)
	if err != nil || res == nil {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	if stats.Commands == nil || stats.Commands["cmd_get"] != 1 || stats.Commands["get_hits"] != 1 {
		t.Fatalf("Wrong stats handling: ", stats.Commands)
	}

	stats = stat.New(42, "9999", "8888", 1024, 2, true, true)
	testEnum = Ascii_protocol_enum{"cas", []string{"key", }, 1, 0, 42, 424242, false, make([]byte, 42), ""}
	res, err = testEnum.HandleRequest(storage, stats)
	if err != nil || res == nil {
		t.Fatalf("Unexpected returned values of handling: ", err, res)
	}
	if stats.Commands == nil || stats.Commands["cas_badval"] != 1 {
		t.Fatalf("Wrong stats handling: ", stats.Commands)
	}
}

func TestHandlingMissing(t *testing.T){
	if !IsMissed(NOT_FOUND) || !IsMissed(ERROR_TEMP) || !IsMissed("END\r\n") || IsMissed("not_missing") {
		t.Fatalf("Invalid behavior of function.")
	}
}
