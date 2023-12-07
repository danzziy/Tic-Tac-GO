package analyzer

import (
	"fmt"
	"log"
	"tic-tac-go/pkg/manager"
)

type analyzer struct {
}

func NewAnalyzer() manager.Analyzer {
	return &analyzer{}
}

func (a *analyzer) ValidMove(prevGameState string, playerMove string) (bool, error) {
	if len(prevGameState) != 9 || len(playerMove) != 9 || prevGameState == playerMove {
		return false, nil
	}

	ones := 0
	twos := 0

	for i := 0; i < len(prevGameState); i++ {
		if playerMove[i] == '1' {
			ones++
		} else if playerMove[i] == '2' {
			twos++
		}

		if prevGameState[i] != playerMove[i] {
			// Valid move: From '0' (empty) to '1' or '2' (player marker)
			if prevGameState[i] == '1' && playerMove[i] == '2' || prevGameState[i] == '2' && playerMove[i] == '1' {
				return false, nil // Invalid move
			}
		}
	}

	if twos > ones || ones > twos+1 {
		return false, nil // Invalid move
	}

	return true, nil
}

func (a *analyzer) DetermineWinner(playerMove string, players []manager.Player) ([]manager.Player, error) {
	for i := 0; i < 9; i += 3 {
		if playerMove[i:i+3] == "111" {
			players[0].Message = fmt.Sprintf("%s:Win", playerMove)
			players[1].Message = fmt.Sprintf("%s:Lose", playerMove)
			return players, nil
		} else if playerMove[i:i+3] == "222" {
			players[0].Message = fmt.Sprintf("%s:Lose", playerMove)
			players[1].Message = fmt.Sprintf("%s:Win", playerMove)
			return players, nil
		}
	}
	for i := 0; i < 3; i++ {
		verticalLine := fmt.Sprintf("%c%c%c", playerMove[i], playerMove[i+3], playerMove[i+6])
		log.Printf("VErtical: %s", verticalLine)

		if verticalLine == "111" {
			players[0].Message = fmt.Sprintf("%s:Win", playerMove)
			players[1].Message = fmt.Sprintf("%s:Lose", playerMove)
			return players, nil
		} else if verticalLine == "222" {
			players[0].Message = fmt.Sprintf("%s:Lose", playerMove)
			players[1].Message = fmt.Sprintf("%s:Win", playerMove)
			return players, nil
		}
	}
	for i := 0; i <= 2; i += 2 {
		diagonalLine := fmt.Sprintf("%c%c%c", playerMove[i], playerMove[4], playerMove[8-i])
		log.Printf("%s", diagonalLine)
		if diagonalLine == "111" {
			players[0].Message = fmt.Sprintf("%s:Win", playerMove)
			players[1].Message = fmt.Sprintf("%s:Lose", playerMove)
			return players, nil
		} else if diagonalLine == "222" {
			players[0].Message = fmt.Sprintf("%s:Lose", playerMove)
			players[1].Message = fmt.Sprintf("%s:Win", playerMove)
			return players, nil
		}
	}
	players[0].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	players[1].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	return players, nil

	checkHorizontalLines(playerMove, players)
	checkVerticalLines(playerMove, players)
	checkDiagonalLines(playerMove, players)
	players[0].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	players[1].Message = fmt.Sprintf("%s:Ongoing", playerMove)

	return []manager.Player{}, nil
}

func checkHorizontalLines(playerMove string, players []manager.Player) []manager.Player {
	for i := 0; i < 9; i += 3 {
		if playerMove[i:i+3] == "111" {
			players[0].Message = fmt.Sprintf("%s:Win", playerMove)
			players[1].Message = fmt.Sprintf("%s:Lose", playerMove)
			return players
		} else if playerMove[i:i+3] == "222" {
			players[0].Message = fmt.Sprintf("%s:Lose", playerMove)
			players[1].Message = fmt.Sprintf("%s:Win", playerMove)
			return players
		}
	}
	players[0].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	players[1].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	return players
}

func checkVerticalLines(playerMove string, players []manager.Player) []manager.Player {
	for i := 0; i < 3; i++ {
		verticalLine := fmt.Sprintf("%c%c%c", playerMove[i], playerMove[i+3], playerMove[i+6])
		log.Printf("VErtical: %s", verticalLine)

		if verticalLine == "111" {
			players[0].Message = fmt.Sprintf("%s:Win", playerMove)
			players[1].Message = fmt.Sprintf("%s:Lose", playerMove)
			return players
		} else if verticalLine == "222" {
			players[0].Message = fmt.Sprintf("%s:Lose", playerMove)
			players[1].Message = fmt.Sprintf("%s:Win", playerMove)
			return players
		}
	}
	players[0].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	players[1].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	return players
}

func checkDiagonalLines(playerMove string, players []manager.Player) []manager.Player {
	for i := 0; i <= 2; i += 2 {
		diagonalLine := fmt.Sprintf("%c%c%c", playerMove[i], playerMove[4], playerMove[i+6])
		log.Printf("%s", diagonalLine)
		if diagonalLine == "111" {
			players[0].Message = fmt.Sprintf("%s:Win", playerMove)
			players[1].Message = fmt.Sprintf("%s:Lose", playerMove)
			return players
		} else if diagonalLine == "222" {
			players[0].Message = fmt.Sprintf("%s:Lose", playerMove)
			players[1].Message = fmt.Sprintf("%s:Win", playerMove)
			return players
		}
	}
	players[0].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	players[1].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	return players
}

// func determinSegmentWinner(line string, playerMove string, players []manager.Player) []manager.Player {
// 	if line == "111" {
// 		players[0].Message = fmt.Sprintf("%s:Win", playerMove)
// 		players[1].Message = fmt.Sprintf("%s:Lose", playerMove)
// 		return players
// 	} else if line == "222" {
// 		players[0].Message = fmt.Sprintf("%s:Lose", playerMove)
// 		players[1].Message = fmt.Sprintf("%s:Win", playerMove)
// 		return players
// 	}
// 	return players
// }
