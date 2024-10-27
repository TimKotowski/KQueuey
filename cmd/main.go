package main

import (
	"flag"
	"fmt"

	"kqueuey"
)

func main() {
	flag.UintVar(&kqueuey.TestFlag, "test-flag", 0, "testing")
	flag.Parse()
	c, err := kqueuey.LoadConfiguration()
	if err != nil {
		panic(err)
	}
	fmt.Println(c)

}
