package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var fileName string
	var timeLimit int
	flag.StringVar(
		&fileName,
		"csv",
		"problems.csv",
		"a csv file in the format of 'question,answer'",
	)
	flag.IntVar(
		&timeLimit,
		"limit",
		30,
		"the time limit for the quiz in seconds",
	)
	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	questionChan := make(chan string)
	answerChan := make(chan string)
	timeUpChan := time.After(time.Duration(timeLimit) * time.Second)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		questionNumber := 1
		for {
			select {
			case q := <-questionChan:
				fmt.Printf("Question %d:\n", questionNumber)
				fmt.Println(q)
				a, _ := reader.ReadString('\n')
				answerChan <- a
				questionNumber += 1
			case <-timeUpChan:
				return
			}
		}
	}()

	correctCount := 0
	totalCount := len(lines)

	for _, line := range lines {
		question := line[0]
		answer := line[1]

		questionChan <- question

		select {
		case a := <-answerChan:
			if strings.Trim(a, "\n ") == answer {
				correctCount += 1
			}
		case <-timeUpChan:
			printResults(correctCount, totalCount)
			return
		}
	}
	printResults(correctCount, totalCount)
}

func printResults(correct, total int) {
    fmt.Printf("You got %d correct out of %d total questions.", correct, total)
}
