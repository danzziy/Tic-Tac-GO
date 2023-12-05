package database

import "tic-tac-go/pkg/manager"

type database struct {
}

func NewDatabase() manager.Database {
	return &database{}
}
