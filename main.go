package main

import (
	"flag"
	"fmt"

	"github.com/pihao/himawari8-desktop/src"
)

const VERSION = "himawari8-desktop version 0.0.2"

func main() {
	v := flag.Bool("v", false, "show version.")
	flag.Parse()
	if *v {
		fmt.Println(VERSION)
	} else {
		src.Run()
	}

}
