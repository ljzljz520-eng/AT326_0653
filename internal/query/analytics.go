package query

import (
	"fmt"
	"sort"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

type Analytics struct {
	db    *store.DB
	board *Leaderboard
}

func NewAnalytics(db *store.DB) *Analytics { return &Analytics{db: db, board: NewLeaderboard(db)} }

type ScoreStats struct {
	MachineID string
	Count     int
	Total     int64
	Average   int64
	Highest   int64
	Lowest    int64
	Verified  int
}

func (a *Analytics) Stats(machineID string) (ScoreStats, error) {
	entries, err := a.db.ListScores(machineID)
	if err != nil {
		return ScoreStats{}, err
	}
	stats := ScoreStats{MachineID: machineID, Count: len(entries)}
	if len(entries) == 0 {
		return stats, nil
	}
	stats.Lowest = entries[0].Score
	for _, entry := range entries {
		stats.Total += entry.Score
		if entry.Score > stats.Highest {
			stats.Highest = entry.Score
		}
		if entry.Score < stats.Lowest {
			stats.Lowest = entry.Score
		}
		if entry.Verified {
			stats.Verified++
		}
	}
	stats.Average = stats.Total / int64(stats.Count)
	return stats, nil
}

type PlayerScoreBreakdown struct {
	PlayerID      string
	Nickname      string
	Machines      int
	Entries       int
	Total         int64
	Best          int64
	MachineScores map[string]int64
}

func (a *Analytics) PlayerBreakdown(playerID string) (PlayerScoreBreakdown, error) {
	entries, err := a.board.ByPlayer(playerID)
	if err != nil {
		return PlayerScoreBreakdown{}, err
	}
	breakdown := PlayerScoreBreakdown{PlayerID: playerID, MachineScores: make(map[string]int64)}
	seen := make(map[string]bool)
	for _, entry := range entries {
		breakdown.Nickname = entry.Nickname
		breakdown.Entries++
		breakdown.Total += entry.Score
		if entry.Score > breakdown.Best {
			breakdown.Best = entry.Score
		}
		if !seen[entry.MachineID] {
			seen[entry.MachineID] = true
			breakdown.Machines++
		}
		if entry.Score > breakdown.MachineScores[entry.MachineID] {
			breakdown.MachineScores[entry.MachineID] = entry.Score
		}
	}
	return breakdown, nil
}

type MachineStanding struct {
	MachineID string
	Players   int
	Scores    int
	Highest   int64
	Leader    string
}

func (a *Analytics) Standings() ([]MachineStanding, error) {
	entries, err := a.db.ListScores("")
	if err != nil {
		return nil, err
	}
	byMachine := make(map[string][]model.HighScoreEntry)
	for _, entry := range entries {
		byMachine[entry.MachineID] = append(byMachine[entry.MachineID], entry)
	}
	result := make([]MachineStanding, 0, len(byMachine))
	for machineID, values := range byMachine {
		ranked := model.ScoreRank(values)
		players := make(map[string]bool)
		standing := MachineStanding{MachineID: machineID, Scores: len(values)}
		for _, value := range ranked {
			players[value.PlayerID] = true
		}
		standing.Players = len(players)
		if len(ranked) > 0 {
			standing.Highest = ranked[0].Score
			standing.Leader = ranked[0].Nickname
		}
		result = append(result, standing)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].MachineID < result[j].MachineID })
	return result, nil
}

func (a *Analytics) ExportText(machineID string) (string, error) {
	stats, err := a.Stats(machineID)
	if err != nil {
		return "", err
	}
	standings, err := a.Standings()
	if err != nil {
		return "", err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s,%d,%d,%d,%d\n", stats.MachineID, stats.Count, stats.Total, stats.Average, stats.Highest)
	for _, standing := range standings {
		fmt.Fprintf(&b, "%s,%d,%d,%s\n", standing.MachineID, standing.Players, standing.Highest, standing.Leader)
	}
	return b.String(), nil
}

func (a *Analytics) RankForPlayer(machineID, playerID string) (int, error) {
	entries, err := a.db.ListScores(machineID)
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if entry.PlayerID == playerID {
			return entry.Rank, nil
		}
	}
	return 0, store.ErrNotFound
}
