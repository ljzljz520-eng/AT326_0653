package service

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

type Snapshot struct {
	Players  []model.PlayerProfile
	Machines []model.ArcadeMachine
	Scores   []model.HighScoreEntry
	Matches  []model.MatchSchedule
	Prizes   []model.PrizeDescription
}

func (s *ArcadeService) ExportSnapshot() (Snapshot, error) {
	players, err := s.db.ListPlayers()
	if err != nil {
		return Snapshot{}, err
	}
	machines, err := s.db.ListMachines()
	if err != nil {
		return Snapshot{}, err
	}
	scores, err := s.db.ListScores("")
	if err != nil {
		return Snapshot{}, err
	}
	matches, err := s.db.ListMatches()
	if err != nil {
		return Snapshot{}, err
	}
	prizes, err := s.db.ListPrizes(false)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Players: players, Machines: machines, Scores: scores, Matches: matches, Prizes: prizes}, nil
}

func (s *ArcadeService) ImportSnapshot(snapshot Snapshot) error {
	for _, player := range snapshot.Players {
		if err := model.ValidatePlayer(player); err != nil {
			return err
		}
		if err := s.db.SavePlayer(player); err != nil {
			return err
		}
	}
	for _, machine := range snapshot.Machines {
		if err := model.ValidateMachine(machine); err != nil {
			return err
		}
		if err := s.db.SaveMachine(machine); err != nil {
			return err
		}
	}
	for _, score := range snapshot.Scores {
		if err := model.ValidateScore(score); err != nil {
			return err
		}
		if err := s.db.SaveScore(score); err != nil {
			return err
		}
	}
	for _, match := range snapshot.Matches {
		if err := model.ValidateSchedule(match); err != nil {
			return err
		}
		if err := s.db.SaveMatch(match); err != nil {
			return err
		}
	}
	for _, prize := range snapshot.Prizes {
		if err := model.ValidatePrize(prize); err != nil {
			return err
		}
		if err := s.db.SavePrize(prize); err != nil {
			return err
		}
	}
	return nil
}

func BuildScoreID(machineID, playerID string, score int64, ordinal int) string {
	return fmt.Sprintf("%s-%s-%d-%d", strings.TrimSpace(machineID), strings.TrimSpace(playerID), score, ordinal)
}

func OpenService(path string) (*ArcadeService, error) {
	db, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	return New(db), nil
}
