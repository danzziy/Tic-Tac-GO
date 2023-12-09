package analyzer

import (
	"fmt"
	"strings"
	"tic-tac-go/pkg/manager"
)

type analyzer struct {
}

func NewAnalyzer() manager.Analyzer {
	return &analyzer{}
}

func (a *analyzer) ValidMove(prevGameState string, playerMove string) (bool, error) {
	// TODO: If the length is not 9 there should be an error.
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
	if rowsContainAWinner(playerMove, &players) || columnsContainAWinner(playerMove, &players) ||
		diagonalsContainAWinner(playerMove, &players) {
		return players, nil
	} else if !strings.Contains(playerMove, "0") {
		players[0].Message = fmt.Sprintf("%s:Tie", playerMove)
		players[1].Message = fmt.Sprintf("%s:Tie", playerMove)
		return players, nil
	}

	players[0].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	players[1].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	return players, nil
}

func rowsContainAWinner(playerMove string, players *[]manager.Player) bool {
	for i := 0; i < 9; i += 3 {
		if isAWinningSegment(playerMove[i:i+3], playerMove, players) {
			return true
		}
	}
	return false
}

func columnsContainAWinner(playerMove string, players *[]manager.Player) bool {
	for i := 0; i < 3; i++ {
		verticalLine := fmt.Sprintf("%c%c%c", playerMove[i], playerMove[i+3], playerMove[i+6])
		if isAWinningSegment(verticalLine, playerMove, players) {
			return true
		}
	}
	return false
}

func diagonalsContainAWinner(playerMove string, players *[]manager.Player) bool {
	for i := 0; i <= 2; i += 2 {
		diagonalLine := fmt.Sprintf("%c%c%c", playerMove[i], playerMove[4], playerMove[8-i])
		if isAWinningSegment(diagonalLine, playerMove, players) {
			return true
		}
	}
	return false
}

func isAWinningSegment(line string, playerMove string, players *[]manager.Player) bool {
	if line == "111" {
		(*players)[0].Message = fmt.Sprintf("%s:Win", playerMove)
		(*players)[1].Message = fmt.Sprintf("%s:Lose", playerMove)
		return true
	} else if line == "222" {
		(*players)[0].Message = fmt.Sprintf("%s:Lose", playerMove)
		(*players)[1].Message = fmt.Sprintf("%s:Win", playerMove)
		return true
	}
	return false
}
