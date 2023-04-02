package cat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Signal channel
var sigChannel = make(chan bool, 0)

var wg sync.WaitGroup
var aM sync.Mutex

// Channels for service goroutine
var (
	wChannel = make(chan string)
	rChannel = make(chan string)
)

// Slices for containing cat's options and filenames
var (
	filenames = make([]string, 0)
	options   = make([]string, 0)
)

// Temporary files
const (
	TEMPFILE1 = "./temp1.txt"
	TEMPFILE2 = "./temp2.txt"
)

func OptionCat(ctx context.Context, args []string) error {
	for _, arg := range args {
		if ok := strings.Contains(arg, "-"); ok {
			options = append(options, arg)
		} else if ok := strings.Contains(arg, "."); ok {
			filenames = append(filenames, arg)
		} else {
			fmt.Println("Wrong argument given!\n[OPTION].. [FILENAME]..")
			os.Exit(1)
		}
	}

	for _, file := range filenames {
		tempFilesCreate(TEMPFILE1, TEMPFILE2, file)

		for _, option := range options {
			switch option {
			case "-b":
				numberNonEmptyLines(ctx, TEMPFILE1)
			case "-n":
				numberAllLines(ctx, TEMPFILE1)
			}
		}

		stdoutResult(TEMPFILE1)
	}

	return nil
}

// Option functionality
func numberAllLines(ctx context.Context, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Didn't manage to open the file %s\n", filename)
		return
	}
	defer file.Close()

	temp2, err := os.OpenFile(TEMPFILE2, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Didn't manage to open the file %s\n", filename)
		return
	}
	defer temp2.Close()

	reader := bufio.NewReader(file)

	go monitor()
	go allLinesProccess(reader, temp2)

	select {
	case <-sigChannel:
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		return
	}

	wg.Wait()
	fileCopy(TEMPFILE2, TEMPFILE1)
}

func numberNonEmptyLines(ctx context.Context, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Didn't manage to open the file %s\n", filename)
		return
	}
	defer file.Close()

	temp2, err := os.OpenFile(TEMPFILE2, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Didn't manage to open the file %s\n", filename)
		return
	}
	defer temp2.Close()

	reader := bufio.NewReader(file)

	go monitor()
	go notEmptyProccess(reader, temp2)

	select {
	case <-sigChannel:
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		return
	}

	wg.Wait()
	fileCopy(TEMPFILE2, TEMPFILE1)
}

// Help functions that use service goroutine
func notEmptyProccess(reader *bufio.Reader, temp2 *os.File) {
	var counter byte = 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		wChannel <- line

		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			aM.Lock() // using because of the counter variable

			if line != "\n" {
				counter += 1
				io.WriteString(temp2, strconv.Itoa(int(counter))+" ")
				io.WriteString(temp2, line)
			} else {
				io.WriteString(temp2, "\n")
			}

			aM.Unlock()

		}(<-rChannel)
	}

	wg.Wait()
	sigChannel <- true
}

func allLinesProccess(reader *bufio.Reader, temp2 *os.File) {
	var counter byte = 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		wChannel <- line

		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			aM.Lock() // using because of the counter variable

			counter++
			io.WriteString(temp2, strconv.Itoa(int(counter))+" ")
			io.WriteString(temp2, line)

			aM.Unlock()

		}(<-rChannel)
	}

	wg.Wait()
	sigChannel <- true
}

// Service goroutine
func monitor() {
	var line string
	for {
		select {
		case temp := <-wChannel:
			line = temp
		case rChannel <- line:
		}
	}
}

// File managagement functionality
func stdoutResult(filename string) {
	b, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))

	err = os.Remove(TEMPFILE1)
	if err != nil {
		log.Println(err)
	}

	err = os.Remove(TEMPFILE2)
	if err != nil {
		log.Println(err)
	}
}

func fileCopy(src string, dest string) {
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Printf("Didn't manage to read the file '%s'\n", src)
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dest, bytesRead, 0644)
	if err != nil {
		fmt.Printf("Didn't manage to write to the file '%s'\n", dest)
		log.Fatal(err)
	}
}

func tempFilesCreate(filename1 string, filename2 string, src string) {
	_, err := os.Create(TEMPFILE1)
	if err != nil {
		fmt.Printf("Didn't manage to create a temp file: '%s'\n", TEMPFILE1)
		os.Exit(1)
	}

	_, err = os.Create(TEMPFILE2)
	if err != nil {
		fmt.Printf("Didn't manage to create a temp file: '%s'\n", TEMPFILE1)
		os.Exit(1)
	}

	fileCopy(src, filename1)
}
