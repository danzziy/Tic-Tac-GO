package manager

type Database interface {
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
