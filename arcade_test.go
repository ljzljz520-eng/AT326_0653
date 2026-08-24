package arcade_test

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
	"arcadeledger/internal/service"
	"arcadeledger/internal/store"
)

func TestArcadeScoreSubmissionIsIdempotent(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := service.New(db)
	if _, err := svc.RegisterPlayer("p1", "Nova", "CN", "day-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RegisterMachine(model.ArcadeMachine{ID: "m1", Name: "cab", Title: "Asteroids", Venue: "Hall", Difficulty: "hard"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.PublishMachine("m1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SubmitScore("m1", "p1", "Nova", 5000, "day-1", "admin"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SubmitScore("m1", "p1", "Nova", 5000, "day-1", "admin"); err != nil {
		t.Fatal(err)
	}
	count, err := svc.ScoreCount("m1")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected idempotent score count 1, got %d", count)
	}
}
