package schedule

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
)

func NormalizeStatus(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func ValidateSlots(start, end string) error {
	start, end = strings.TrimSpace(start), strings.TrimSpace(end)
	if start == "" || end == "" {
		return fmt.Errorf("both start and end slots are required")
	}
	if start >= end {
		return fmt.Errorf("start slot must precede end slot")
	}
	return nil
}

func ValidateCapacity(capacity, registered int) error {
	if capacity < 1 {
		return fmt.Errorf("capacity must be positive")
	}
	if registered < 0 {
		return fmt.Errorf("registered players cannot be negative")
	}
	if registered > capacity {
		return model.ErrCapacityReached
	}
	return nil
}

func ContainsPlayer(match model.MatchSchedule, playerID string) bool {
	return model.ContainsString(match.PlayerIDs, playerID)
}

func ContainsMachine(match model.MatchSchedule, machineID string) bool {
	for _, id := range match.MachineIDs {
		if id == machineID {
			return true
		}
	}
	return false
}
