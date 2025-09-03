package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	const inputFilePath = "messages.txt"
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("error: failed to open file %s: %s\n", inputFilePath, err)
	}
	defer f.Close()

	output_channel := getLinesChannel(f)
	for s := range output_channel {
		fmt.Println("read:", s)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	outputs := make(chan string)
	buff_len := 8
	go func() {
		defer f.Close()
		defer close(outputs)
		buff := make([]byte, buff_len)
		curr_line := ""
		for {
			n, err := f.Read(buff)
			if err != nil {
				if curr_line != "" {
					outputs <- curr_line
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println("error: %s\n", err.Error())
				return

			}
			str_buff := string(buff[:n])
			parts := strings.Split(str_buff, "\n")
			for i := 0; i < len(parts)-1; i++ {
				outputs <- fmt.Sprintf("%s%s", curr_line, parts[i])
				curr_line = ""
			}
			curr_line += parts[len(parts)-1]
		}
	}()
	return outputs
}
