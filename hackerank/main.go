package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func logParser(inputFileName string, normalFileName string, errorFileName string) {
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Fail to open input file:", err)
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			fmt.Println("Fail to close input file:", err)
		}
	}(inputFile)

	normalFile, err := os.OpenFile(normalFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Fail to create normal log file:", err)
		return
	}
	defer func(normalFile *os.File) {
		err := normalFile.Close()
		if err != nil {
			fmt.Println("Fail to close normal log file:", err)
		}
	}(normalFile)

	errorFile, err := os.OpenFile(errorFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Fail to create error log file:", err)
		return
	}
	defer func(errorFile *os.File) {
		err := errorFile.Close()
		if err != nil {
			fmt.Println("Fail to close error log file:", err)
		}
	}(errorFile)

	normFileCh := make(chan string, 100)
	errorFileCh := make(chan string, 100)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		logWriter(normalFile, normFileCh)
	}()

	go func() {
		defer wg.Done()
		logWriter(errorFile, errorFileCh)
	}()

	scan := bufio.NewScanner(inputFile)
	for scan.Scan() {
		text := scan.Text()
		if strings.Contains(text, "ERROR:") {
			errorFileCh <- text + "\n"
		} else {
			normFileCh <- text + "\n"
		}
	}

	close(normFileCh)
	close(errorFileCh)
	wg.Wait()
}

//func logWriter(file io.Writer, ch <-chan string) {
//	for line := range ch {
//		_, err := file.Write([]byte(line))
//		if err != nil {
//			fmt.Println("Write error:", err)
//		}
//	}
//}

func logWriter(file io.StringWriter, ch <-chan string) {
	for line := range ch {
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}

func main() {

	inputFile := "log.txt"
	normalFileName := "normalLog.txt"
	errorFileName := "errorLog.txt"

	fmt.Println("🔍 Log ban dau:", inputFile)
	logParser(inputFile, normalFileName, errorFileName)
	fmt.Println("Sau khi duoc phan loai.")
	fmt.Println("➡ Log luu o normalFile:", normalFileName)
	fmt.Println("➡ Log luu tai eroor:", errorFileName)

}
