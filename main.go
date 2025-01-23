package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("input team results:")
	scanner := bufio.NewScanner(os.Stdin)

	var teamPoints map[string]int
	counter := 0

	for {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 || line == "end" {
			break
		}
		counter++
		team1, team1Score, team2, team2Score, calcErr := calculateTeamPoints(line)
		if calcErr != nil {
			//inform user of error and how to fix it
			fmt.Println("Line ", counter, " has an error and was not added to the tally, please ensure your format for every line is team1 score, team2 score")
		} else {
			//add values to teams in teamPoints map
			if teamScore, ok := teamPoints[team1]; ok {
				teamPoints[team1] = teamScore + team1Score
			}
			if teamScore, ok := teamPoints[team2]; ok {
				teamPoints[team2] = teamScore + team2Score
			}
		}

	}

	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output: ")
	// for team,_ := range teamPoints {

	// }

	fmt.Println("The Result is also written to the results file")
}

func calculateTeamPoints(inputLine string) (team1 string, team1Score int, team2 string, team2Score int, calculationError error) {
	// do regex later for proof of concept stuff, teams with comma and or numbers in name
	// splitRegex := regexp.MustCompile("[0-9], [0-9a-zA-Z]");
	var teamNamesAndScores = strings.Split(inputLine, ",")
	if len(teamNamesAndScores) > 2 {
		return "", 0, "", 0, errors.New("invalid line")
	}

}
