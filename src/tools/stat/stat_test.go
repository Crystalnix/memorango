package stat

import (
	"testing"
	"fmt"
)

func TestRusage(t *testing.T){
	var x ServerStat
	a, b := x.rusage(0)
	fmt.Println(a, b)
	a, b = x.rusage(1)
	fmt.Println(a, b)
	if a == 0 && b == 0 || b >= 1000000 {
		t.Fatalf("Unexpected returned values")
	}
}
