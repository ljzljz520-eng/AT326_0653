package query

import (
	"sort"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

type Leaderboard struct{ db *store.DB }

func NewLeaderboard(db *store.DB) *Leaderboard { return &Leaderboard{db: db} }

func (l *Leaderboard) Top(machineID string, limit int) ([]model.HighScoreEntry, error) {
	if limit < 1 {
		limit = 10
	}
	entries, err := l.db.ListScores(machineID)
	if err != nil {
		return nil, err
	}
	if len(entries) > limit {
		entries = entries[:limit]
	}
	return append([]model.HighScoreEntry(nil), entries...), nil
}

func (l *Leaderboard) ByPlayer(playerID string) ([]model.HighScoreEntry, error) {
	entries, err := l.db.ListScores("")
	if err != nil {
		return nil, err
	}
	result := entries[:0]
	for _, entry := range entries {
		if entry.PlayerID == playerID {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (l *Leaderboard) Search(machineID, nickname string) ([]model.HighScoreEntry, error) {
	entries, err := l.db.ListScores(machineID)
	if err != nil {
		return nil, err
	}
	needle := model.NormalizeNickname(nickname)
	result := make([]model.HighScoreEntry, 0)
	for _, entry := range entries {
		if needle == "" || strings.Contains(model.NormalizeNickname(entry.Nickname), needle) {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (l *Leaderboard) MachineSummary(machineID string) (MachineSummary, error) {
	entries, err := l.db.ListScores(machineID)
	if err != nil {
		return MachineSummary{}, err
	}
	summary := MachineSummary{MachineID: machineID, Entries: len(entries)}
	if len(entries) > 0 {
		summary.HighestScore = entries[0].Score
		summary.LeadingPlayer = entries[0].Nickname
	}
	unique := map[string]bool{}
	for _, entry := range entries {
		unique[entry.PlayerID] = true
	}
	summary.UniquePlayers = len(unique)
	return summary, nil
}

type MachineSummary struct {
	MachineID     string
	Entries       int
	UniquePlayers int
	HighestScore  int64
	LeadingPlayer string
}

func SortByMachine(entries []model.HighScoreEntry) []model.HighScoreEntry {
	result := append([]model.HighScoreEntry(nil), entries...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].MachineID == result[j].MachineID {
			return result[i].Score > result[j].Score
		}
		return result[i].MachineID < result[j].MachineID
	})
	return result
}
