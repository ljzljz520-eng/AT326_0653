package schedule

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

type Manager struct{ db *store.DB }

func NewManager(db *store.DB) *Manager { return &Manager{db: db} }

func (m *Manager) Create(id, title string, machines []string, start, end string, capacity int, description string) (model.MatchSchedule, error) {
	if m == nil || m.db == nil {
		return model.MatchSchedule{}, fmt.Errorf("schedule manager is unavailable")
	}
	if strings.TrimSpace(id) == "" {
		return model.MatchSchedule{}, model.ErrMissingID
	}
	if err := ValidateSlots(start, end); err != nil {
		return model.MatchSchedule{}, err
	}
	if err := ValidateCapacity(capacity, 0); err != nil {
		return model.MatchSchedule{}, err
	}
	if len(machines) == 0 {
		return model.MatchSchedule{}, fmt.Errorf("machines cannot be empty")
	}
	match := model.MatchSchedule{ID: id, Title: strings.TrimSpace(title), MachineIDs: model.CopyStringSlice(machines), StartSlot: strings.TrimSpace(start), EndSlot: strings.TrimSpace(end), Status: "draft", Capacity: capacity, Description: model.NormalizeDescription(description)}
	if err := model.ValidateSchedule(match); err != nil {
		return model.MatchSchedule{}, err
	}
	if _, err := m.db.GetMatch(id); err == nil {
		return model.MatchSchedule{}, fmt.Errorf("schedule already exists")
	}
	return match, m.db.SaveMatch(match)
}

func (m *Manager) Get(id string) (model.MatchSchedule, error) { return m.db.GetMatch(id) }

func (m *Manager) List() ([]model.MatchSchedule, error) { return m.db.ListMatches() }

func (m *Manager) Calendar() (*Calendar, error) {
	matches, err := m.List()
	if err != nil {
		return nil, err
	}
	return NewCalendar(matches), nil
}

func (m *Manager) Transition(id, next string) (model.MatchSchedule, error) {
	match, err := m.Get(id)
	if err != nil {
		return model.MatchSchedule{}, err
	}
	next = NormalizeStatus(next)
	if !model.CanTransition(match.Status, next) {
		return model.MatchSchedule{}, fmt.Errorf("cannot transition %s to %s", match.Status, next)
	}
	match.Status = next
	return match, m.db.SaveMatch(match)
}

func (m *Manager) RegisterPlayer(id, playerID string) (model.MatchSchedule, error) {
	match, err := m.Get(id)
	if err != nil {
		return model.MatchSchedule{}, err
	}
	if match.Status != "open" {
		return model.MatchSchedule{}, fmt.Errorf("schedule is not open")
	}
	if ContainsPlayer(match, playerID) {
		return match, nil
	}
	if err := ValidateCapacity(match.Capacity, len(match.PlayerIDs)+1); err != nil {
		return model.MatchSchedule{}, err
	}
	match.PlayerIDs = append(match.PlayerIDs, playerID)
	return match, m.db.SaveMatch(match)
}

func (m *Manager) UnregisterPlayer(id, playerID string) (model.MatchSchedule, error) {
	match, err := m.Get(id)
	if err != nil {
		return model.MatchSchedule{}, err
	}
	filtered := match.PlayerIDs[:0]
	removed := false
	for _, value := range match.PlayerIDs {
		if value == playerID {
			removed = true
			continue
		}
		filtered = append(filtered, value)
	}
	if !removed {
		return match, nil
	}
	match.PlayerIDs = filtered
	return match, m.db.SaveMatch(match)
}

func (m *Manager) OpenForRegistration(id string) (model.MatchSchedule, error) {
	return m.Transition(id, "open")
}

func (m *Manager) Start(id string) (model.MatchSchedule, error) { return m.Transition(id, "running") }

func (m *Manager) Complete(id string) (model.MatchSchedule, error) {
	return m.Transition(id, "completed")
}

func (m *Manager) Cancel(id string) (model.MatchSchedule, error) {
	return m.Transition(id, "cancelled")
}
