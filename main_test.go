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

	//use error lines so we can have a few lines that we know will fail but not make the test fail
	type filesAndExpectedFails struct {
		fileName string
		errLines []int
	}

	var testCaseFileNames = []filesAndExpectedFails{
		{"resultnamesAndScores.txt", []int{}},
		{"resultnamesAndScores2.txt", []int{6, 7, 8}},
		{"resultnamesAndScores3.txt", []int{1}},
		{"resultnamesAndScores4.txt", []int{}},
	}

	var debugFile = "result3.txt"

	for _, value := range testCaseFileNames {
		fmt.Println("Testing File: ", value.fileName)
		f, err := os.Open("./gameResults/" + value.fileName)
		if err != nil {
			log.Fatal(err)
		}
		// close the file at the end of the program
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineCounter := 0

		for scanner.Scan() {
			//use the line counter so accurate feedback can be given on where something went wrong
			lineCounter++
			line := scanner.Text()
			team1, team2, nameErr := getTeamNamesAndScores(line)
			if nameErr != nil {
				//use contains function to check whether the line is in the slice of lines we know should fail
				if !contains(value.errLines, lineCounter) {
					t.Errorf("Line %d should not have an error", lineCounter)
				}
				fmt.Println("Invalid line: ", line)
			} else {
				t.Log("LINE PASSED: " + scanner.Text())
				//use debug file in case we would like to see the result if we think something should have a different result
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
	scannerInputTests := []struct {
		fileName       string
		expectedResult map[string]int
	}{
		{"result1.txt", map[string]int{"Lions": 5, "Snakes": 1, "Tarantulas": 6, "FC Awesome": 1, "Grouches": 0}},
		{"result2.txt", map[string]int{"Lions": 5, "Snakes": 1, "Tarantulas": 6, "FC Awesome": 1, "Grouches": 0}},
		{"result3.txt", map[string]int{"the, c4,ts": 3, "the, lions": 0}},
		{"result4.txt", map[string]int{"1970 Coca Cola": 9, "Pepsi": 0, "Fanta 1899": 6}},
	}

	for i, scannerInputTest := range scannerInputTests {

		reader, writer, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
			t.FailNow()
		}
		defer reader.Close()
		defer writer.Close()

		// Save the original os.Stdin for restoration later
		originalStdin := os.Stdin
		defer func() { os.Stdin = originalStdin }()

		os.Stdin = reader

		input := scannerInputTest.fileName
		go func() {
			fmt.Fprint(writer, input)
			writer.Close() // Close the writer to signal EOF
		}()

		scannerInputResult := receiveScannerInputs()

		if !reflect.DeepEqual(scannerInputResult, scannerInputTest.expectedResult) {
			t.Errorf("The results are not the same as expected for scannerInputTests %d\n expected: %v\n actual: %v", i+1, scannerInputTest.expectedResult, scannerInputResult)
		}
	}

}

func TestIntegration(t *testing.T) {

	scannerInputTests := []struct {
		fileName       string
		expectedResult map[string]int
	}{
		{"result1.txt", map[string]int{"Lions": 5, "Snakes": 1, "Tarantulas": 6, "FC Awesome": 1, "Grouches": 0}},
		{"result2.txt", map[string]int{"Lions": 5, "Snakes": 1, "Tarantulas": 6, "FC Awesome": 1, "Grouches": 0}},
		{"result3.txt", map[string]int{"the, c4,ts": 3, "the, lions": 0}},
		{"result4.txt", map[string]int{"1970 Coca Cola": 9, "Pepsi": 0, "Fanta 1899": 6}},
	}

	for i, scannerInputTest := range scannerInputTests {

		reader, writer, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
			t.FailNow()
		}
		defer reader.Close()
		defer writer.Close()

		// Save the original os.Stdin for restoration later
		originalStdin := os.Stdin
		defer func() { os.Stdin = originalStdin }()

		os.Stdin = reader

		input := scannerInputTest.fileName
		go func() {
			fmt.Fprint(writer, input)
			writer.Close() // Close the writer to signal EOF
		}()

		scannerInputResult := receiveScannerInputs()

		if !reflect.DeepEqual(scannerInputResult, scannerInputTest.expectedResult) {
			t.Errorf("The results are not the same as expected for scannerInputTests %d\n expected: %v\n actual: %v", i+1, scannerInputTest.expectedResult, scannerInputResult)
		}
		sortedTeams := calculateTeamsOrder(scannerInputResult)
		printResults(sortedTeams, os.Stdout)

	}

}
