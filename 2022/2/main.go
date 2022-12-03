package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed challenge.in
var input string

type ChallengeMode int64

const (
	rockOp     string = "A "
	paperOp    string = "B "
	scissorsOp string = "C "

	rock     string = "X"
	paper    string = "Y"
	scissors string = "Z"

	loss string = "X"
	draw string = "Y"
	win  string = "Z"
)

var outcomeValue = map[string]int{
	win:  6,
	draw: 3,
	loss: 0,
}

var objectValue = map[string]int{
	rock:     1,
	paper:    2,
	scissors: 3,
}

var objectMatrix = map[string]int{
	rockOp + rock:     outcomeValue[draw],
	rockOp + paper:    outcomeValue[win],
	rockOp + scissors: outcomeValue[loss],

	paperOp + rock:     outcomeValue[loss],
	paperOp + paper:    outcomeValue[draw],
	paperOp + scissors: outcomeValue[win],

	scissorsOp + rock:     outcomeValue[win],
	scissorsOp + paper:    outcomeValue[loss],
	scissorsOp + scissors: outcomeValue[draw],
}

var outcomeMatrix = map[string]int{
	rockOp + win:  objectValue[paper],
	rockOp + draw: objectValue[rock],
	rockOp + loss: objectValue[scissors],

	paperOp + win:  objectValue[scissors],
	paperOp + draw: objectValue[paper],
	paperOp + loss: objectValue[rock],

	scissorsOp + win:  objectValue[rock],
	scissorsOp + draw: objectValue[scissors],
	scissorsOp + loss: objectValue[paper],
}

func runChallenge(challengePart int) int {
	scoreSum := 0
	playPlan := strings.Split(input, "\n")
	for _, move := range playPlan {
		mine := string(move[2])
		if challengePart == 1 {
			scoreSum += objectMatrix[move] + objectValue[mine]
		} else {
			scoreSum += outcomeMatrix[move] + outcomeValue[mine]
		}
	}
	return scoreSum
}

func main() {
	highscore := runChallenge(2)
	fmt.Println(highscore)
}
