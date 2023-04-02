package wc

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

func RegularWc(ctx context.Context, wrChannel chan string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Didn't manage to open the file!")
		return err
	}

	reader := bufio.NewReader(file)
	go writeLines(wrChannel, reader)
	go processLines(wrChannel)

	select {
	case <-signalChannel:
		fmt.Println("Signal caught!")
		fmt.Println(linesCount, " ", wordsCount, " ", charsCount)
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

func processLines(rChannel <-chan string) {
	for line := range rChannel {
		waitGroup.Add(1)
		go func(str string) {
			aMutex.Lock()

			linesCount++

			if str != "\n" {
				for _, word := range str {
					if word == ' ' {
						wordsCount++
					}
					charsCount++
				}
				wordsCount++
			} else {
				charsCount++
			}

			aMutex.Unlock()
			waitGroup.Done()
		}(line)
	}

	waitGroup.Wait()
	signalChannel <- true
}
