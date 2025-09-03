package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	const filePath = "messages.txt"
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error: could not open file %s: %s", filePath, err.Error())
	}
	outputs := getChannelLines(f)

	for s := range outputs {
		fmt.Printf("read: %s\n", s)
	}
}

func getChannelLines(f io.ReadCloser) <-chan string {
	buff_size := 8
	buff := make([]byte, buff_size)
	outputs := make(chan string)

	go func() {
		defer f.Close()
		defer close(outputs)
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
				fmt.Printf("error: %s\n", err.Error())
			}

			str_buff := string(buff[:n])
			parts := strings.Split(str_buff, "\n")
			for i := 0; i < len(parts)-1; i++ {
				curr_line = curr_line + parts[i]
				outputs <- curr_line
				curr_line = ""
			}
			curr_line += parts[len(parts)-1]
		}
	}()
	return outputs
}
