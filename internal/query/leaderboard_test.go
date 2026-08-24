package query

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

func TestLeaderboardSearchAndSummary(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, score := range []model.HighScoreEntry{{ID: "a", MachineID: "m1", PlayerID: "p1", Nickname: "Nova", Score: 90, RecordedOn: "01"}, {ID: "b", MachineID: "m1", PlayerID: "p2", Nickname: "Zed", Score: 80, RecordedOn: "02"}} {
		if err := db.SaveScore(score); err != nil {
			t.Fatal(err)
		}
	}
	l := NewLeaderboard(db)
	entries, err := l.Search("m1", "nov")
	if err != nil || len(entries) != 1 {
		t.Fatalf("unexpected search %#v %v", entries, err)
	}
	summary, err := l.MachineSummary("m1")
	if err != nil || summary.HighestScore != 90 || summary.UniquePlayers != 2 {
		t.Fatalf("unexpected summary %#v %v", summary, err)
	}
}
