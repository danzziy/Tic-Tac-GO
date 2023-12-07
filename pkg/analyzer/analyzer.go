package analyzer

import "tic-tac-go/pkg/manager"

type analyzer struct {
}

func (a *analyzer) ValidGameState(prevGameState string, playerMove string) (bool, error) {
	return false, nil
}
func (a *analyzer) DetermineWinner(playerMove string, players []manager.Player) ([]manager.Player, error) {
	return []manager.Player{}, nil
}
