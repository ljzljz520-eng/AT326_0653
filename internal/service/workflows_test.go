package service

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
	"arcadeledger/internal/schedule"
	"arcadeledger/internal/store"
)

func TestWorkflowScoreDay(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db)
	result, err := svc.RunScoreDay(ScoreDayInput{PlayerID: "p1", Nickname: "Nova", Country: "CN", CreatedAt: "day-1", Machine: model.ArcadeMachine{ID: "m1", Name: "cab-1", Title: "Asteroids", Venue: "Hall", Difficulty: "hard"}, MatchID: "match-1", MatchTitle: "Qualifier", StartSlot: "10:00", EndSlot: "11:00", Capacity: 2, Description: "qualifier", Score: 9000, RecordedOn: "day-1", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.ValidateWorkflow(result); err != nil {
		t.Fatal(err)
	}
	if len(result.TopScores) != 1 || result.TopScores[0].Score != 9000 {
		t.Fatalf("unexpected board %#v", result.TopScores)
	}
}

func TestWorkflowMatchRegistration(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db)
	if _, err := svc.RegisterPlayer("p1", "Nova", "CN", "day-1"); err != nil {
		t.Fatal(err)
	}
	m := &model.MatchSchedule{}
	managerResult, err := svc.AdvanceMatch("missing", "open")
	if err == nil || managerResult.ID != "" {
		t.Fatal("missing match unexpectedly advanced")
	}
	_ = m
	manager := newScheduleManager(db)
	match, err := manager.Create("match-1", "Open", []string{"m1"}, "09:00", "10:00", 1, "open")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := manager.OpenForRegistration(match.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RegisterForMatch(match.ID, "p1"); err != nil {
		t.Fatal(err)
	}
	if running, err := svc.AdvanceMatch(match.ID, "running"); err != nil || running.Status != "running" {
		t.Fatalf("unexpected state %#v %v", running, err)
	}
	if completed, err := svc.AdvanceMatch(match.ID, "completed"); err != nil || completed.Status != "completed" {
		t.Fatalf("unexpected state %#v %v", completed, err)
	}
}

func TestWorkflowPrizeClaim(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db)
	if _, err := svc.PublishPrize("pr1", "Champion Pin", "champion", "pin", 2); err != nil {
		t.Fatal(err)
	}
	prize, err := svc.ClaimPrizeForRank(1, "")
	if err != nil {
		t.Fatal(err)
	}
	if prize.Claimed != 1 {
		t.Fatalf("unexpected prize %#v", prize)
	}
}

func TestSnapshot(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := New(db)
	if _, err := svc.RegisterPlayer("p1", "Nova", "CN", "day-1"); err != nil {
		t.Fatal(err)
	}
	snapshot, err := svc.ExportSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Players) != 1 {
		t.Fatalf("unexpected snapshot %#v", snapshot)
	}
	if err := svc.ImportSnapshot(snapshot); err != nil {
		t.Fatal(err)
	}
}

func newScheduleManager(db *store.DB) *schedule.Manager { return schedule.NewManager(db) }
