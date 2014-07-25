package tools

import (
	"testing"
	"time"
)

func TestStoringData(t *testing.T){
	x := NewStoredData([]byte("111"), "1")
	if string(x.Value()) != "111" || x.Key() != "1" || x.Size() != 3 {
		t.Fatalf("The results of methods are unexpected.")
	}
}

func TestEntrance(t *testing.T){
	var x = "a"
	var y = "aa"
	var collection1 = []string{"a", "a", "a", "a"}
	if !In(x, collection1) || In(y, collection1) {
		t.Fatalf("The result is unexpected.")
	}

}

func TestStringToInt32Success(t *testing.T){
	res, err := StringToInt32("32000")
	if res != 32000 || err != nil {
		t.Fatalf("The result is unexpected.")
	}
	res, err = StringToInt32("-32000")
	if res != -32000 || err != nil {
		t.Fatalf("The result is unexpected.")
	}
	res, err = StringToInt32("-0")
	if res != 0 || err != nil {
		t.Fatalf("The result is unexpected.")
	}
}

func TestStringToInt32Fail(t *testing.T){
	_, err := StringToInt32("omfgbbq21321")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
	_, err = StringToInt32("--32")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
	_, err = StringToInt32("32.5")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
	_, err = StringToInt32("32e-5")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
}

func TestStringToInt64Success(t *testing.T){
	res, err := StringToInt64("64000")
	if res != int64(64000) || err != nil {
		t.Fatalf("The result is unexpected.")
	}
	res, err = StringToInt64("-64000")
	if res != int64(-64000) || err != nil {
		t.Fatalf("The result is unexpected.")
	}
	res, err = StringToInt64("-0")
	if res != int64(0) || err != nil {
		t.Fatalf("The result is unexpected.")
	}
}

func TestStringToInt64Fail(t *testing.T){
	_, err := StringToInt64("omfgbbq21641")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
	_, err = StringToInt64("--64")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
	_, err = StringToInt64("64.5")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
	_, err = StringToInt64("64e-5")
	if err == nil {
		t.Fatalf("The result is unexpected.")
	}
}

func TestIntToString(t *testing.T){
	res := IntToString(42)
	if res != "42" {
		t.Fatalf("The result is unexpected.")
	}
	res = IntToString(-42)
	if res != "-42" {
		t.Fatalf("The result is unexpected.")
	}
	res = IntToString(-0)
	if res != "0" {
		t.Fatalf("The result is unexpected.")
	}
}

func TestToTimeStamp(t *testing.T){
	// 2592000 == 3600 * 24 * 30 - is 30 days in seconds.
	if ToTimeStampFromNow(2592001) != 2592001 {
		t.Fatalf("The result is unexpected.")
	}
	if ToTimeStampFromNow(0) != 0 {
		t.Fatalf("The result is unexpected.")
	}
	now_ts := time.Now().Unix()
	if ToTimeStampFromNow(42) != now_ts + 42 {
		t.Fatalf("The result is unexpected.")
	}
	now_ts = time.Now().Unix()
	if ToTimeStampFromNow(-42) != now_ts - 42 {
		t.Fatalf("The result is unexpected.")
	}

}

func TestCasGenerator(t *testing.T){
	cas1 := GenerateCasId([]byte("Test"))
	cas2 := GenerateCasId([]byte("Test"))
	if cas1 == cas2 {
		t.Fatalf("Two identificators with same conditions are matched: %d", cas1)
	}
}

func TestUintToString(t *testing.T){
	if UIntToString(+424242) != "424242"{
		t.Fatalf("The result is unexpected.")
	}
}

func TestStringToUInt64(t *testing.T){
	res, err := StringToUInt64("424242")
	if err != nil || res != uint64(424242) {
		t.Fatalf("The result is unexpected: ", res, err)
	}
	res, err = StringToUInt64("-424242")
	if err == nil {
		t.Fatalf("The result is unexpected: ", res, err)
	}
}


func TestExtraction(t *testing.T){
	data := NewStoredData([]byte("42"), "42")
	if ExtractStoredData(data) == nil || ExtractStoredData("42") != nil {
		t.Fatalf("The result is unexpected.")
	}
}
