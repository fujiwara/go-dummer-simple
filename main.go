package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fujiwara/shapeio"
)

const (
	defaultMessage = "time:2013-11-20 23:39:42 +0900\tlevel:ERROR\tmethod:POST\turi:/api/v1/people\treqtime:3.1983877060667103"
	timeResolition = 20
)

var (
	messages       [][]byte
	defaultBufSize = 1024 * 1024
	LF             = []byte{10}
)

func main() {
	var (
		second  int64
		output  string
		input   string
		message string
		rate    float64
	)

	flag.Int64Var(&second, "s", 1, "Duration of running in second")
	flag.StringVar(&output, "o", "dummy.log", "Output file")
	flag.StringVar(&input, "i", "", "Input file (Output messages by reading lines of the file in rotation)")
	flag.StringVar(&message, "m", defaultMessage, "Output message")
	flag.Float64Var(&rate, "r", 0, "Number of generating messages per second")
	flag.Parse()

	if input != "" {
		err := loadMessages(input)
		if err != nil {
			die(err)
		}
	} else {
		m := []byte(message + "\n")
		messages = [][]byte{m}
	}
	f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		die(err)
	}
	defer f.Close()

	var bufSize int
	w := shapeio.NewWriter(f)
	if rate != 0 {
		avgMessageSize := 0
		for _, m := range messages {
			avgMessageSize += len(m)
		}
		avgMessageSize = avgMessageSize / len(messages)
		limit := float64(avgMessageSize) * rate
		w.SetRateLimit(limit)
		if limit > timeResolition {
			bufSize = int(limit / timeResolition)
		} else {
			bufSize = int(limit)
		}
	} else {
		bufSize = defaultBufSize
	}
	bw := bufio.NewWriterSize(w, bufSize)

	running := true
	done := make(chan interface{})
	timer := time.NewTimer(time.Duration(second) * time.Second)
	go func() {
		n := len(messages)
		for i := 0; running; i++ {
			bw.Write(messages[i%n])
		}
		bw.Flush()
		done <- true
	}()
	<-timer.C
	running = false
	<-done
}

func die(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func loadMessages(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		line = append(line, LF...)
		messages = append(messages, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
