package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	filename := flag.String("file", "problems.csv", "CSV file with problems to solve")
	seconds := flag.Uint("time", 30, "timeout (seconds) for solving all problems")

	flag.Parse()

	fmt.Println("press Enter to start...")

	r := bufio.NewReader(os.Stdin)

	if _, err := r.ReadString('\n'); err != nil {
		fmt.Printf("\ncould not read stdin: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Solve the problems:")

	file, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("could not open file: %v\n", err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = 2

	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("could not read csv file: %v\n", err)
		os.Exit(1)
	}

	correct := 0

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(*seconds)*time.Second,
	)
	defer cancel()

	answerCh := make(chan a)
	go readAnswer(r, answerCh)

loop:
	for _, record := range records {
		fmt.Printf("%v = ", record[0])

		select {
		case a := <-answerCh:
			if a.err != nil {
				fmt.Printf("could not read answer: %v\n", err)
				os.Exit(1)
			}
			if strings.EqualFold(a.answer, record[1]) {
				correct++
			}
		case <-ctx.Done():
			fmt.Printf("\ntimeout of %d seconds is expired\n", *seconds)
			break loop
		}
	}

	fmt.Printf("Your score is %d of %d\n", correct, len(records))
}

func readAnswer(r *bufio.Reader, answerCh chan<- a) {
	line, err := r.ReadString('\n')
	if err != nil {
		answerCh <- a{err: err}
		return
	}
	answerCh <- a{answer: strings.TrimSpace(line)}
}

type a struct {
	answer string
	err    error
}
