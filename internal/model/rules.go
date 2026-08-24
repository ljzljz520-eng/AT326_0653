package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingID       = errors.New("missing identifier")
	ErrInvalidNickname = errors.New("nickname must be 2-24 characters")
	ErrInvalidScore    = errors.New("score must be positive")
	ErrInvalidMachine  = errors.New("machine must be published before use")
	ErrInvalidStatus   = errors.New("invalid schedule status")
	ErrCapacityReached = errors.New("schedule capacity reached")
	ErrPrizeExhausted  = errors.New("prize quantity exhausted")
)

func ValidatePlayer(p PlayerProfile) error {
	if PlayerKey(p.Nickname) == "" {
		return ErrMissingID
	}
	length := len([]rune(strings.TrimSpace(p.Nickname)))
	if length < 2 || length > 24 {
		return ErrInvalidNickname
	}
	if p.ID == "" {
		return ErrMissingID
	}
	return nil
}

func ValidateMachine(m ArcadeMachine) error {
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.Name) == "" {
		return ErrMissingID
	}
	if strings.TrimSpace(m.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(m.Venue) == "" {
		return fmt.Errorf("venue is required")
	}
	if strings.TrimSpace(m.Difficulty) == "" {
		return fmt.Errorf("difficulty is required")
	}
	return nil
}

func ValidateScore(s HighScoreEntry) error {
	if s.ID == "" || s.MachineID == "" || s.PlayerID == "" {
		return ErrMissingID
	}
	if s.Score <= 0 {
		return ErrInvalidScore
	}
	if strings.TrimSpace(s.RecordedOn) == "" {
		return fmt.Errorf("recorded_on is required")
	}
	return nil
}

func ValidateSchedule(s MatchSchedule) error {
	if s.ID == "" || strings.TrimSpace(s.Title) == "" {
		return ErrMissingID
	}
	if len(s.MachineIDs) == 0 {
		return fmt.Errorf("at least one machine is required")
	}
	if s.Capacity < 1 {
		return fmt.Errorf("capacity must be positive")
	}
	if s.Status == "" {
		return fmt.Errorf("status is required")
	}
	if !ValidStatus(s.Status) {
		return ErrInvalidStatus
	}
	if strings.TrimSpace(s.StartSlot) == "" || strings.TrimSpace(s.EndSlot) == "" {
		return fmt.Errorf("schedule slots are required")
	}
	return nil
}

func ValidatePrize(p PrizeDescription) error {
	if p.ID == "" || strings.TrimSpace(p.Name) == "" {
		return ErrMissingID
	}
	if strings.TrimSpace(p.Tier) == "" || strings.TrimSpace(p.Description) == "" {
		return fmt.Errorf("prize metadata is required")
	}
	if p.Quantity < 0 || p.Claimed < 0 || p.Claimed > p.Quantity {
		return fmt.Errorf("invalid prize inventory")
	}
	return nil
}

func ValidStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "draft", "open", "running", "completed", "cancelled":
		return true
	default:
		return false
	}
}

func CanTransition(from, to string) bool {
	from, to = strings.ToLower(strings.TrimSpace(from)), strings.ToLower(strings.TrimSpace(to))
	if !ValidStatus(from) || !ValidStatus(to) {
		return false
	}
	switch from {
	case "draft":
		return to == "open" || to == "cancelled"
	case "open":
		return to == "running" || to == "cancelled"
	case "running":
		return to == "completed" || to == "cancelled"
	case "completed", "cancelled":
		return false
	default:
		return false
	}
}

func ScoreRank(entries []HighScoreEntry) []HighScoreEntry {
	result := append([]HighScoreEntry(nil), entries...)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Score > result[i].Score || (result[j].Score == result[i].Score && result[j].RecordedOn < result[i].RecordedOn) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	for i := range result {
		result[i].Rank = i + 1
	}
	return result
}
