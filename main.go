package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type OperationResult struct {
	operation string
	isCorrect bool
	answer    int
}

func main() {
	files := getFiles()

	for _, filename := range files {
		fmt.Printf("Processing the page: %s\n", filename)
		lines, err := readFileLines(filepath.Join("pages", filename))
		if err != nil {
			log.Fatalf("An error ocurred processing the file %s\n", filename)
		}

		processOperations(lines, filename)
	}
}

func processOperations(lines []string, filename string) {
	var isReady string
	for isReady != "yes" {
		fmt.Print("Are you ready to start? ")
		_, err := fmt.Scanln(&isReady)
		if err != nil {
			log.Fatalf("error reading the input: %s\n", err)
		}
	}

	var operationResults []OperationResult
	var solvedInLessTwoMinutes int

	starTime := time.Now()
	for _, operation := range lines {
		var userAnswer int
		fmt.Printf("What is the result for the next operation: %s\n", operation)
		_, err := fmt.Scanln(&userAnswer)
		if err != nil {
			log.Fatalf("error reading the input: %s\n", err)
		}

		isCorrect, answer := reviewOperationResult(operation, userAnswer)
		elapsed := time.Since(starTime)
		if elapsed < 2*time.Minute {
			solvedInLessTwoMinutes++
		}
		operationResults = append(operationResults, OperationResult{
			operation: operation,
			isCorrect: isCorrect,
			answer:    answer,
		})

	}

	var correctOperations int
	for _, result := range operationResults {
		if result.isCorrect {
			correctOperations++
		}
	}

	fmt.Printf("The user had %d correct operations\n", correctOperations)
	fmt.Printf("operations solved in less than two minutes %d\n", solvedInLessTwoMinutes)

	for _, result := range operationResults {
		if !result.isCorrect {
			fmt.Printf("The correct result for the operation: %s = %d\n", result.operation, result.answer)
		}
	}

	saveResultInCSVFile(filename, correctOperations, solvedInLessTwoMinutes)
}

func saveResultInCSVFile(filename string, correctOperations int, lessTwoMinutes int) {
	file, err := os.OpenFile("record.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open the file: %v", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	wrongOperations := strconv.Itoa(60 - correctOperations)
	currentTime := time.Now()
	date := currentTime.Format("02-01-2006")

	row := []string{filename, strconv.Itoa(correctOperations), wrongOperations, strconv.Itoa(lessTwoMinutes), date}
	if err := writer.Write(row); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatalf("Error writing csv: %v", err)
	}

	fmt.Println("Results saved in the CSV file")
}

func reviewOperationResult(operation string, userAnswer int) (bool, int) {
	parts := strings.Split(operation, " ")
	if len(parts) < 3 {
		return false, 0
	}

	firstNumber, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, 0
	}

	secondNumber, err := strconv.Atoi(parts[2])
	if err != nil {
		return false, 0
	}

	var correctAnswer int
	switch parts[1] {
	case "/":
		if secondNumber == 0 {
			return false, 0
		}

		correctAnswer = firstNumber / secondNumber
	case "+":
		correctAnswer = firstNumber + secondNumber
	case "x":
		correctAnswer = firstNumber * secondNumber
	case "-":
		correctAnswer = firstNumber - secondNumber
	}

	return userAnswer == correctAnswer, correctAnswer

}

func getFiles() []string {
	var fileNames []string
	dirPath := "pages"

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filename := file.Name()
		fileNames = append(fileNames, filename)
	}

	return fileNames
}

func readFileLines(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	return lines, nil
}
