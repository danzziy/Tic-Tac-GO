package analyzer

import (
	"fmt"
	"testing"
	"tic-tac-go/pkg/manager"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TODO: None of these need to return an error.
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

// TODO: Rename playerMove to gamestate for more clarity.
// TODO: DetermineWinner seems overloaded, perhaps moving the append logic into manager will simplify it.
func TestDeterminesWinner(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		playerMove    string
		player1Result string
		player2Result string
	}{
		{"000000000", "000000000:Ongoing", "000000000:Ongoing"}, // No change
		{"111220000", "111220000:Win", "111220000:Lose"},        // P1 wins first row
		{"222110100", "222110100:Lose", "222110100:Win"},        // P2 wins first row
		{"220111000", "220111000:Win", "220111000:Lose"},        // P1 wins second row
		{"110222100", "110222100:Lose", "110222100:Win"},        // P2 wins second row
		{"200022111", "200022111:Win", "200022111:Lose"},        // P1 wins third row
		{"100011222", "100011222:Lose", "100011222:Win"},        // P2 wins third row
		{"122100100", "122100100:Win", "122100100:Lose"},        // P1 wins first column
		{"211210200", "211210200:Lose", "211210200:Win"},        // P2 wins first column
		{"210210010", "210210010:Win", "210210010:Lose"},        // P1 wins second column
		{"121120020", "121120020:Lose", "121120020:Win"},        // P2 wins second column
		{"001001221", "001001221:Win", "001001221:Lose"},        // P1 wins third column
		{"002012112", "002012112:Lose", "002012112:Win"},        // P2 wins third column
		{"120210001", "120210001:Win", "120210001:Lose"},        // P1 wins left to right diagonal
		{"211120002", "211120002:Lose", "211120002:Win"},        // P2 wins left to right diagonal
		{"201210100", "201210100:Win", "201210100:Lose"},        // P1 wins right to left diagonal
		{"102021210", "102021210:Lose", "102021210:Win"},        // P2 wins right to left diagonal
		{"221112211", "221112211:Tie", "221112211:Tie"},         // P2 wins right to left diagonal
	} {
		tc := tc
		t.Run(fmt.Sprintf("game state %s", tc.playerMove), func(t *testing.T) {
			t.Parallel()

			player1ID := uuid.NewString()
			player2ID := uuid.NewString()

			players := []manager.Player{{ID: player1ID, Message: ""}, {ID: player2ID, Message: ""}}
			expectedPlayers := []manager.Player{{ID: player1ID, Message: tc.player1Result}, {ID: player2ID, Message: tc.player2Result}}

			// Act
			analyzer := NewAnalyzer()
			actualPlayers, err := analyzer.DetermineWinner(tc.playerMove, players)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedPlayers, actualPlayers)
		})
	}
}
