package service

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/query"
	"arcadeledger/internal/store"
)

type ArcadeService struct {
	db        *store.DB
	board     *query.Leaderboard
	report    *query.Report
	analytics *query.Analytics
	sequence  int
}

func New(db *store.DB) *ArcadeService {
	return &ArcadeService{db: db, board: query.NewLeaderboard(db), report: query.NewReport(db), analytics: query.NewAnalytics(db)}
}

func (s *ArcadeService) RegisterPlayer(id, nickname, country, createdAt string) (model.PlayerProfile, error) {
	player := model.PlayerProfile{ID: strings.TrimSpace(id), Nickname: strings.TrimSpace(nickname), DisplayName: strings.TrimSpace(nickname), Country: strings.TrimSpace(country), CreatedAt: strings.TrimSpace(createdAt), Active: true}
	if err := model.ValidatePlayer(player); err != nil {
		return model.PlayerProfile{}, err
	}
	if existing, err := s.db.FindPlayerByNickname(player.Nickname); err == nil {
		return existing, fmt.Errorf("nickname already registered")
	}
	return player, s.db.SavePlayer(player)
}

func (s *ArcadeService) UpdatePlayer(id, nickname, country string, active bool) (model.PlayerProfile, error) {
	player, err := s.db.GetPlayer(id)
	if err != nil {
		return model.PlayerProfile{}, err
	}
	if strings.TrimSpace(nickname) != "" {
		player.Nickname = strings.TrimSpace(nickname)
		player.DisplayName = strings.TrimSpace(nickname)
	}
	if strings.TrimSpace(country) != "" {
		player.Country = strings.TrimSpace(country)
	}
	player.Active = active
	if err := model.ValidatePlayer(player); err != nil {
		return model.PlayerProfile{}, err
	}
	return player, s.db.SavePlayer(player)
}

func (s *ArcadeService) RegisterMachine(machine model.ArcadeMachine) (model.ArcadeMachine, error) {
	if err := model.ValidateMachine(machine); err != nil {
		return model.ArcadeMachine{}, err
	}
	if _, err := s.db.GetMachine(machine.ID); err == nil {
		return model.ArcadeMachine{}, fmt.Errorf("machine already exists")
	}
	return machine, s.db.SaveMachine(machine)
}

func (s *ArcadeService) PublishMachine(id string) (model.ArcadeMachine, error) {
	machine, err := s.db.GetMachine(id)
	if err != nil {
		return model.ArcadeMachine{}, err
	}
	machine.Published = true
	return machine, s.db.SaveMachine(machine)
}

func (s *ArcadeService) SubmitScore(machineID, playerID, nickname string, score int64, recordedOn, source string) (model.SubmissionReceipt, error) {
	machine, err := s.db.GetMachine(machineID)
	if err != nil {
		return model.SubmissionReceipt{}, err
	}
	if !machine.Published {
		return model.SubmissionReceipt{}, model.ErrInvalidMachine
	}
	player, err := s.db.GetPlayer(playerID)
	if err != nil {
		return model.SubmissionReceipt{}, err
	}
	if !player.Active {
		return model.SubmissionReceipt{}, fmt.Errorf("player is inactive")
	}
	entry := model.HighScoreEntry{MachineID: machineID, PlayerID: playerID, Nickname: nickname, Score: score, RecordedOn: recordedOn, Verified: true, Source: source}
	count, countErr := s.db.ScoreCount(machineID)
	if countErr != nil {
		return model.SubmissionReceipt{}, countErr
	}
	s.sequence++
	entry.ID = fmt.Sprintf("%s-%s-%d-%d", machineID, playerID, count, s.sequence)
	if err := model.ValidateScore(entry); err != nil {
		return model.SubmissionReceipt{}, err
	}
	if err := s.db.SaveScore(entry); err != nil {
		return model.SubmissionReceipt{}, err
	}
	receipt := model.SubmissionReceipt{ID: fmt.Sprintf("receipt-%s-%d", entry.ID, s.sequence), MachineID: machineID, PlayerID: playerID, Score: score, Accepted: true, ScoreID: entry.ID, Message: "score recorded"}
	if err := s.db.SaveReceipt(receipt); err != nil {
		return model.SubmissionReceipt{}, err
	}
	return receipt, nil
}

func (s *ArcadeService) Leaderboard(machineID string, limit int) ([]model.HighScoreEntry, error) {
	return s.board.Top(machineID, limit)
}

func (s *ArcadeService) Dashboard(machineID string, limit int) (query.Dashboard, error) {
	return s.report.Dashboard(machineID, limit)
}

func (s *ArcadeService) ScoreCount(machineID string) (int, error) { return s.db.ScoreCount(machineID) }

func (s *ArcadeService) ReceiptCount() (int, error) { return s.db.ReceiptCount() }

func (s *ArcadeService) ScoreStats(machineID string) (query.ScoreStats, error) {
	return s.analytics.Stats(machineID)
}

func (s *ArcadeService) PlayerBreakdown(playerID string) (query.PlayerScoreBreakdown, error) {
	return s.analytics.PlayerBreakdown(playerID)
}

func (s *ArcadeService) MachineStandings() ([]query.MachineStanding, error) {
	return s.analytics.Standings()
}

func (s *ArcadeService) Player(id string) (model.PlayerProfile, error) { return s.db.GetPlayer(id) }

func (s *ArcadeService) Machine(id string) (model.ArcadeMachine, error) { return s.db.GetMachine(id) }
