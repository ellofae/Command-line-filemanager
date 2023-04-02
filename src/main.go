package main

import (
	"fmt"
	"main/wc"
	"os"
)

var wrChannel = make(chan string, 5)

func main() {
	arguments := os.Args
	if len(arguments) < 3 {
		fmt.Println("You need minimum 3 arguments!")
		return
	}

	if arguments[1] == "wc" {
		wc.RegularWc(wrChannel, os.Args[2])
	}
}
