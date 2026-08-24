package transport

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
)

func RenderBoard(entries []model.HighScoreEntry) string {
	if len(entries) == 0 {
		return "no scores"
	}
	var b strings.Builder
	for _, entry := range entries {
		fmt.Fprintf(&b, "%d. %s %s (%s)\n", entry.Rank, entry.Nickname, model.FormatScore(entry.Score), entry.MachineID)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func RenderPlayers(players []model.PlayerProfile) string {
	var b strings.Builder
	for _, player := range players {
		fmt.Fprintf(&b, "%s: %s [%s]\n", player.ID, player.Nickname, player.Country)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func RenderMatches(matches []model.MatchSchedule) string {
	var b strings.Builder
	for _, match := range matches {
		fmt.Fprintf(&b, "%s %s %s-%s (%s)\n", match.ID, match.Title, match.StartSlot, match.EndSlot, match.Status)
	}
	return strings.TrimSuffix(b.String(), "\n")
}
