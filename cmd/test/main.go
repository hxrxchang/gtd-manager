package main

import (
	"flag"
	"fmt"
)

func main() {
	label := flag.String("label", "", "label name")
	flag.Parse()
	fmt.Println(*label)
}
