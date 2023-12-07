package manager

import (
	"log"

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
}

type Manager interface {
	StartGame(message string) (GameRoom, error)
	ExecutePlayerMove(roomID string, message string) (GameRoom, error)
	EndGame(roomID string) GameRoom
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
func (m *manager) ExecutePlayerMove(roomID string, playerMove string) (GameRoom, error) {
	gameRoom, _ := m.database.RetrieveGame(roomID)
	prevGameState := gameRoom.Players[0].Message
	log.Printf("sasdfjas;klfjwepi\n\n %s", prevGameState)
	log.Printf("ALFJASLFasdfsadfdsfadsf %s", playerMove)

	validMove, _ := m.analyzer.ValidMove(prevGameState, playerMove)

	if validMove {
		_ = m.database.ExecutePlayerMove(roomID, playerMove)
		players, _ := m.analyzer.DetermineWinner(playerMove, gameRoom.Players)
		return GameRoom{gameRoom.RoomID, players}, nil
	}
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
