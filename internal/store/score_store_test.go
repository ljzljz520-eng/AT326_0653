package store

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
)

func TestScoreStoreRanksAndFilters(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	values := []model.HighScoreEntry{{ID: "s1", MachineID: "m1", PlayerID: "p1", Nickname: "Nova", Score: 100, RecordedOn: "02"}, {ID: "s2", MachineID: "m1", PlayerID: "p2", Nickname: "Zed", Score: 200, RecordedOn: "01"}, {ID: "s3", MachineID: "m2", PlayerID: "p1", Nickname: "Nova", Score: 300, RecordedOn: "03"}}
	for _, value := range values {
		if err := db.SaveScore(value); err != nil {
			t.Fatal(err)
		}
	}
	scores, err := db.ListScores("m1")
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) != 2 || scores[0].Score != 200 || scores[0].Rank != 1 {
		t.Fatalf("unexpected scores %#v", scores)
	}
	exists, err := db.ScoreExists("m1", "p1", 100)
	if err != nil || !exists {
		t.Fatalf("score existence failed %v", err)
	}
}
