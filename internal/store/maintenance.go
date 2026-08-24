package store

import (
	"fmt"
	"sort"

	"arcadeledger/internal/model"
)

type StorageSummary struct {
	Players  int
	Machines int
	Scores   int
	Matches  int
	Prizes   int
	Receipts int
}

func (d *DB) Summary() (StorageSummary, error) {
	players, err := d.PlayerCount()
	if err != nil {
		return StorageSummary{}, err
	}
	machines, err := d.MachineCount()
	if err != nil {
		return StorageSummary{}, err
	}
	scores, err := d.ScoreCount("")
	if err != nil {
		return StorageSummary{}, err
	}
	matches, err := d.MatchCount()
	if err != nil {
		return StorageSummary{}, err
	}
	prizes, err := d.PrizeCount()
	if err != nil {
		return StorageSummary{}, err
	}
	receipts, err := d.ReceiptCount()
	if err != nil {
		return StorageSummary{}, err
	}
	return StorageSummary{Players: players, Machines: machines, Scores: scores, Matches: matches, Prizes: prizes, Receipts: receipts}, nil
}

func (d *DB) ValidateStorage() error {
	players, err := d.ListPlayers()
	if err != nil {
		return err
	}
	for _, player := range players {
		if err := model.ValidatePlayer(player); err != nil {
			return fmt.Errorf("player %s: %w", player.ID, err)
		}
	}
	machines, err := d.ListMachines()
	if err != nil {
		return err
	}
	for _, machine := range machines {
		if err := model.ValidateMachine(machine); err != nil {
			return fmt.Errorf("machine %s: %w", machine.ID, err)
		}
	}
	scores, err := d.ListScores("")
	if err != nil {
		return err
	}
	for _, score := range scores {
		if err := model.ValidateScore(score); err != nil {
			return fmt.Errorf("score %s: %w", score.ID, err)
		}
	}
	matches, err := d.ListMatches()
	if err != nil {
		return err
	}
	for _, match := range matches {
		if err := model.ValidateSchedule(match); err != nil {
			return fmt.Errorf("match %s: %w", match.ID, err)
		}
	}
	prizes, err := d.ListPrizes(false)
	if err != nil {
		return err
	}
	for _, prize := range prizes {
		if err := model.ValidatePrize(prize); err != nil {
			return fmt.Errorf("prize %s: %w", prize.ID, err)
		}
	}
	return nil
}

func (d *DB) ScoreHistory(playerID string) ([]model.HighScoreEntry, error) {
	scores, err := d.ListScores("")
	if err != nil {
		return nil, err
	}
	result := make([]model.HighScoreEntry, 0)
	for _, score := range scores {
		if score.PlayerID == playerID {
			result = append(result, score)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].RecordedOn == result[j].RecordedOn {
			return result[i].ID < result[j].ID
		}
		return result[i].RecordedOn < result[j].RecordedOn
	})
	return result, nil
}

func (d *DB) ReceiptForScore(scoreID string) (model.SubmissionReceipt, error) {
	receipts, err := d.ListReceipts()
	if err != nil {
		return model.SubmissionReceipt{}, err
	}
	for _, receipt := range receipts {
		if receipt.ScoreID == scoreID {
			return receipt, nil
		}
	}
	return model.SubmissionReceipt{}, ErrNotFound
}

func (d *DB) RemoveReceiptsForPlayer(playerID string) (int, error) {
	receipts, err := d.ListReceipts()
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, receipt := range receipts {
		if receipt.PlayerID == playerID {
			if deleteErr := d.remove(receiptsBucket, []byte(receipt.ID)); deleteErr != nil {
				return removed, deleteErr
			}
			removed++
		}
	}
	return removed, nil
}
