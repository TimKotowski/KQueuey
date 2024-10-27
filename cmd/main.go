package main

import (
	"flag"

	"kqueuey"
)

func main() {
	flag.UintVar(&kqueuey.TestFlag, "test-flag", 0, "testing")
	flag.Parse()
}
