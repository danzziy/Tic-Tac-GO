package analyzer

import "tic-tac-go/pkg/manager"

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
	return []manager.Player{}, nil
}
