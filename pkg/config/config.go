package config

import (
	"strconv"
	"strings"
)

type config struct {
	env map[string]string
}

func NewConfig(env []string) *config {
	mappedEnv := make(map[string]string)
	for _, t := range env {
		key, value, _ := strings.Cut(t, "=")
		mappedEnv[key] = value
	}
	return &config{mappedEnv}
}

func (c config) GetListeningPort() (int, error) {

	env := c.env["LISTENING_PORT"]
	port, _ := strconv.Atoi(env)

	return port, nil
}
