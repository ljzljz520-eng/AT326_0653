package store

import (
	"sort"
	"strings"

	"arcadeledger/internal/model"
)

func (d *DB) SaveScore(score model.HighScoreEntry) error {
	data, err := score.Marshal()
	if err != nil {
		return err
	}
	return d.write(scoresBucket, []byte(score.ID), data)
}

func (d *DB) GetScore(id string) (model.HighScoreEntry, error) {
	data, err := d.read(scoresBucket, []byte(id))
	if err != nil {
		return model.HighScoreEntry{}, err
	}
	return model.DecodeScore(data)
}

func (d *DB) ListScores(machineID string) ([]model.HighScoreEntry, error) {
	values, err := d.list(scoresBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.HighScoreEntry, 0, len(values))
	for _, data := range values {
		score, decodeErr := model.DecodeScore(data)
		if decodeErr != nil {
			return nil, decodeErr
		}
		if machineID == "" || score.MachineID == machineID {
			result = append(result, score)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return strings.Compare(result[i].RecordedOn, result[j].RecordedOn) < 0
		}
		return result[i].Score > result[j].Score
	})
	for i := range result {
		result[i].Rank = i + 1
	}
	return result, nil
}

func (d *DB) DeleteScore(id string) error { return d.remove(scoresBucket, []byte(id)) }

func (d *DB) ScoreCount(machineID string) (int, error) {
	scores, err := d.ListScores(machineID)
	return len(scores), err
}

func (d *DB) ScoreExists(machineID, playerID string, score int64) (bool, error) {
	scores, err := d.ListScores(machineID)
	if err != nil {
		return false, err
	}
	for _, entry := range scores {
		if entry.PlayerID == playerID && entry.Score == score {
			return true, nil
		}
	}
	return false, nil
}
