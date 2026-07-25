package score

import (
	"encoding/json"
	"number-guessing-games/game"
	"os"
)

type HighScores struct {
	Easy   int `json:"easy"`
	Medium int `json:"medium"`
	Hard   int `json:"hard"`
}

const filePath = "score/highscore.json"

func Load() (HighScores, error) {
	dataByte, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return HighScores{}, nil
	} else if err != nil {
		return HighScores{}, err
	}
	var highScores HighScores
	if err := json.Unmarshal(dataByte, &highScores); err != nil {
		return highScores, err
	}
	return highScores, nil
}

func Save(scores HighScores) error {
	dataByte, err := json.MarshalIndent(scores, "", " ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, dataByte, 0644); err != nil {
		return err
	}
	return nil
}

func (h *HighScores) UpdateIfBetter(difficulty game.Difficulty, attempts int) bool {
	switch difficulty {
	case game.Easy:
		if h.Easy == 0 || attempts < h.Easy {
			h.Easy = attempts
			return true
		}
	case game.Medium:
		if h.Medium == 0 || attempts < h.Medium {
			h.Medium = attempts
			return true
		}
	case game.Hard:
		if h.Hard == 0 || attempts < h.Hard {
			h.Hard = attempts
			return true
		}
	default:
		return false
	}
	return false
}
