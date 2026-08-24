package store

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arcade.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.SavePlayer(model.PlayerProfile{ID: "p1", Nickname: "Nova", Active: true}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveMachine(model.ArcadeMachine{ID: "m1", Name: "cab-1", Title: "Asteroids", Venue: "Hall", Difficulty: "hard", Published: true}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if err := db.Reopen(); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	player, err := db.GetPlayer("p1")
	if err != nil {
		t.Fatal(err)
	}
	if player.Nickname != "Nova" {
		t.Fatalf("unexpected player: %#v", player)
	}
	machine, err := db.GetMachine("m1")
	if err != nil {
		t.Fatal(err)
	}
	if !machine.Published {
		t.Fatal("machine did not persist")
	}
}

func TestDatabaseCountsAndRemoval(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.SavePrize(model.PrizeDescription{ID: "pr1", Name: "Token", Tier: "bronze", Description: "token", Quantity: 1, Active: true}); err != nil {
		t.Fatal(err)
	}
	if count, err := db.PrizeCount(); err != nil || count != 1 {
		t.Fatalf("unexpected count %d %v", count, err)
	}
	if err := db.DeletePrize("pr1"); err != nil {
		t.Fatal(err)
	}
	if err := db.DeletePrize("pr1"); err != ErrNotFound {
		t.Fatalf("unexpected delete result %v", err)
	}
}
