package tools

import "testing"

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
