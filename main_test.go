package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

func contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func TestGetTeamNamesAndScores(t *testing.T) {

	type filesAndExpectedFails struct {
		fileName string
		errLines []int
	}

	var testCaseFileNames = []filesAndExpectedFails{
		{"result1.txt", []int{}},
		{"result2.txt", []int{6, 7, 8}},
		{"result3.txt", []int{1}},
		{"result4.txt", []int{}},
	}

	var debugFile = "result3.txt"

	for _, value := range testCaseFileNames {
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
		lineCounter := 0

		for scanner.Scan() {
			lineCounter++
			// do something with a line
			line := scanner.Text()
			team1, team2, nameErr := getTeamNamesAndScores(line)
			if nameErr != nil {
				if !contains(value.errLines, lineCounter) {
					t.Errorf("Line %d should not have an error", lineCounter)
				}
				fmt.Println("Invalid line: ", line)
			} else {
				t.Log("LINE PASSED: " + scanner.Text())
				if value.fileName == debugFile {
					fmt.Println("Line: ", line)
					fmt.Printf("%s - %d \n", team1.Name, team1.Score)
					fmt.Printf("%s - %d \n", team2.Name, team2.Score)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

}

func TestInitAndUpdateTeamPoints(t *testing.T) {
	teams := []struct {
		teamPoints                map[string]int
		team1                     string
		team2                     string
		result                    string
		teamPointsAfterProcessing map[string]int
	}{
		{map[string]int{"Anton": 1}, "Anton", "SPAN", "tie", map[string]int{"Anton": 2, "SPAN": 1}},
		{map[string]int{"Anton": 1, "SPAN": 1}, "Anton", "SPAN", "tie", map[string]int{"Anton": 2, "SPAN": 2}},
		{map[string]int{}, "Anton", "SPAN", "tie", map[string]int{"Anton": 1, "SPAN": 1}},
		{map[string]int{"Anton": 1}, "Anton", "SPAN", "Anton", map[string]int{"Anton": 4, "SPAN": 0}},
		{map[string]int{"Anton": 1}, "Anton", "SPAN", "SPAN", map[string]int{"Anton": 1, "SPAN": 3}},
	}

	for i, team := range teams {
		initAndUpdateTeamPoints(team.teamPoints, team.team1, team.team2, team.result)
		if !reflect.DeepEqual(team.teamPoints, team.teamPointsAfterProcessing) {
			t.Errorf("The team points result is not the same as expected for result %d\n expected result: %v\n actual result: %v", i+1, teams[i].teamPointsAfterProcessing, teams[i].teamPoints)
		}
	}
}

func TestCalculateTeamsOrder(t *testing.T) {

	teamPointsTestCases := []struct {
		teamPoints     map[string]int
		expectedResult []teamNamePoints
	}{
		{map[string]int{"Ben": 1, "Anton": 1, "Span": 3}, []teamNamePoints{{"Span", 3}, {"Anton", 1}, {"Ben", 1}}},
		{map[string]int{"Ben": 3, "Anton": 3, "Span": 3}, []teamNamePoints{{"Anton", 3}, {"Ben", 3}, {"Span", 3}}},
		{map[string]int{"Ben": 1, "Anton": 2, "Span": 3}, []teamNamePoints{{"Span", 3}, {"Anton", 2}, {"Ben", 1}}},
	}

	for i, teamPointTestCase := range teamPointsTestCases {
		result := calculateTeamsOrder(teamPointTestCase.teamPoints)
		for j, resultVal := range result {
			if teamPointTestCase.expectedResult[j] != resultVal {
				t.Errorf("The order of the expected result for testcase %d is not the same as the order of the return result from the calculateTeamsOrder function\n expected result: %v\n actual result: %v", i+1, teamPointTestCase.expectedResult, result)
			}
		}
	}
}

func TestHandleError(t *testing.T) {
	testErrors := []struct {
		errMsg          error
		lineNumber      int
		expectedResult1 string
		expectedResult2 string
	}{
		{errors.New("an error has occurred"), 3, "Error: an error has occurred\n", "Line 3 has an error and was not added to the tally, please ensure your format for every line is team1 score, team2 score\n"},
		{errors.New("could not find team name"), 6, "Error: could not find team name\n", "Line 6 has an error and was not added to the tally, please ensure your format for every line is team1 score, team2 score\n"},
	}

	for i, testError := range testErrors {
		var output bytes.Buffer
		var output2 bytes.Buffer
		handleError(testError.errMsg, testError.lineNumber, &output, &output2)
		if output.String() != testError.expectedResult1 {
			t.Errorf("The errors are not the same as expected for testError %d\n expected: %s\n actual: %s", i+1, testError.expectedResult1, output.String())
		}
		if output2.String() != testError.expectedResult2 {
			t.Errorf("The errors are not the same as expected for testError %d\n expected: %s\n actual: %s", i+1, testError.expectedResult2, output2.String())
		}
	}
}

func TestCalculateTeamPoints(t *testing.T) {
	testTeamPoints := []struct {
		team1          teamNameScore
		team2          teamNameScore
		expectedResult string
	}{
		{teamNameScore{"Anton", 1}, teamNameScore{"Span", 2}, "Span"},
		{teamNameScore{"Anton", 2}, teamNameScore{"Span", 2}, "tie"},
		{teamNameScore{"Anton", 2}, teamNameScore{"Span", 1}, "Anton"},
	}

	for i, testTeamPoint := range testTeamPoints {

		result := calculateTeamPoints(testTeamPoint.team1, testTeamPoint.team2)
		if result != testTeamPoint.expectedResult {
			t.Errorf("The result is not the same as expected for testTeamPoints %d\n expected: %s\n actual: %s", i+1, testTeamPoint.expectedResult, result)
		}
	}
}

func TestPrintResults(t *testing.T) {
	testSortedTeams := []struct {
		teams          []teamNamePoints
		expectedResult string
	}{
		{[]teamNamePoints{{"Span", 3}, {"Anton", 1}, {"Ben", 1}}, "1. Span, 3 pts\n2. Anton, 1 pt\n2. Ben, 1 pt"},
		{[]teamNamePoints{{"Anton", 3}, {"Ben", 3}, {"Span", 3}}, "1. Anton, 3 pts\n1. Ben, 3 pts\n1. Span, 3 pts"},
		{[]teamNamePoints{{"Span", 3}, {"Anton", 2}, {"Ben", 1}}, "1. Span, 3 pts\n2. Anton, 2 pts\n3. Ben, 1 pt"},
	}
	for i, testTeamPoint := range testSortedTeams {
		var output bytes.Buffer
		printResults(testTeamPoint.teams, &output)
		if output.String() != testTeamPoint.expectedResult {
			t.Errorf("The result is not the same as expected for testSortedTeams %d\n expected: %s\n actual: %s", i+1, testTeamPoint.expectedResult, output.String())
		}
	}

}

func TestReceiveScannerInputs(t *testing.T) {

}
