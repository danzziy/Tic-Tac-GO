package analyzer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidMoves(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		prevGameState string
		playerMove    string
	}{
		{"000000000", "000010000"},
		{"022110000", "022111000"},
		{"220110001", "222110001"},
	} {
		tc := tc
		t.Run(fmt.Sprintf("%s -> %s", tc.prevGameState, tc.playerMove), func(t *testing.T) {
			t.Parallel()

			// Act
			analyzer := NewAnalyzer()
			validMove, err := analyzer.ValidMove(tc.prevGameState, tc.playerMove)

			// Assert
			assert.NoError(t, err)
			assert.True(t, validMove)
		})
	}
}

func TestInValidMoves(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		prevGameState string
		playerMove    string
	}{
		{"000000000", "000000000"}, // No change
		{"000000000", "000020000"}, // Starting with 2
		// {"000020000", "000120000"}, I would like this condition to result in a fail but it will never happen so I won't implement it.
		{"000010000", "000110000"}, // Too many ones
		{"000210000", "000212000"}, // Too many twos
		{"000010000", "000020000"}, // Overwrite
	} {
		tc := tc
		t.Run(fmt.Sprintf("%s -> %s", tc.prevGameState, tc.playerMove), func(t *testing.T) {
			t.Parallel()

			// Act
			analyzer := NewAnalyzer()
			validMove, err := analyzer.ValidMove(tc.prevGameState, tc.playerMove)

			// Assert
			assert.NoError(t, err)
			assert.False(t, validMove)
		})
	}
}
