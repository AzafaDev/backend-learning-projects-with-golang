package game

import (
	"fmt"
	"math/rand"
	"strings"
)

type Game struct {
	SecretNumber int
	MaxAttempts  int
	AttemptsUsed int
	Difficulty   Difficulty
}

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

type GuessResult int

const (
	TooLow GuessResult = iota
	TooHigh
	Correct
)

func NewGame(difficulty Difficulty) *Game {
	secretNumber := rand.Intn(100) + 1
	maxAttempts := 0
	switch difficulty {
	case Easy:
		maxAttempts = 10
	case Medium:
		maxAttempts = 5
	case Hard:
		maxAttempts = 3
	default:
		panic(fmt.Errorf("NewGame: unknown difficulty value %d", difficulty))
	}
	return &Game{
		SecretNumber: secretNumber,
		MaxAttempts:  maxAttempts,
		AttemptsUsed: 0,
		Difficulty:   difficulty,
	}
}

func (g *Game) Guess(n int) GuessResult {
	g.AttemptsUsed++
	switch {
	case n < g.SecretNumber:
		return TooLow
	case n > g.SecretNumber:
		return TooHigh
	default:
		return Correct
	}
}

func (g *Game) AttemptsRemaining() int {
	return g.MaxAttempts - g.AttemptsUsed
}

func ParseDifficulty(s string) (Difficulty, error) {
	input := strings.ToLower(strings.TrimSpace(s))
	switch input {
	case "easy":
		return Easy, nil
	case "medium":
		return Medium, nil
	case "hard":
		return Hard, nil
	default:
		return 0, fmt.Errorf("unknown difficulty: %q", s)
	}
}
