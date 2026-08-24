package store

import (
	"sort"

	"arcadeledger/internal/model"
)

func (d *DB) SaveMatch(match model.MatchSchedule) error {
	data, err := match.Marshal()
	if err != nil {
		return err
	}
	return d.write(matchesBucket, []byte(match.ID), data)
}

func (d *DB) GetMatch(id string) (model.MatchSchedule, error) {
	data, err := d.read(matchesBucket, []byte(id))
	if err != nil {
		return model.MatchSchedule{}, err
	}
	return model.DecodeMatch(data)
}

func (d *DB) ListMatches() ([]model.MatchSchedule, error) {
	values, err := d.list(matchesBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.MatchSchedule, 0, len(values))
	for _, data := range values {
		match, decodeErr := model.DecodeMatch(data)
		if decodeErr != nil {
			return nil, decodeErr
		}
		result = append(result, match)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].StartSlot == result[j].StartSlot {
			return result[i].ID < result[j].ID
		}
		return result[i].StartSlot < result[j].StartSlot
	})
	return result, nil
}

func (d *DB) DeleteMatch(id string) error { return d.remove(matchesBucket, []byte(id)) }

func (d *DB) MatchCount() (int, error) { return d.count(matchesBucket) }
