package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrievesListeningPort(t *testing.T) {
	for _, expectedPort := range []int{8080, 6969} {
		expectedPort := expectedPort
		t.Run(fmt.Sprintf("with port %d", expectedPort), func(t *testing.T) {
			t.Parallel()

			// Act
			config := NewConfig([]string{fmt.Sprintf("LISTENING_PORT=%d", expectedPort)})
			actualPort, err := config.ListeningPort()

			// Assert
			assert.Equal(t, actualPort, expectedPort)
			assert.NoError(t, err)
		})
	}
}

func TestRetrievesDatabaseHost(t *testing.T) {
	for _, expectedHost := range []string{"localhost:6379", "127.0.0.1:2750"} {
		expectedHost := expectedHost
		t.Run(fmt.Sprintf("with host %s", expectedHost), func(t *testing.T) {
			t.Parallel()

			// Act
			config := NewConfig([]string{fmt.Sprintf("DATABASE_HOST=%s", expectedHost)})
			actualHost, err := config.DatabaseHost()

			// Assert
			assert.Equal(t, expectedHost, actualHost)
			assert.NoError(t, err)
		})
	}
}

func TestRetrievesDatabasePassword(t *testing.T) {
	for _, expectedPassword := range []string{"", "some_password"} {
		expectedPassword := expectedPassword
		t.Run(fmt.Sprintf("with host %s", expectedPassword), func(t *testing.T) {
			t.Parallel()

			// Act
			config := NewConfig([]string{fmt.Sprintf("DATABASE_PASSWORD=%s", expectedPassword)})
			actualPassword, err := config.DatabasePassword()

			// Assert
			assert.Equal(t, expectedPassword, actualPassword)
			assert.NoError(t, err)
		})
	}
}
