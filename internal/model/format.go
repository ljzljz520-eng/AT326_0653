package model

import (
	"fmt"
	"strconv"
	"strings"
)

func FormatScore(score int64) string {
	if score < 0 {
		return "0"
	}
	value := strconv.FormatInt(score, 10)
	for index := len(value) - 3; index > 0; index -= 3 {
		value = value[:index] + "," + value[index:]
	}
	return value
}

func FormatPlayer(player PlayerProfile) string {
	status := "inactive"
	if player.Active {
		status = "active"
	}
	return fmt.Sprintf("%s (%s, %s)", player.DisplayName, player.Country, status)
}

func FormatMachine(machine ArcadeMachine) string {
	status := "hidden"
	if machine.Published {
		status = "published"
	}
	return fmt.Sprintf("%s - %s at %s [%s]", machine.Name, machine.Title, machine.Venue, status)
}

func FormatMatch(match MatchSchedule) string {
	return fmt.Sprintf("%s: %s-%s %s (%d/%d)", match.Title, match.StartSlot, match.EndSlot, match.Status, len(match.PlayerIDs), match.Capacity)
}

func FormatPrize(prize PrizeDescription) string {
	state := "inactive"
	if prize.Active {
		state = "active"
	}
	return fmt.Sprintf("%s [%s] %d/%d claimed (%s)", prize.Name, prize.Tier, prize.Claimed, prize.Quantity, state)
}

func ParseScoreLine(line string) (HighScoreEntry, error) {
	parts := strings.Split(strings.TrimSpace(line), "|")
	if len(parts) != 6 {
		return HighScoreEntry{}, fmt.Errorf("score line requires six fields")
	}
	score, err := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
	if err != nil {
		return HighScoreEntry{}, err
	}
	entry := HighScoreEntry{ID: strings.TrimSpace(parts[0]), MachineID: strings.TrimSpace(parts[1]), PlayerID: strings.TrimSpace(parts[2]), Score: score, RecordedOn: strings.TrimSpace(parts[4]), Source: strings.TrimSpace(parts[5]), Verified: true}
	if err := ValidateScore(entry); err != nil {
		return HighScoreEntry{}, err
	}
	return entry, nil
}

func EncodeScoreLine(entry HighScoreEntry) string {
	return strings.Join([]string{entry.ID, entry.MachineID, entry.PlayerID, strconv.FormatInt(entry.Score, 10), entry.RecordedOn, entry.Source}, "|")
}

func NormalizeCountry(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if len(value) > 3 {
		return value[:3]
	}
	return value
}

func NormalizeDescription(value string) string {
	words := strings.Fields(value)
	return strings.Join(words, " ")
}

func CopyStringSlice(values []string) []string { return append([]string(nil), values...) }

func ContainsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func RemoveString(values []string, target string) ([]string, bool) {
	result := make([]string, 0, len(values))
	removed := false
	for _, value := range values {
		if value == target {
			removed = true
			continue
		}
		result = append(result, value)
	}
	return result, removed
}
