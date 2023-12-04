package manager

type Manager interface {
	StartGame(message string) (GameRoom, error)
	MakePlayerMove()
	EndGame()
}

type manager struct {
}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) StartGame(message string) (GameRoom, error) {
	return GameRoom{}, nil
}

func (m *manager) MakePlayerMove() {

}

func (m *manager) EndGame() {

}

type GameRoom struct {
	RoomID    string
	PlayerIDs []string
	Message   string
}
