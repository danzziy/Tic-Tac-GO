package manager

type Manager interface {
	StartGame(message string) (GameRoom, error)
	MakePlayerMove(roomID string, message string) (GameRoom, error)
	EndGame(roomID string) GameRoom
}

type manager struct {
}

func NewManager() Manager {
	return &manager{}
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
	RoomID    string
	PlayerIDs []string
	Message   []string
}
