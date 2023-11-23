package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortNumber(t *testing.T) {
	for _, expectedPort := range []int{8080, 6969} {
		expectedPort := expectedPort
		t.Run(fmt.Sprintf("with port %d", expectedPort), func(t *testing.T) {
			t.Parallel()

			// Act
			config := NewConfig([]string{fmt.Sprintf("LISTENING_PORT=%d", expectedPort)})
			actualPort, _ := config.GetListeningPort()

			// Assert
			assert.Equal(t, actualPort, expectedPort)
		})
	}
}
