package model

import "testing"

func TestEntityValidationAndRanking(t *testing.T) {
	player := PlayerProfile{ID: "p1", Nickname: "  Nova  ", Active: true}
	if err := ValidatePlayer(player); err != nil {
		t.Fatal(err)
	}
	if PlayerKey(player.Nickname) != "nova" {
		t.Fatalf("unexpected key")
	}
	if err := ValidateScore(HighScoreEntry{ID: "s", MachineID: "m", PlayerID: "p", Score: 0, RecordedOn: "1"}); err != ErrInvalidScore {
		t.Fatalf("expected score error, got %v", err)
	}
	ranked := ScoreRank([]HighScoreEntry{{ID: "a", Score: 10, RecordedOn: "02"}, {ID: "b", Score: 10, RecordedOn: "01"}, {ID: "c", Score: 5, RecordedOn: "03"}})
	if ranked[0].ID != "b" || ranked[0].Rank != 1 || ranked[2].Rank != 3 {
		t.Fatalf("unexpected ranking: %#v", ranked)
	}
	if !CanTransition("draft", "open") || CanTransition("completed", "open") {
		t.Fatal("unexpected transition rules")
	}
}
