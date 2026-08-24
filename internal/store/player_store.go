package store

import (
	"sort"

	"arcadeledger/internal/model"
)

func (d *DB) SavePlayer(player model.PlayerProfile) error {
	data, err := player.Marshal()
	if err != nil {
		return err
	}
	return d.write(playersBucket, []byte(player.ID), data)
}

func (d *DB) GetPlayer(id string) (model.PlayerProfile, error) {
	data, err := d.read(playersBucket, []byte(id))
	if err != nil {
		return model.PlayerProfile{}, err
	}
	return model.DecodePlayer(data)
}

func (d *DB) FindPlayerByNickname(nickname string) (model.PlayerProfile, error) {
	players, err := d.ListPlayers()
	if err != nil {
		return model.PlayerProfile{}, err
	}
	key := model.PlayerKey(nickname)
	for _, player := range players {
		if model.PlayerKey(player.Nickname) == key {
			return player, nil
		}
	}
	return model.PlayerProfile{}, ErrNotFound
}

func (d *DB) ListPlayers() ([]model.PlayerProfile, error) {
	values, err := d.list(playersBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.PlayerProfile, 0, len(values))
	for _, data := range values {
		player, decodeErr := model.DecodePlayer(data)
		if decodeErr != nil {
			return nil, decodeErr
		}
		result = append(result, player)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (d *DB) DeletePlayer(id string) error { return d.remove(playersBucket, []byte(id)) }

func (d *DB) PlayerCount() (int, error) { return d.count(playersBucket) }
