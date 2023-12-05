package database

import "tic-tac-go/pkg/manager"

type database struct {
}

func NewDatabase() manager.Database {
	return &database{}
}

func (d *database) PublicRoomAvailable() bool {
	return false
}

func (d *database) CreatePublicRoom(roomID string, playerID string) error {
	return nil
}
