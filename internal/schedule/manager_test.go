package schedule

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

func TestScheduleManagerLifecycle(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	m := NewManager(db)
	match, err := m.Create("match-1", "Friday Frenzy", []string{"m1"}, "10:00", "11:00", 1, "qualifier")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.OpenForRegistration(match.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.RegisterPlayer(match.ID, "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.RegisterPlayer(match.ID, "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.RegisterPlayer(match.ID, "p2"); err != model.ErrCapacityReached {
		t.Fatalf("expected capacity error, got %v", err)
	}
	if _, err := m.Start(match.ID); err != nil {
		t.Fatal(err)
	}
	if completed, err := m.Complete(match.ID); err != nil || completed.Status != "completed" {
		t.Fatalf("unexpected completion %#v %v", completed, err)
	}
}
