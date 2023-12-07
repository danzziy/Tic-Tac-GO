package database

import "tic-tac-go/pkg/manager"

type database struct {
}

func NewDatabase() manager.Database {
	return &database{}
}

func (d *database) PublicRoomAvailable() (bool, error) {
	return false, nil
}

func (d *database) CreatePublicRoom(roomID string, playerID string) error {
	return nil
}

func (d *database) JoinPublicRoom(playerID string) (string, string, error) {
	return "", "", nil
}

func (d *database) RetrieveGameState(roomID string) (string, error) {
	return "", nil
}

func (d *database) ExecutePlayerMove(GameRoom string, roomID string) error {
	return nil
}
