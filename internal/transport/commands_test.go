package transport

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
	"arcadeledger/internal/service"
	"arcadeledger/internal/store"
)

func TestCLICommandsAndRendering(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := service.New(db)
	command, err := Parse("register-player p1 Nova CN day-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Execute(svc, command); err != nil {
		t.Fatal(err)
	}
	if _, err := Execute(svc, Command{Name: "register-machine", Args: []string{"m1", "cab", "Asteroids", "Hall", "hard"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := Execute(svc, Command{Name: "score", Args: []string{"m1", "p1", "500", "day-1", "admin"}}); err != nil {
		t.Fatal(err)
	}
	board, err := svc.Leaderboard("m1", 10)
	if err != nil {
		t.Fatal(err)
	}
	if RenderBoard(board) == "no scores" || RenderPlayers([]model.PlayerProfile{{ID: "p1", Nickname: "Nova"}}) == "" {
		t.Fatal("rendering failed")
	}
}
