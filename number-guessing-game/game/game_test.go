package game_test

import (
	"number-guessing-games/game"
	"testing"
)

func TestParseDifficulty(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    game.Difficulty
		wantErr bool
	}{
		{"lowercase easy", "easy", game.Easy, false},
		{"lowercase medium", "medium", game.Medium, false},
		{"lowercase hard", "hard", game.Hard, false},
		{"invalid", "xxx", 0, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := game.ParseDifficulty(c.input)

			gotErr := err != nil
			if gotErr != c.wantErr {
				t.Fatalf("error mismatch: got err=%v, want err=%v", err, c.wantErr)
			}

			if !c.wantErr && got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestNewGame(t *testing.T) {
	cases := []struct {
		name            string
		input           game.Difficulty
		wantMaxAttempts int
	}{
		{"easy", game.Easy, 10},
		{"medium", game.Medium, 5},
		{"hard", game.Hard, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := game.NewGame(c.input)

			if g.MaxAttempts != c.wantMaxAttempts {
				t.Errorf("MaxAttempts: got%d, want %d", g.MaxAttempts, c.wantMaxAttempts)
			}

			if g.SecretNumber < 1 || g.SecretNumber > 100 {
				t.Errorf("SecretNumber out of range: got %d", g.SecretNumber)
			}

			if g.Difficulty != c.input {
				t.Errorf("Difficulty: got %v, want %v", g.Difficulty, c.input)
			}
		})
	}
}

func TestGuess(t *testing.T) {

	cases := []struct {
		name  string
		input int
		want  game.GuessResult
	}{
		{"less than n", 9, game.TooLow},
		{"more than n", 11, game.TooHigh},
		{"correct", 10, game.Correct},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := game.Game{
				SecretNumber: 10,
				MaxAttempts:  10,
			}

			result := g.Guess(c.input)

			if result != c.want {
				t.Errorf("Result do not match, got: %v, want: %v", result, c.want)
			}

			if g.AttemptsUsed != 1 {
				t.Errorf("AttemptsUsed: got %d, want 1", g.AttemptsUsed)
			}

		})
	}
}

func TestAttemptsRemaining(t *testing.T) {
	cases := []struct {
		name    string
		maxAtt  int
		attUsed int
		want    int
	}{
		{"no guess at all", 10, 0, 10},
		{"three guess", 10, 3, 7},
		{"out of attempts", 10, 10, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := game.Game{
				MaxAttempts:  c.maxAtt,
				AttemptsUsed: c.attUsed,
			}
			remaining := g.AttemptsRemaining()
			if remaining != c.want {
				t.Errorf("AttemptsRemaining got: %d, want: %d", remaining, c.want)
			}
		})
	}
}

func TestHint(t *testing.T) {
	cases := []struct {
		name         string
		secretNumber int
		want         string
	}{
		{"even", 10, "even"},
		{"odd", 9, "odd"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := game.Game{
				SecretNumber: c.secretNumber,
			}
			result := g.Hint()
			if result != c.want {
				t.Errorf("Hint got: %v, want: %v", result, c.want)
			}
		})
	}
}
