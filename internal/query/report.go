package query

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

type Report struct {
	db    *store.DB
	board *Leaderboard
}

func NewReport(db *store.DB) *Report { return &Report{db: db, board: NewLeaderboard(db)} }

func (r *Report) Dashboard(machineID string, limit int) (Dashboard, error) {
	players, err := r.db.ListPlayers()
	if err != nil {
		return Dashboard{}, err
	}
	matches, err := r.db.ListMatches()
	if err != nil {
		return Dashboard{}, err
	}
	prizes, err := r.db.ListPrizes(true)
	if err != nil {
		return Dashboard{}, err
	}
	entries, err := r.board.Top(machineID, limit)
	if err != nil {
		return Dashboard{}, err
	}
	return Dashboard{MachineID: machineID, Players: len(players), Matches: len(matches), ActivePrizes: len(prizes), TopScores: entries}, nil
}

type Dashboard struct {
	MachineID    string
	Players      int
	Matches      int
	ActivePrizes int
	TopScores    []model.HighScoreEntry
}

func (r *Report) RenderText(d Dashboard) string {
	var b strings.Builder
	fmt.Fprintf(&b, "ARCADE LEDGER\nMachine: %s\nPlayers: %d\nMatches: %d\nActive prizes: %d\n", d.MachineID, d.Players, d.Matches, d.ActivePrizes)
	for _, score := range d.TopScores {
		fmt.Fprintf(&b, "%d. %s %d\n", score.Rank, score.Nickname, score.Score)
	}
	return b.String()
}

func (r *Report) PlayerCard(playerID string) (PlayerCard, error) {
	player, err := r.db.GetPlayer(playerID)
	if err != nil {
		return PlayerCard{}, err
	}
	scores, err := r.board.ByPlayer(playerID)
	if err != nil {
		return PlayerCard{}, err
	}
	return PlayerCard{Player: player, Scores: scores, TotalScore: total(scores)}, nil
}

type PlayerCard struct {
	Player     model.PlayerProfile
	Scores     []model.HighScoreEntry
	TotalScore int64
}

func total(scores []model.HighScoreEntry) int64 {
	var value int64
	for _, score := range scores {
		value += score.Score
	}
	return value
}
