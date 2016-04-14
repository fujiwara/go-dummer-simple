package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fujiwara/shapeio"
)

const (
	defaultMessage = "time:2013-11-20 23:39:42 +0900\tlevel:ERROR\tmethod:POST\turi:/api/v1/people\treqtime:3.1983877060667103"
)

var (
	messages []string
	bufSize  = 1024 * 1024
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
		messages = []string{message + "\n"}
	}
	f, err := os.Create(output)
	if err != nil {
		die(err)
	}
	defer f.Close()

	w := shapeio.NewWriter(f)
	if rate != 0 {
		avgMessageSize := 0
		for _, m := range messages {
			avgMessageSize += len(m)
		}
		avgMessageSize = avgMessageSize / len(messages)
		limit := float64(avgMessageSize) * rate
		w.SetRateLimit(limit)
	}
	bw := bufio.NewWriterSize(w, bufSize)

	timer := time.NewTimer(time.Duration(second) * time.Second)
	go func() {
		n := len(messages)
		for i := 0; ; i++ {
			io.WriteString(bw, messages[i%n])
		}
	}()
	<-timer.C
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
		messages = append(messages, scanner.Text()+"\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
