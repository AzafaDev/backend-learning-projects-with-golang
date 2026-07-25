package score_test

import (
	"number-guessing-games/game"
	"number-guessing-games/score"
	"testing"
)

func TestUpdateIfBetter(t *testing.T) {
	cases := []struct {
		name       string
		initial    score.HighScores // state awal sebelum panggil
		difficulty game.Difficulty
		attempts   int
		wantOk     bool
		wantScore  int // nilai field yang relevan SETELAH panggilan
	}{
		{"first record easy", score.HighScores{}, game.Easy, 5, true, 5},
		{"worse attempt not updated", score.HighScores{Easy: 3}, game.Easy, 5, false, 3},
		{"better attempt updated", score.HighScores{Easy: 5}, game.Easy, 3, true, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := c.initial
			ok := h.UpdateIfBetter(c.difficulty, c.attempts)
			if ok != c.wantOk {
				t.Errorf("Error OK got: %v, wantOK: %v", ok, c.wantOk)
			}

			if h.Easy != c.wantScore {
				t.Errorf("Error score got: %d, want: %d", h.Easy, c.wantScore)
			}
		})
	}
}

func TestUpdateIfBetterDoesNotAffectOtherFields(t *testing.T) {
	h := score.HighScores{Medium: 4, Hard: 2}

	h.UpdateIfBetter(game.Easy, 5)

	if h.Medium != 4 {
		t.Errorf("Medium changed unexpectedly: got %d, want 4", h.Medium)
	}
	if h.Hard != 2 {
		t.Errorf("Hard changed unexpectedly: got %d, want 2", h.Hard)
	}
}
