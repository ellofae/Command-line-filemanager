package wc

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func Test() {
	fmt.Println("Test")
}

var (
	linesCount = 0
	wordsCount = 0
	charsCount = 0
)

var signalChannel = make(chan bool, 0) // signal channel

func RegularWc(wrChannel chan string, filename string) error {
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
		linesCount++

		for _, word := range line {
			if word == ' ' {
				wordsCount++
			}
			charsCount++
		}
		wordsCount++
	}

	signalChannel <- true
}
