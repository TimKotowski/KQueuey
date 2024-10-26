package main

import (
	"fmt"

	"kqueuey"
)

func main() {
	c, err := kqueuey.LoadConfiguration()
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
	for _, node := range c.RaftOpts.Nodes {
		fmt.Println(node[])
	}
}
