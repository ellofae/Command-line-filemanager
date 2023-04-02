package main

import (
	"context"
	"fmt"
	"main/cat"
	"main/wc"
	"os"
	"strings"
	"time"
)

var wrChannel = make(chan string, 5) // buffer channel

func main() {
	arguments := os.Args
	if len(arguments) < 3 {
		fmt.Println("You need minimum 3 arguments!")
		return
	}

	ctx := context.Background()
	ctx, closed := context.WithTimeout(ctx, time.Duration(4)*time.Second)
	defer closed()

	if arguments[1] == "wc" {
		wc.RegularWc(ctx, wrChannel, os.Args[2])
	} else if arguments[1] == "cat" {
		temp := strings.Split(arguments[2], ".")
		if len(temp) == 2 {
			cat.RegularCat(ctx, wrChannel, os.Args[2])
			return
		}

		cat.OptionCat(ctx, os.Args[2:])
	}
}
