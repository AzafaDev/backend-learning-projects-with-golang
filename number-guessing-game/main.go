package main

import (
	"bufio"
	"fmt"
	"number-guessing-games/game"
	"number-guessing-games/score"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		playRound(reader)

		fmt.Println("Play again? (y/n)")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error to read the input:", err)
			os.Exit(1)
		}
		answer := strings.ToLower(strings.TrimSpace(input))
		if answer != "y" && answer != "yes" {
			fmt.Println("See you!")
			break
		}
	}

}

func playRound(reader *bufio.Reader) {
	var difficulty game.Difficulty
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
		if result != game.Correct && g.AttemptsRemaining() <= g.MaxAttempts/2 {
			fmt.Println("Hint: the number is", g.Hint())
		}
		switch result {
		case game.TooLow:
			fmt.Println("Too low")
		case game.TooHigh:
			fmt.Println("Too High")
		case game.Correct:
			fmt.Printf("Correct! used %d attempts\n", g.AttemptsUsed)
			fmt.Println("Duration of the playing:", g.Elapsed().Round(time.Second))
			h, err := score.Load()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if ok := h.UpdateIfBetter(g.Difficulty, g.AttemptsUsed); ok {
				if err := score.Save(h); err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
			}

			fmt.Printf(`Your High Score
Easy: %d,
Medium: %d,
Hard: %d
			`, h.Easy, h.Medium, h.Hard)

			break GameLoop
		}
		if result != game.Correct && g.AttemptsRemaining() <= 0 {
			fmt.Println("Attempts are out! the answer is", g.SecretNumber)
			fmt.Println("Duration of the playing:", g.Elapsed().Round(time.Second))
			break GameLoop
		}

	}
}
