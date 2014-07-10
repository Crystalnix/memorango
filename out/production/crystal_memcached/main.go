package main

import (
	"fmt"
	"os"
	"server/core"
	"client/core"
)

func main() {

	fmt.Printf("Servers: %d \n", len(os.Args [1 : ]))
	for _, elem := range os.Args [1 : ] {
		fmt.Printf(string(elem) + "\n")
	}


}
