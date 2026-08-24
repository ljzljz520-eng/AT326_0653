package store

import (
	"sort"

	"arcadeledger/internal/model"
)

func (d *DB) SavePrize(prize model.PrizeDescription) error {
	data, err := prize.Marshal()
	if err != nil {
		return err
	}
	return d.write(prizesBucket, []byte(prize.ID), data)
}

func (d *DB) GetPrize(id string) (model.PrizeDescription, error) {
	data, err := d.read(prizesBucket, []byte(id))
	if err != nil {
		return model.PrizeDescription{}, err
	}
	return model.DecodePrize(data)
}

func (d *DB) ListPrizes(activeOnly bool) ([]model.PrizeDescription, error) {
	values, err := d.list(prizesBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.PrizeDescription, 0, len(values))
	for _, data := range values {
		prize, decodeErr := model.DecodePrize(data)
		if decodeErr != nil {
			return nil, decodeErr
		}
		if !activeOnly || prize.Active {
			result = append(result, prize)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Tier == result[j].Tier {
			return result[i].Name < result[j].Name
		}
		return result[i].Tier < result[j].Tier
	})
	return result, nil
}

func (d *DB) DeletePrize(id string) error { return d.remove(prizesBucket, []byte(id)) }

func (d *DB) PrizeCount() (int, error) { return d.count(prizesBucket) }
