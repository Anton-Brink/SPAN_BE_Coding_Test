package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("input team results:")
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

			if result == "tie" {
				if _, exists := teamPoints[teamName1]; exists {
					teamPoints[teamName1] += 1
				} else {
					teamPoints[teamName1] += 1
				}
				if _, exists := teamPoints[teamName2]; exists {
					teamPoints[teamName2] += 1
				} else {
					teamPoints[teamName2] += 1
				}
			} else if result == teamName1 {
				if _, exists := teamPoints[teamName1]; exists {
					teamPoints[teamName1] += 3
				} else {
					teamPoints[teamName1] += 3
				}
			} else {
				if _, exists := teamPoints[teamName2]; exists {
					teamPoints[teamName2] += 3
				} else {
					teamPoints[teamName2] += 3
				}
			}
		}

	}

	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output: ")
	for team, score := range teamPoints {
		fmt.Println(team, " - ", score)
	}

	fmt.Println("The Result is also written to the results file")
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
	// do regex later for proof of concept stuff, teams with comma and or numbers in name
	// splitRegex := regexp.MustCompile("[0-9], [0-9a-zA-Z]");
	var bothNamesAndScores = strings.Split(inputLine, ",")

	if len(bothNamesAndScores) != 2 {
		return "", "", 0, 0, errors.New("invalid line, too many commas")
	} else {

		getTeamNameAndScore := func(nameAndScoreString string) (teamName string, teamScore int) {
			nameScoreSplitRegex := regexp.MustCompile(`\s[0-9]+`)
			teamNameAndScore := nameScoreSplitRegex.Split(nameAndScoreString, -1)
			// teamNameAndScore := strings.Split(nameAndScoreString,nameScoreSplitRegex);
			var scoreErr error
			if len(teamNameAndScore) > 2 {
				teamName = strings.Join(teamNameAndScore[0:len(teamNameAndScore)-2], ",")
				teamScore, scoreErr = strconv.Atoi(teamNameAndScore[len(teamNameAndScore)-1])
				if scoreErr != nil {
					return
				}
			} else if len(teamNameAndScore) == 2 {
				teamName = teamNameAndScore[0]
				teamScore, scoreErr = strconv.Atoi(teamNameAndScore[1])
			}
			return teamName, teamScore
		}

		teamName1, teamScore1 := getTeamNameAndScore(bothNamesAndScores[0])
		teamName2, teamScore2 := getTeamNameAndScore(bothNamesAndScores[1])
		return teamName1, teamName2, teamScore1, teamScore2, nil
	}

}
