package main

import (
	"bufio"
	"fmt"
	"number-guessing-games/game"
	"os"
	"strconv"
	"strings"
)

func main() {
	var difficulty game.Difficulty
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Choose difficulty level (easy/medium/hard): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error to read the input:", err)
			continue
		}
		d, err := game.ParseDifficulty(input)
		if err != nil {
			fmt.Println("Error in parsing difficulty:", err)
			continue
		}
		difficulty = d
		break
	}
	g := game.NewGame(difficulty)
GameLoop:
	for {
		fmt.Println("Guess number (1-100)")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error to read the input:", err)
			os.Exit(1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Error in converting the guess number:", err)
			continue
		}
		result := g.Guess(n)
		switch result {
		case game.TooLow:
			fmt.Println("Too low")
		case game.TooHigh:
			fmt.Println("Too High")
		case game.Correct:
			fmt.Printf("Correct! used %d attempts\n", g.AttemptsUsed)
			break GameLoop
		}
		if result != game.Correct && g.AttemptsRemaining() <= 0 {
			fmt.Println("Attempts are out! the answer is", g.SecretNumber)
			break GameLoop
		}

	}
}
