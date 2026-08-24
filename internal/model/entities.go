package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PlayerProfile struct {
	ID          string `json:"id"`
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display_name"`
	Country     string `json:"country"`
	CreatedAt   string `json:"created_at"`
	Active      bool   `json:"active"`
}

type ArcadeMachine struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Venue       string `json:"venue"`
	Difficulty  string `json:"difficulty"`
	Published   bool   `json:"published"`
	HighScoreID string `json:"high_score_id"`
}

type HighScoreEntry struct {
	ID         string `json:"id"`
	MachineID  string `json:"machine_id"`
	PlayerID   string `json:"player_id"`
	Nickname   string `json:"nickname"`
	Score      int64  `json:"score"`
	RecordedOn string `json:"recorded_on"`
	Rank       int    `json:"rank"`
	Verified   bool   `json:"verified"`
	Source     string `json:"source"`
}

type MatchSchedule struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	MachineIDs  []string `json:"machine_ids"`
	StartSlot   string   `json:"start_slot"`
	EndSlot     string   `json:"end_slot"`
	Status      string   `json:"status"`
	Capacity    int      `json:"capacity"`
	PlayerIDs   []string `json:"player_ids"`
	Description string   `json:"description"`
}

type PrizeDescription struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tier        string `json:"tier"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Claimed     int    `json:"claimed"`
	Active      bool   `json:"active"`
}

type SubmissionReceipt struct {
	ID        string `json:"id"`
	MachineID string `json:"machine_id"`
	PlayerID  string `json:"player_id"`
	Score     int64  `json:"score"`
	Accepted  bool   `json:"accepted"`
	Duplicate bool   `json:"duplicate"`
	Message   string `json:"message"`
	ScoreID   string `json:"score_id"`
}

func (p PlayerProfile) Marshal() ([]byte, error)     { return json.Marshal(p) }
func (m ArcadeMachine) Marshal() ([]byte, error)     { return json.Marshal(m) }
func (s HighScoreEntry) Marshal() ([]byte, error)    { return json.Marshal(s) }
func (m MatchSchedule) Marshal() ([]byte, error)     { return json.Marshal(m) }
func (p PrizeDescription) Marshal() ([]byte, error)  { return json.Marshal(p) }
func (r SubmissionReceipt) Marshal() ([]byte, error) { return json.Marshal(r) }

func DecodePlayer(data []byte) (PlayerProfile, error) {
	var v PlayerProfile
	err := json.Unmarshal(data, &v)
	return v, err
}
func DecodeMachine(data []byte) (ArcadeMachine, error) {
	var v ArcadeMachine
	err := json.Unmarshal(data, &v)
	return v, err
}
func DecodeScore(data []byte) (HighScoreEntry, error) {
	var v HighScoreEntry
	err := json.Unmarshal(data, &v)
	return v, err
}
func DecodeMatch(data []byte) (MatchSchedule, error) {
	var v MatchSchedule
	err := json.Unmarshal(data, &v)
	return v, err
}
func DecodePrize(data []byte) (PrizeDescription, error) {
	var v PrizeDescription
	err := json.Unmarshal(data, &v)
	return v, err
}
func DecodeReceipt(data []byte) (SubmissionReceipt, error) {
	var v SubmissionReceipt
	err := json.Unmarshal(data, &v)
	return v, err
}

func NormalizeNickname(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func PlayerKey(nickname string) string {
	key := NormalizeNickname(nickname)
	if key == "" {
		return ""
	}
	return key
}

func ScoreKey(machineID, playerID string, score int64) string {
	return fmt.Sprintf("%s|%s|%d", strings.TrimSpace(machineID), strings.TrimSpace(playerID), score)
}

func ScheduleKey(title string) string { return strings.ToLower(strings.TrimSpace(title)) }

func PrizeKey(name string) string { return strings.ToLower(strings.TrimSpace(name)) }
