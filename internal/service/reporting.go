package service

import (
	"fmt"
	"sort"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/query"
	"arcadeledger/internal/schedule"
)

type CompetitionDigest struct {
	Match       model.MatchSchedule
	Scores      []model.HighScoreEntry
	Winner      model.HighScoreEntry
	PrizeTier   string
	Eligible    bool
	PlayerCount int
}

func (s *ArcadeService) CompetitionDigest(matchID string, machineID string) (CompetitionDigest, error) {
	match, err := s.db.GetMatch(matchID)
	if err != nil {
		return CompetitionDigest{}, err
	}
	scores, err := s.db.ListScores(machineID)
	if err != nil {
		return CompetitionDigest{}, err
	}
	ranked := model.ScoreRank(scores)
	digest := CompetitionDigest{Match: match, Scores: ranked, PlayerCount: len(match.PlayerIDs), Eligible: match.Status == "completed"}
	if len(ranked) > 0 {
		digest.Winner = ranked[0]
		digest.PrizeTier = tierForWinner(ranked[0].Rank, digest.Eligible)
	}
	return digest, nil
}

func tierForWinner(rank int, eligible bool) string {
	if !eligible {
		return ""
	}
	if rank <= 1 {
		return "champion"
	}
	if rank <= 3 {
		return "gold"
	}
	if rank <= 10 {
		return "silver"
	}
	return "bronze"
}

type AdminDigest struct {
	Overview AdminOverview
	Machines []query.MachineStanding
	Matches  []model.MatchSchedule
	Prizes   []model.PrizeDescription
	Lines    []string
}

func (s *ArcadeService) BuildAdminDigest(machineID string) (AdminDigest, error) {
	overview, err := s.Overview(machineID, 10)
	if err != nil {
		return AdminDigest{}, err
	}
	machines, err := s.MachineStandings()
	if err != nil {
		return AdminDigest{}, err
	}
	matches, err := s.db.ListMatches()
	if err != nil {
		return AdminDigest{}, err
	}
	prizes, err := s.db.ListPrizes(true)
	if err != nil {
		return AdminDigest{}, err
	}
	lines, err := s.ExportScoreLines(machineID)
	if err != nil {
		return AdminDigest{}, err
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].StartSlot < matches[j].StartSlot })
	return AdminDigest{Overview: overview, Machines: machines, Matches: matches, Prizes: prizes, Lines: lines}, nil
}

func (s *ArcadeService) ScheduleCalendar() (*schedule.Calendar, error) {
	return schedule.NewManager(s.db).Calendar()
}

func (s *ArcadeService) RenderAdminDigest(digest AdminDigest) string {
	var b strings.Builder
	fmt.Fprintf(&b, "players=%d machines=%d scores=%d\n", digest.Overview.Storage.Players, digest.Overview.Storage.Machines, digest.Overview.Storage.Scores)
	for _, machine := range digest.Machines {
		fmt.Fprintf(&b, "%s: %s (%d scores)\n", machine.MachineID, machine.Leader, machine.Scores)
	}
	for _, match := range digest.Matches {
		fmt.Fprintf(&b, "%s %s %s\n", match.ID, match.Title, match.Status)
	}
	for _, prize := range digest.Prizes {
		fmt.Fprintf(&b, "%s %s %d/%d\n", prize.ID, prize.Name, prize.Claimed, prize.Quantity)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func (s *ArcadeService) ValidateMachineReady(id string) error {
	machine, err := s.db.GetMachine(id)
	if err != nil {
		return err
	}
	if !machine.Published {
		return model.ErrInvalidMachine
	}
	if strings.TrimSpace(machine.Title) == "" || strings.TrimSpace(machine.Venue) == "" {
		return fmt.Errorf("machine metadata incomplete")
	}
	return nil
}
