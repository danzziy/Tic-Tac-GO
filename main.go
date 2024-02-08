package main

import (
	"os"
	"tic-tac-go/pkg/analyzer"
	"tic-tac-go/pkg/config"
	"tic-tac-go/pkg/database"
	"tic-tac-go/pkg/manager"
	game "tic-tac-go/pkg/server"
)

func main() {
	server := initializeGameServer()
	_ = server.Start()
	defer func() { _ = server.Stop() }()
}

func initializeGameServer() *game.HTTPServer {
	env := config.NewConfig(
		os.Environ(),
	)
	port, _ := env.ListeningPort()
	databaseHost, _ := env.DatabaseHost()
	databasePassword, _ := env.DatabasePassword()
	return game.NewHTTPServer(port, manager.NewManager(
		database.NewDatabase(databaseHost, databasePassword), analyzer.NewAnalyzer()),
	)
}
