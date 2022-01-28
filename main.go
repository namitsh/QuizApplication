package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type inputData struct {
	filePath string
	timer    int
}

type Problem struct {
	q   string
	ans string
}

func readInputData() inputData {
	defaultFile := "problems.csv"
	fileName := flag.String("csv", defaultFile, "csv file in format of `question,ans`")
	timer := flag.Int("timer", 30, "timer file in seconds. default :30")

	flag.Parse()
	return inputData{
		filePath: *fileName,
		timer:    *timer,
	}
}

//read the csv file
func readCsvFile(fileName string) [][]string {
	if !strings.HasSuffix(fileName, ".csv") {
		panic("Invalid file type")
	}

	if _, err := os.Stat(fileName); err != nil && os.IsNotExist(err) {
		fmt.Printf("File Not Exist: %s\n", fileName)
		os.Exit(1)
	}
	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error occurred while opening file %s", fileName)
		os.Exit(1)
	}
	fmt.Println("Successfully opened a file", fileName)
	defer csvFile.Close()
	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println("Error occurred while reading csv file", err)
		os.Exit(1)
	}
	return lines
}

func startQuiz(problems []Problem, t int) int {
	timer := getTimer(t)
	result := 0
	for i, problem := range problems {
		fmt.Printf("Q #%d: %v = ", i+1, problem.q)
		answerCh := make(chan string)
		go func() {
			var a string
			fmt.Scanln(&a)
			answerCh <- a
		}()
		select {
		case <-timer.C:
			return result
		case answer := <-answerCh:
			if answer == problem.ans {
				result += 1
			}
		}
	}
	return result
}

func parseLines(lines [][]string) []Problem {
	problems := make([]Problem, len(lines))
	for i, val := range lines {
		problems[i] = Problem{q: val[0], ans: val[1]}
	}
	return problems
}

func getTimer(t int) *time.Timer {
	return time.NewTimer(time.Duration(t) * time.Second)
}

func main() {
	fmt.Println("Hello World!")
	inputData := readInputData()
	lines := readCsvFile(inputData.filePath)
	problems := parseLines(lines)
	result := startQuiz(problems, inputData.timer)
	fmt.Printf("Score: %d out of %d", result, len(problems))
}
