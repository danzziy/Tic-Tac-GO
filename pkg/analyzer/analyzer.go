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

func (a *analyzer) ValidMove(prevGameState string, playerMove string) bool {
	if len(prevGameState) != 9 || len(playerMove) != 9 || prevGameState == playerMove {
		return false
	}

	ones := 0
	twos := 0
	moves := 0
	for i := 0; i < len(prevGameState); i++ {
		if playerMove[i] == '1' {
			ones++
		} else if playerMove[i] == '2' {
			twos++
		}

		if prevGameState[i] != playerMove[i] {
			moves++
			// Valid move: From '0' (empty) to '1' or '2' (player marker)
			if prevGameState[i] == '1' && playerMove[i] == '2' || prevGameState[i] == '2' && playerMove[i] == '1' || moves > 1 {
				return false // Invalid move
			}
		}
	}

	if twos > ones || ones > twos+1 {
		return false // Invalid move
	}

	return true
}

func (a *analyzer) DetermineWinner(playerMove string, players []manager.Player) []manager.Player {
	if rowsContainAWinner(playerMove, &players) || columnsContainAWinner(playerMove, &players) ||
		diagonalsContainAWinner(playerMove, &players) {
		return players
	} else if !strings.Contains(playerMove, "0") {
		players[0].Message = fmt.Sprintf("%s:Tie", playerMove)
		players[1].Message = fmt.Sprintf("%s:Tie", playerMove)
		return players
	}

	players[0].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	players[1].Message = fmt.Sprintf("%s:Ongoing", playerMove)
	return players
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
