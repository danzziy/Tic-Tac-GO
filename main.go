package main

import (
	"fmt"
	"tic-tac-go/pkg/analyzer"
	"tic-tac-go/pkg/config"
	"tic-tac-go/pkg/database"
	"tic-tac-go/pkg/manager"
	game "tic-tac-go/pkg/server"
)

func main() {
	server := initializeGameServer("")
	_ = server.Start()
	defer func() { _ = server.Stop() }()
}

func initializeGameServer(dbAddr string) *game.HTTPServer {
	env := config.NewConfig(
		[]string{"LISTENING_PORT=8080", fmt.Sprintf("DATABASE_HOST=%s", dbAddr), "DATABASE_PASSWORD=something"},
	)
	port, _ := env.ListeningPort()
	databaseHost, _ := env.DatabaseHost()
	databasePassword, _ := env.DatabasePassword()
	return game.NewHTTPServer(port, manager.NewManager(
		database.NewDatabase(databaseHost, databasePassword), analyzer.NewAnalyzer()),
	)
}
