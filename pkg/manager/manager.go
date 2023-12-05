package manager

import (
	"github.com/google/uuid"
)

type Database interface {
	PublicRoomAvailable() bool
	CreatePublicRoom(roomID string, playerID string) error
}

type Manager interface {
	StartGame(message string) (GameRoom, error)
	MakePlayerMove(roomID string, message string) (GameRoom, error)
	EndGame(roomID string) GameRoom
}

type manager struct {
	database Database
}

func NewManager(database Database) Manager {
	return &manager{database}
}

func (m *manager) StartGame(message string) (GameRoom, error) {
	playerID := uuid.NewString()
	if !m.database.PublicRoomAvailable() {
		roomID := uuid.NewString()

		_ = m.database.CreatePublicRoom(roomID, playerID)
		return GameRoom{roomID, []Player{{playerID, "Waiting for Player"}}}, nil
	}
	return GameRoom{}, nil
}

func (m *manager) MakePlayerMove(roomID string, message string) (GameRoom, error) {
	return GameRoom{}, nil
}

func (m *manager) EndGame(roomID string) GameRoom {
	return GameRoom{}
}

type GameRoom struct {
	RoomID  string
	Players []Player
}

type Player struct {
	ID      string
	Message string
}
