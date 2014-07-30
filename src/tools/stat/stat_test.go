package stat

import (
	"testing"
	"fmt"
)

func TestRusage(t *testing.T){
	var x ServerStat
	a, b, c, d := x.rusage()
	fmt.Println(a, b)
	fmt.Println(c, d)
	if a == 0 && b == 0 || b >= 1000000 || c == 0 && d == 0 || d >= 1000000 {
		t.Fatalf("Unexpected returned values")
	}
}
