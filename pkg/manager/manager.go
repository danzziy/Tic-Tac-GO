package manager

import (
	"github.com/google/uuid"
)

type Analyzer interface {
	ValidMove(prevGameState string, playerMove string) bool
	DetermineWinner(playerMove string, players []Player) []Player
}

type Database interface {
	PublicRoomAvailable() (bool, error)
	CreatePublicRoom(roomID string, playerID string) error
	JoinPublicRoom(playerID string) (string, error)
	RetrieveGame(roomID string) (GameRoom, error)
	ExecutePlayerMove(roomID string, playerMove string) error
	DeleteGameRoom(roomID string) error
}

type Manager interface {
	StartGame(message string) (GameRoom, error)
	ExecutePlayerMove(roomID string, message string) (GameRoom, error)
	EndGame(roomID string) (GameRoom, error)
}

type manager struct {
	database Database
	analyzer Analyzer
}

func NewManager(database Database, analyzer Analyzer) Manager {
	return &manager{database, analyzer}
}

// TODO: Consider having another object to store players and gamestate and exclude the message.
// Perhaps call it a GameRoomState and rename the other to GameRoomMessenger.
// TODO: If an opponent sends a websocket message prior to the current players turn, they could hijack
// their opponents move. Perhaps have each user send in their ids to verify themselves.
// TODO: Implement retry logic.
// TODO: Implement error handling tests.

func (m *manager) StartGame(message string) (GameRoom, error) {
	playerID := uuid.NewString()

	roomAvailable, err := m.database.PublicRoomAvailable()
	if err != nil {
		return GameRoom{}, err
	}

	if !roomAvailable {
		roomID := uuid.NewString()
		if err = m.database.CreatePublicRoom(roomID, playerID); err != nil {
			return GameRoom{}, err
		}
		return GameRoom{roomID, []Player{{playerID, "Waiting for Player"}}}, nil
	}

	roomID, err := m.database.JoinPublicRoom(playerID)
	if err != nil {
		return GameRoom{}, err
	}
	gameRoom, err := m.database.RetrieveGame(roomID)
	if err != nil {
		return GameRoom{}, err
	}

	gameRoom.Players[0].Message = "Start Game"
	gameRoom.Players[1].Message = "Start Game"

	return gameRoom, nil
}

func (m *manager) ExecutePlayerMove(roomID string, playerMove string) (GameRoom, error) {
	gameRoom, err := m.database.RetrieveGame(roomID)
	if err != nil {
		return GameRoom{}, err
	}
	prevGameState := gameRoom.Players[0].Message

	validMove := m.analyzer.ValidMove(prevGameState, playerMove)

	if validMove {
		if err := m.database.ExecutePlayerMove(roomID, playerMove); err != nil {
			return GameRoom{}, err
		}
		players := m.analyzer.DetermineWinner(playerMove, gameRoom.Players)
		return GameRoom{gameRoom.RoomID, players}, nil
	}
	return GameRoom{}, nil
}

func (m *manager) EndGame(roomID string) (GameRoom, error) {
	gameRoom, err := m.database.RetrieveGame(roomID)
	if err != nil {
		return GameRoom{}, err
	}
	if err := m.database.DeleteGameRoom(roomID); err != nil {
		return GameRoom{}, err
	}

	gameRoom.Players[0].Message = "Terminate Connection"
	gameRoom.Players[1].Message = "Terminate Connection"

	return gameRoom, nil
}

type GameRoom struct {
	RoomID  string
	Players []Player
}

type Player struct {
	ID      string
	Message string
}
