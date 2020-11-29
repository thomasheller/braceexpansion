package main

import (
	"fmt"
	"os"
	be "github.com/thomasheller/braceexpansion"
)

func main() {
	tree, err := be.New().Parse(os.Args[1])
	if err != nil {
		panic(err)
	}
	for _, s := range tree.Expand() {
		fmt.Println(s)
	}
}
