package main

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	// fmt.Println("input team results:")
	scanner := bufio.NewScanner(os.Stdin)

	teamPoints := make(map[string]int)
	counter := 0

	for {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 || line == "end" {
			break
		}
		counter++
		teamName1, teamName2, team1Score, team2Score, nameErr := getTeamNamesAndScores(line)
		result, calcErr := calculateTeamPoints(team1Score, team2Score, teamName1, teamName2)

		if calcErr != nil || nameErr != nil {
			//inform user of error and how to fix it
			fmt.Println("CalcErr", calcErr)
			fmt.Println("NameErr: ", nameErr)
			fmt.Println("Line ", counter, " has an error and was not added to the tally, please ensure your format for every line is team1 score, team2 score")
		} else {
			if _, exists := teamPoints[teamName1]; !exists {
				teamPoints[teamName1] = 0
			}
			if _, exists := teamPoints[teamName2]; !exists {
				teamPoints[teamName2] = 0
			}
			// println("Result: ", result)
			if result == "tie" {
				teamPoints[teamName1] += 1
				teamPoints[teamName2] += 1
			} else if result == teamName1 {
				teamPoints[teamName1] += 3
			} else {
				teamPoints[teamName2] += 3
			}
		}

	}

	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("output: ")
	type teamNameScore struct {
		Name  string
		Score int
	}

	var sortedTeams []teamNameScore

	for team, score := range teamPoints {
		sortedTeams = append(sortedTeams, teamNameScore{team, score})
	}

	slices.SortFunc(sortedTeams, func(i, j teamNameScore) int {
		return cmp.Compare(j.Score, i.Score)
	})

	pos := 0
	previousScore := 0
	for i, team := range sortedTeams {
		if previousScore != team.Score || pos == 0 {
			pos = i + 1
		}
		if i != len(sortedTeams)-1 {
			fmt.Printf("%d. %s, %d pts\n", pos, team.Name, team.Score)
			previousScore = team.Score
		} else {
			fmt.Printf("%d. %s, %d pts", pos, team.Name, team.Score)
		}
	}
}

func calculateTeamPoints(team1Score, team2Score int, team1Name, team2Name string) (result string, calculationError error) {
	if team1Score == team2Score {
		return "tie", nil
	} else if team1Score > team2Score {
		return team1Name, nil
	} else {
		return team2Name, nil
	}
}

func getTeamNamesAndScores(inputLine string) (teamName1, teamName2 string, teamScore1, teamScore2 int, calculationError error) {

	// get positions of commas that are preceded by a score
	splitPosRegex := regexp.MustCompile(`\s[0-9]+,\s`)
	var teamDelimPos = splitPosRegex.FindAllStringIndex(inputLine, -1)
	var teamDelimSplitArray = splitPosRegex.Split(inputLine, -1)

	if teamDelimPos == nil {
		return "", "", 0, 0, errors.New("invalid line, teams have to be separated by a comma , for example team1 3, team2, 4")
	} else if len(teamDelimSplitArray) > 2 {
		return "", "", 0, 0, errors.New("invalid line, team names cannot contain commas and numbers in the following format ' 1234, '")
	}
	var teamNamesAndScores []string

	team1 := strings.Trim(inputLine[0:teamDelimPos[0][1]], " ")
	team1 = team1[0 : len(team1)-1]
	team2 := strings.Trim(inputLine[teamDelimPos[0][1]:], " ")
	teamNamesAndScores = append(teamNamesAndScores, team1)
	teamNamesAndScores = append(teamNamesAndScores, team2)

	if len(teamNamesAndScores) < 2 {
		return "", "", 0, 0, errors.New("invalid line, no commas to separate teams")
	} else {

		getTeamNameAndScore := func(nameAndScoreString string) (teamName string, teamScore int, scoreErr error) {
			nameScoreSplitRegex := regexp.MustCompile(`\s[0-9]+`)
			// println("name and score string: ", nameAndScoreString)
			scoreIndex := nameScoreSplitRegex.FindAllStringIndex(nameAndScoreString, -1)

			if scoreIndex == nil {
				return "", 0, errors.New("invalid line, could not find a team score")
			}
			//find last index of a space followed by a number
			// fmt.Printf("Score Index: %v", scoreIndex)

			// println(scoreIndex)
			teamName = nameAndScoreString[0:scoreIndex[len(scoreIndex)-1][0]]
			teamScore, scoreErr = strconv.Atoi(strings.TrimSpace(nameAndScoreString[scoreIndex[len(scoreIndex)-1][0]:scoreIndex[len(scoreIndex)-1][1]]))
			if scoreErr != nil {
				return "", 0, scoreErr
			}
			return teamName, teamScore, nil
		}

		teamName1, teamScore1, scoreErr := getTeamNameAndScore(teamNamesAndScores[0])
		if scoreErr != nil {
			return teamName1, teamName2, teamScore1, teamScore2, scoreErr
		}
		teamName2, teamScore2, scoreErr := getTeamNameAndScore(teamNamesAndScores[1])
		if scoreErr != nil {
			return teamName1, teamName2, teamScore1, teamScore2, scoreErr
		}
		return teamName1, teamName2, teamScore1, teamScore2, nil
	}

}
