package cat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

var signalChannel = make(chan bool, 0) // signal channel

var waitGroup sync.WaitGroup
var aMutex sync.Mutex

var (
	linesCount = 0
	wordsCount = 0
	charsCount = 0
)

func RegularCat(ctx context.Context, wrChannel chan string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Didn't manage to open the file!")
		return err
	}

	reader := bufio.NewReader(file)
	go writeLines(wrChannel, reader)
	go readLines(wrChannel)

	select {
	case <-signalChannel:
		fmt.Println("Signal caught!")
		return nil
	case <-ctx.Done():
		fmt.Println("Program execution ended not having being finished: time out!")
		return ctx.Err()
	}
}

func writeLines(wChannel chan<- string, reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("An error while reading the file occured!")
			return
		}

		wChannel <- line
	}

	close(wChannel)
}

func readLines(rChannel <-chan string) {
	for line := range rChannel {
		fmt.Print(line)
	}

	signalChannel <- true
}
