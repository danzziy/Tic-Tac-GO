package manager

import (
	"github.com/google/uuid"
)

type Analyzer interface {
	ValidMove(prevGameState string, playerMove string) (bool, error)
	DetermineWinner(playerMove string, players []Player) ([]Player, error)
}

type Database interface {
	PublicRoomAvailable() (bool, error)
	CreatePublicRoom(roomID string, playerID string) error
	JoinPublicRoom(playerID string) (string, string, error)
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

func (m *manager) StartGame(message string) (GameRoom, error) {
	playerID := uuid.NewString()
	roomAvailable, _ := m.database.PublicRoomAvailable()
	if !roomAvailable {
		roomID := uuid.NewString()

		_ = m.database.CreatePublicRoom(roomID, playerID)
		return GameRoom{roomID, []Player{{playerID, "Waiting for Player"}}}, nil
	}
	roomID, opponentID, _ := m.database.JoinPublicRoom(playerID)
	return GameRoom{roomID, []Player{{opponentID, "Start Game"}, {playerID, "Start Game"}}}, nil
}

// TODO: Consider having another object to store players and gamestate and exclude the message.
// Perhaps call it a GameRoomState and rename the other to GameRoomMessenger.
// TODO: If an opponent sends a websocket message prior to the current players turn, they could hijack
// their opponents move. Perhaps have each user send in their ids to verify themselves.
func (m *manager) ExecutePlayerMove(roomID string, playerMove string) (GameRoom, error) {
	gameRoom, _ := m.database.RetrieveGame(roomID)
	prevGameState := gameRoom.Players[0].Message

	validMove, _ := m.analyzer.ValidMove(prevGameState, playerMove)

	if validMove {
		_ = m.database.ExecutePlayerMove(roomID, playerMove)
		players, _ := m.analyzer.DetermineWinner(playerMove, gameRoom.Players)
		return GameRoom{gameRoom.RoomID, players}, nil
	}
	return GameRoom{}, nil
}

func (m *manager) EndGame(roomID string) (GameRoom, error) {
	gameRoom, _ := m.database.RetrieveGame(roomID)
	_ = m.database.DeleteGameRoom(roomID)

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
