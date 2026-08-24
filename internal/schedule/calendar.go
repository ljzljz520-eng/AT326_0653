package schedule

import (
	"fmt"
	"sort"
	"strings"

	"arcadeledger/internal/model"
)

type Calendar struct{ matches []model.MatchSchedule }

func NewCalendar(matches []model.MatchSchedule) *Calendar {
	return &Calendar{matches: append([]model.MatchSchedule(nil), matches...)}
}

func (c *Calendar) Add(match model.MatchSchedule) error {
	if err := model.ValidateSchedule(match); err != nil {
		return err
	}
	if c.Overlaps(match.StartSlot, match.EndSlot) {
		return fmt.Errorf("schedule overlaps an existing match")
	}
	c.matches = append(c.matches, match)
	return nil
}

func (c *Calendar) Remove(id string) bool {
	result := make([]model.MatchSchedule, 0, len(c.matches))
	removed := false
	for _, match := range c.matches {
		if match.ID == id {
			removed = true
			continue
		}
		result = append(result, match)
	}
	c.matches = result
	return removed
}

func (c *Calendar) Overlaps(start, end string) bool {
	for _, match := range c.matches {
		if start < match.EndSlot && end > match.StartSlot {
			return true
		}
	}
	return false
}

func (c *Calendar) ForMachine(machineID string) []model.MatchSchedule {
	result := make([]model.MatchSchedule, 0)
	for _, match := range c.matches {
		if ContainsMachine(match, machineID) {
			result = append(result, match)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StartSlot < result[j].StartSlot })
	return result
}

func (c *Calendar) OpenMatches() []model.MatchSchedule {
	result := make([]model.MatchSchedule, 0)
	for _, match := range c.matches {
		if strings.EqualFold(match.Status, "open") {
			result = append(result, match)
		}
	}
	return result
}

func (c *Calendar) NextForPlayer(playerID string) (model.MatchSchedule, bool) {
	choices := make([]model.MatchSchedule, 0)
	for _, match := range c.matches {
		if match.Status == "open" && !ContainsPlayer(match, playerID) {
			choices = append(choices, match)
		}
	}
	if len(choices) == 0 {
		return model.MatchSchedule{}, false
	}
	sort.Slice(choices, func(i, j int) bool { return choices[i].StartSlot < choices[j].StartSlot })
	return choices[0], true
}

func (c *Calendar) All() []model.MatchSchedule {
	return append([]model.MatchSchedule(nil), c.matches...)
}
