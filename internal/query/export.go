package query

import (
	"fmt"
	"sort"
	"strings"

	"arcadeledger/internal/model"
)

type LeaderboardPage struct {
	MachineID string
	Page      int
	PageSize  int
	Total     int
	Entries   []model.HighScoreEntry
	HasNext   bool
}

func Page(entries []model.HighScoreEntry, machineID string, page, pageSize int) LeaderboardPage {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	if start > len(entries) {
		start = len(entries)
	}
	end := start + pageSize
	if end > len(entries) {
		end = len(entries)
	}
	result := LeaderboardPage{MachineID: machineID, Page: page, PageSize: pageSize, Total: len(entries), Entries: append([]model.HighScoreEntry(nil), entries[start:end]...)}
	result.HasNext = end < len(entries)
	return result
}

func (l *Leaderboard) Page(machineID string, page, pageSize int) (LeaderboardPage, error) {
	entries, err := l.db.ListScores(machineID)
	if err != nil {
		return LeaderboardPage{}, err
	}
	return Page(entries, machineID, page, pageSize), nil
}

func FilterVerified(entries []model.HighScoreEntry, verified bool) []model.HighScoreEntry {
	result := make([]model.HighScoreEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.Verified == verified {
			result = append(result, entry)
		}
	}
	return result
}

func FilterSource(entries []model.HighScoreEntry, source string) []model.HighScoreEntry {
	needle := strings.ToLower(strings.TrimSpace(source))
	if needle == "" {
		return append([]model.HighScoreEntry(nil), entries...)
	}
	result := make([]model.HighScoreEntry, 0, len(entries))
	for _, entry := range entries {
		if strings.ToLower(strings.TrimSpace(entry.Source)) == needle {
			result = append(result, entry)
		}
	}
	return result
}

func DistinctPlayers(entries []model.HighScoreEntry) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, entry := range entries {
		if !seen[entry.PlayerID] {
			seen[entry.PlayerID] = true
			result = append(result, entry.PlayerID)
		}
	}
	sort.Strings(result)
	return result
}

func GroupByPlayer(entries []model.HighScoreEntry) map[string][]model.HighScoreEntry {
	result := make(map[string][]model.HighScoreEntry)
	for _, entry := range entries {
		result[entry.PlayerID] = append(result[entry.PlayerID], entry)
	}
	for playerID := range result {
		sort.Slice(result[playerID], func(i, j int) bool { return result[playerID][i].Score > result[playerID][j].Score })
	}
	return result
}

func RenderTable(entries []model.HighScoreEntry) string {
	if len(entries) == 0 {
		return "| rank | player | score |\n| --- | --- | --- |"
	}
	var b strings.Builder
	b.WriteString("| rank | player | score |\n| --- | --- | --- |\n")
	for _, entry := range entries {
		fmt.Fprintf(&b, "| %d | %s | %s |\n", entry.Rank, entry.Nickname, model.FormatScore(entry.Score))
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func Compare(entries []model.HighScoreEntry, left, right string) (int64, int64, int64) {
	var leftBest, rightBest int64
	for _, entry := range entries {
		if entry.PlayerID == left && entry.Score > leftBest {
			leftBest = entry.Score
		}
		if entry.PlayerID == right && entry.Score > rightBest {
			rightBest = entry.Score
		}
	}
	return leftBest, rightBest, leftBest - rightBest
}

func ScoreBands(entries []model.HighScoreEntry, width int64) map[int]int {
	result := make(map[int]int)
	if width < 1 {
		width = 1000
	}
	for _, entry := range entries {
		result[int(entry.Score/width)]++
	}
	return result
}
