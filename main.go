package main

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type teamNameScore struct {
	Name  string
	Score int
}

// teamNamePoints is used for the team league results and teamNameScore is used for the teams score per game
type teamNamePoints struct {
	Name   string
	Points int
}

func main() {
	teamPoints := receiveScannerInputs()
	sortedTeams := calculateTeamsOrder(teamPoints)
	printResults(sortedTeams, os.Stdout)
}

func receiveScannerInputs() map[string]int {
	scanner := bufio.NewScanner(os.Stdin)

	//initialize map that will be teamName: teamScore
	teamPoints := make(map[string]int)
	//use a counter so we can identify which line an error occurs when feeding in files for testing, could obviously also be used for production if its changed to use files
	counter := 0
	fmt.Println("Please enter your filename, the file lines must be in the following format 'teamname1 teamscore1, teamname2 teamscore2', you can end the process by entering 'end' or just pressing enter with no text")
	scanner.Scan()
	filename := scanner.Text()
	err1 := scanner.Err()
	if err1 != nil {
		log.Fatal(err1)
	}
	file, err2 := os.Open("./gameResults/" + filename)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(line) == 0 || line == "end" {
			break
		}
		counter++
		result := ""
		team1, team2, nameErr := getTeamNamesAndScores(line)
		if nameErr != nil {
			handleError(nameErr, counter, os.Stdout, os.Stdout)
		} else {
			//first determine which team wins or whether its a tie
			result = calculateTeamPoints(team1, team2)
			//update the team values in the teamPoints array based on the result from above
			initAndUpdateTeamPoints(teamPoints, team1.Name, team2.Name, result)
		}

	}
	return teamPoints
}

func initAndUpdateTeamPoints(teamPoints map[string]int, teamName1, teamName2, result string) {
	//make sure both teams exist before updating
	if _, exists := teamPoints[teamName1]; !exists {
		teamPoints[teamName1] = 0
	}
	if _, exists := teamPoints[teamName2]; !exists {
		teamPoints[teamName2] = 0
	}
	//if its a tie increase both team scores by 1, otherwise increase the winning team's points by 3
	if result == "tie" {
		teamPoints[teamName1] += 1
		teamPoints[teamName2] += 1
	} else {
		teamPoints[result] += 3
	}
}

func handleError(err error, lineNumber int, w, w2 io.Writer) {
	fmt.Fprintf(w, "Error: %s\n", err)
	fmt.Fprintf(w2, "Line %d has an error and was not added to the tally, please ensure your format for every line is team1 score, team2 score\n", lineNumber)
}

func calculateTeamsOrder(teamPoints map[string]int) []teamNamePoints {

	var sortedTeams []teamNamePoints

	for team, points := range teamPoints {
		sortedTeams = append(sortedTeams, teamNamePoints{team, points})
	}

	slices.SortFunc(sortedTeams, func(i, j teamNamePoints) int {
		if result := cmp.Compare(j.Points, i.Points); result != 0 {
			return result
		}
		return strings.Compare(i.Name, j.Name)
	})

	return sortedTeams
}

func printResults(sortedTeams []teamNamePoints, w io.Writer) {
	pos := 0
	previousPoints := 0
	//go through the teams and print the results
	for i, team := range sortedTeams {
		// if the previous team did not have the same amount of points, increase their position to their position in the sorted array + 1
		if previousPoints != team.Points || pos == 0 {
			pos = i + 1
		}
		//check whether its the last value in the array and don't add a \n, since the expected output did not have an empty line at the bottom
		if i != len(sortedTeams)-1 {
			//check whether the team does not have only one point, if that is the case write pt instead of pts
			if team.Points != 1 {
				fmt.Fprintf(w, "%d. %s, %d pts\n", pos, team.Name, team.Points)
			} else {
				fmt.Fprintf(w, "%d. %s, %d pt\n", pos, team.Name, team.Points)
			}
			previousPoints = team.Points
		} else {
			if team.Points != 1 {
				fmt.Fprintf(w, "%d. %s, %d pts", pos, team.Name, team.Points)
			} else {
				fmt.Fprintf(w, "%d. %s, %d pt", pos, team.Name, team.Points)
			}
		}
	}
}

func calculateTeamPoints(team1, team2 teamNameScore) (result string) {
	if team1.Score == team2.Score {
		return "tie"
	} else if team1.Score > team2.Score {
		return team1.Name
	} else {
		return team2.Name
	}
}

func getTeamNamesAndScores(inputLine string) (teamNameScore1, teamNameScore2 teamNameScore, calculationError error) {

	// get positions of commas that are preceded by a score
	splitPosRegex := regexp.MustCompile(`\s[0-9]+,\s`)
	var teamDelimPos = splitPosRegex.FindAllStringIndex(inputLine, -1)
	var teamDelimSplitArray = splitPosRegex.Split(inputLine, -1)

	if teamDelimPos == nil {
		return teamNameScore{"", 0}, teamNameScore{"", 0}, errors.New("invalid line, teams have to be separated by a comma , for example team1 3, team2 4")
	} else if len(teamDelimSplitArray) > 2 {
		return teamNameScore{"", 0}, teamNameScore{"", 0}, errors.New("invalid line, team names cannot contain commas and numbers in the following format ' 1234, '")
	}
	var teamNamesAndScores []string

	//copy the team segment and remove the comma
	team1 := strings.TrimSpace(inputLine[0:teamDelimPos[0][1]])
	team1 = team1[0 : len(team1)-1]
	team2 := strings.TrimSpace(inputLine[teamDelimPos[0][1]:])
	teamNamesAndScores = append(teamNamesAndScores, team1)
	teamNamesAndScores = append(teamNamesAndScores, team2)

	if len(teamNamesAndScores) < 2 {
		return teamNameScore{"", 0}, teamNameScore{"", 0}, errors.New("invalid line, no commas to separate teams")
	} else {

		getTeamNameAndScore := func(nameAndScoreString string) (team teamNameScore, scoreErr error) {
			nameScoreSplitRegex := regexp.MustCompile(`\s[0-9]+`)
			//find all the indexes of numbers preceded by a space
			scoreIndex := nameScoreSplitRegex.FindAllStringIndex(nameAndScoreString, -1)

			if scoreIndex == nil {
				return teamNameScore{"", 0}, errors.New("invalid line, could not find a team score")
			}
			//find last index of a space followed by a number this allows for teams with numbers in them
			teamName := nameAndScoreString[0:scoreIndex[len(scoreIndex)-1][0]]
			teamScore, scoreErr := strconv.Atoi(strings.TrimSpace(nameAndScoreString[scoreIndex[len(scoreIndex)-1][0]:scoreIndex[len(scoreIndex)-1][1]]))
			if scoreErr != nil {
				return teamNameScore{teamName, teamScore}, scoreErr
			}
			return teamNameScore{teamName, teamScore}, nil
		}
		var team1NameScore teamNameScore
		var team2NameScore teamNameScore
		var scoreErr error
		team1NameScore, scoreErr = getTeamNameAndScore(teamNamesAndScores[0])
		if scoreErr != nil {
			return teamNameScore1, teamNameScore2, scoreErr
		}
		team2NameScore, scoreErr = getTeamNameAndScore(teamNamesAndScores[1])
		if scoreErr != nil {
			return team1NameScore, team2NameScore, scoreErr
		}
		return team1NameScore, team2NameScore, nil
	}

}
