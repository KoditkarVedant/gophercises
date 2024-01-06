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

	args := parseArguments()
	problems := parseProblems(args.csvFileName)

	answerCh := make(chan result)
	finishCh := make(chan int)
	correct := 0

	go func() {
		startQuiz(problems, answerCh, finishCh)
	}()

	timer := time.NewTimer(time.Duration(args.time) * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Print("\n\nTime over!")
			fmt.Printf("\nYou scored %v out of %v", correct, len(problems))
			return
		case <-finishCh:
			fmt.Print("\nQuiz completed!")
			fmt.Printf("\nYou scored %v out of %v", correct, len(problems))
			return
		case result := <-answerCh:
			if result.res == result.ans {
				correct++
			}
		}
	}
}

func startQuiz(problems []problem, answerCh chan result, finishCh chan int) {
	for i, problem := range problems {

		fmt.Printf("Problem #%v: %v = ", i+1, problem.question)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err := scanner.Err()

		if err != nil {
			log.Fatal(err)
		}

		text := strings.TrimSpace(scanner.Text())

		answerCh <- result{
			res: text,
			ans: problem.answer,
		}
	}
	finishCh <- 0
}

type result struct {
	res string
	ans string
}

func parseArguments() arguments {
	args := arguments{}
	flag.StringVar(&args.csvFileName, "f", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.IntVar(&args.time, "time", 30, "time in seconds for an quiz")
	flag.Parse()
	return args
}

func parseProblems(csvFileName string) []problem {
	file, err := os.Open(csvFileName)

	if err != nil {
		log.Fatalf("Failed to open the CSV file: %s\n", csvFileName)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal("Failed to parse the provided CSV file.")
	}

	return parseLines(records)
}

func parseLines(lines [][]string) []problem {
	res := make([]problem, len(lines))
	for i, line := range lines {
		res[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return res
}

type arguments struct {
	csvFileName string
	time        int
}

type problem struct {
	question string
	answer   string
}
