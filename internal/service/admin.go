package service

import (
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/query"
	"arcadeledger/internal/store"
)

type AdminOverview struct {
	Storage store.StorageSummary
	Stats   query.ScoreStats
	Top     []model.HighScoreEntry
	Healthy bool
}

func (s *ArcadeService) Overview(machineID string, limit int) (AdminOverview, error) {
	storage, err := s.db.Summary()
	if err != nil {
		return AdminOverview{}, err
	}
	stats, err := s.ScoreStats(machineID)
	if err != nil {
		return AdminOverview{}, err
	}
	top, err := s.Leaderboard(machineID, limit)
	if err != nil {
		return AdminOverview{}, err
	}
	healthy := s.db.ValidateStorage() == nil
	return AdminOverview{Storage: storage, Stats: stats, Top: top, Healthy: healthy}, nil
}

func (s *ArcadeService) DeactivatePlayer(id string) (model.PlayerProfile, error) {
	player, err := s.db.GetPlayer(id)
	if err != nil {
		return model.PlayerProfile{}, err
	}
	player.Active = false
	return player, s.db.SavePlayer(player)
}

func (s *ArcadeService) ArchiveMachine(id string) (model.ArcadeMachine, error) {
	machine, err := s.db.GetMachine(id)
	if err != nil {
		return model.ArcadeMachine{}, err
	}
	machine.Published = false
	return machine, s.db.SaveMachine(machine)
}

func (s *ArcadeService) ExportScoreLines(machineID string) ([]string, error) {
	scores, err := s.db.ListScores(machineID)
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0, len(scores))
	for _, score := range scores {
		lines = append(lines, model.EncodeScoreLine(score))
	}
	return lines, nil
}

func (s *ArcadeService) ImportScoreLines(lines []string) (int, error) {
	imported := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		score, err := model.ParseScoreLine(line)
		if err != nil {
			return imported, err
		}
		if err := s.db.SaveScore(score); err != nil {
			return imported, err
		}
		imported++
	}
	return imported, nil
}

func (s *ArcadeService) ScoreHistory(playerID string) ([]model.HighScoreEntry, error) {
	return s.db.ScoreHistory(playerID)
}

func (s *ArcadeService) ReceiptForScore(scoreID string) (model.SubmissionReceipt, error) {
	return s.db.ReceiptForScore(scoreID)
}

func (s *ArcadeService) RemovePlayerReceipts(playerID string) (int, error) {
	return s.db.RemoveReceiptsForPlayer(playerID)
}
