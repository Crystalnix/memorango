package main

import (
	"fmt"
	//"client"
	"server/tools"
	"server/tools/cache"
)

func main() {

	//port := "9999"

	//client.Client(port, []byte("get key 123 123 123 9292929292 noreply\r\n"))

	mc := cache.New(100)
	if !mc.Set(tools.NewStoredData([]byte("This is a value"), "KEY")) {
		fmt.Println("Cache couldn't store data.")
		return
	}else{
		x := mc.Get("KEY")
		y := x.Cacheable
		fmt.Println("Stored value='",string(tools.ExtractStoredData(y)),"'")
	}

	var input string
	fmt.Scanln(&input)
}
