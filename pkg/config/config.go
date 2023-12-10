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

func (c config) ListeningPort() (int, error) {
	env, _ := c.env["LISTENING_PORT"]
	port, _ := strconv.Atoi(env)
	return port, nil
}

func (c config) DatabaseHost() (string, error) {
	host, _ := c.env["DATABASE_HOST"]
	return host, nil
}

func (c config) DatabasePassword() (string, error) {
	password, _ := c.env["DATABASE_PASSWORD"]
	return password, nil
}
