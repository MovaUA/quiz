package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := flag.String("file", "problems.csv", "CSV file with problems to solve")
	flag.Parse()

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

	inReader := bufio.NewReader(os.Stdin)
	correct := 0

	for _, record := range records {
		fmt.Printf("%v = ", record[0])

		line, err := inReader.ReadString('\n')
		if err != nil {
			fmt.Printf("could not read answer: %v\n", err)
		}

		answer := strings.TrimSpace(line)
		if answer == record[1] {
			correct++
		}
	}

	fmt.Printf("Your score is %d of %d\n", correct, len(records))
}
