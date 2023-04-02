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
)

//var signalChannel = make(chan bool, 0) // signal channel
/*
var waitGroup sync.WaitGroup
var aMutex sync.Mutex
*/

var (
	filenames = make([]string, 0)
	options   = make([]string, 0)
)

const (
	TEMPFILE1 = "./temp1.txt"
	TEMPFILE2 = "./temp2.txt"
)

const PARAMS = 0644

func OptionCat(ctx context.Context, wrChannel chan string, args []string) error {
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
				numberNonEmptyLines(TEMPFILE1)
			case "-n":
				numberAllLines(TEMPFILE1)
			}
		}

		stdoutResult(file)
	}

	select {
	/*
		case <-signalChannel:
			fmt.Println("Signal caught!")
			return nil
	*/
	case <-ctx.Done():
		fmt.Println("Program execution ended not having being finished: time out!")
		return ctx.Err()
	}
}

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

func numberAllLines(filename string) {
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
	var counter byte = 0

	for {
		counter += 1

		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}
		io.WriteString(temp2, strconv.Itoa(int(counter))+" ")
		io.WriteString(temp2, line)
	}

	fileCopy(TEMPFILE2, TEMPFILE1)
}

func numberNonEmptyLines(filename string) {
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
	var counter byte = 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return
		}

		if line != "\n" {
			counter += 1
			io.WriteString(temp2, strconv.Itoa(int(counter))+" ")
			io.WriteString(temp2, line)
		} else {
			io.WriteString(temp2, "\n")
		}
	}

	fileCopy(TEMPFILE2, TEMPFILE1)
}

/*
func argumentReader(args []string) {
	for _, arg := range args {
		if ok := strings.Contains(arg, "-"); ok {
			options = append(options, arg)
		} else if ok := strings.Contains(arg, "."); ok {
			filename = append(filename, arg)
		} else {
			fmt.Println("Wrong argument given!\n[OPTION].. [FILENAME]..")
			os.Exit(1)
		}
	}
}
*/
