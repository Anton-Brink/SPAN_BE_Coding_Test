package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

type filesAndExpectedFails struct {
	fileName       string
	expectedFailed int
}

var testCaseFileNames = []filesAndExpectedFails{
	{"result1.txt", 0},
	{"result2.txt", 3},
	// {"result3.txt", 1},
	{"result4.txt", 0},
}

var debugFile = "result4.txt"

func TestGetTeamNamesAndScores(t *testing.T) {
	errCounter := 0
	for _, value := range testCaseFileNames {
		errCounter = 0
		// open file
		fmt.Println("Testing File: ", value.fileName)
		f, err := os.Open("./gameResults/" + value.fileName)
		if err != nil {
			log.Fatal(err)
		}
		// remember to close the file at the end of the program
		defer f.Close()

		// read the file line by line using scanner
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			// do something with a line
			line := scanner.Text()
			teamName1, teamName2, teamScore1, teamScore2, nameErr := getTeamNamesAndScores(line)
			if nameErr != nil {
				errCounter++
				fmt.Println("Invalid line: ", line)
			} else {
				t.Log("LINE PASSED: " + scanner.Text())
				if value.fileName == debugFile {
					fmt.Println("Line: ", line)
					fmt.Printf("%s - %d \n", teamName1, teamScore1)
					fmt.Printf("%s - %d \n", teamName2, teamScore2)
				}
			}
		}
		if errCounter != value.expectedFailed {
			t.Errorf("There are %d errors, but there should be %d errors", errCounter, value.expectedFailed)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

}
