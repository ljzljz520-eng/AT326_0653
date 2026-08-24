package service

import (
	"fmt"

	"arcadeledger/internal/catalog"
	"arcadeledger/internal/model"
	"arcadeledger/internal/schedule"
)

type WorkflowResult struct {
	Player    model.PlayerProfile
	Machine   model.ArcadeMachine
	Match     model.MatchSchedule
	Receipt   model.SubmissionReceipt
	Prize     model.PrizeDescription
	TopScores []model.HighScoreEntry
}

func (s *ArcadeService) RunScoreDay(input ScoreDayInput) (WorkflowResult, error) {
	player, err := s.RegisterPlayer(input.PlayerID, input.Nickname, input.Country, input.CreatedAt)
	if err != nil {
		return WorkflowResult{}, err
	}
	machine, err := s.RegisterMachine(input.Machine)
	if err != nil {
		return WorkflowResult{}, err
	}
	machine, err = s.PublishMachine(machine.ID)
	if err != nil {
		return WorkflowResult{}, err
	}
	manager := schedule.NewManager(s.db)
	match, err := manager.Create(input.MatchID, input.MatchTitle, []string{machine.ID}, input.StartSlot, input.EndSlot, input.Capacity, input.Description)
	if err != nil {
		return WorkflowResult{}, err
	}
	match, err = manager.OpenForRegistration(match.ID)
	if err != nil {
		return WorkflowResult{}, err
	}
	match, err = manager.RegisterPlayer(match.ID, player.ID)
	if err != nil {
		return WorkflowResult{}, err
	}
	receipt, err := s.SubmitScore(machine.ID, player.ID, player.Nickname, input.Score, input.RecordedOn, "score-day")
	if err != nil {
		return WorkflowResult{}, err
	}
	board, err := s.Leaderboard(machine.ID, input.Limit)
	if err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{Player: player, Machine: machine, Match: match, Receipt: receipt, TopScores: board}, nil
}

type ScoreDayInput struct {
	PlayerID    string
	Nickname    string
	Country     string
	CreatedAt   string
	Machine     model.ArcadeMachine
	MatchID     string
	MatchTitle  string
	StartSlot   string
	EndSlot     string
	Capacity    int
	Description string
	Score       int64
	RecordedOn  string
	Limit       int
}

func (s *ArcadeService) PublishPrize(id, name, tier, description string, quantity int) (model.PrizeDescription, error) {
	return catalog.New(s.db).Add(id, name, tier, description, quantity)
}

func (s *ArcadeService) ClaimPrizeForRank(rank int, tier string) (model.PrizeDescription, error) {
	return catalog.New(s.db).ClaimForRank(rank, tier)
}

func (s *ArcadeService) AdvanceMatch(id, target string) (model.MatchSchedule, error) {
	return schedule.NewManager(s.db).Transition(id, target)
}

func (s *ArcadeService) RegisterForMatch(id, playerID string) (model.MatchSchedule, error) {
	return schedule.NewManager(s.db).RegisterPlayer(id, playerID)
}

func (s *ArcadeService) ValidateWorkflow(result WorkflowResult) error {
	if result.Player.ID == "" || result.Machine.ID == "" || result.Match.ID == "" {
		return fmt.Errorf("workflow result missing identity")
	}
	if !result.Receipt.Accepted || result.Receipt.ScoreID == "" {
		return fmt.Errorf("workflow result missing receipt")
	}
	return nil
}
