package store

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
)

func TestPlayerStoreFindAndList(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, p := range []model.PlayerProfile{{ID: "p2", Nickname: "Zed", Active: true}, {ID: "p1", Nickname: "Nova", Active: true}} {
		if err := db.SavePlayer(p); err != nil {
			t.Fatal(err)
		}
	}
	found, err := db.FindPlayerByNickname(" nova ")
	if err != nil || found.ID != "p1" {
		t.Fatalf("unexpected find %#v %v", found, err)
	}
	players, err := db.ListPlayers()
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 2 || players[0].ID != "p1" {
		t.Fatalf("unexpected players %#v", players)
	}
}
