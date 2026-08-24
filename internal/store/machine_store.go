package store

import (
	"sort"

	"arcadeledger/internal/model"
)

func (d *DB) SaveMachine(machine model.ArcadeMachine) error {
	data, err := machine.Marshal()
	if err != nil {
		return err
	}
	return d.write(machinesBucket, []byte(machine.ID), data)
}

func (d *DB) GetMachine(id string) (model.ArcadeMachine, error) {
	data, err := d.read(machinesBucket, []byte(id))
	if err != nil {
		return model.ArcadeMachine{}, err
	}
	return model.DecodeMachine(data)
}

func (d *DB) ListMachines(publishedOnly ...bool) ([]model.ArcadeMachine, error) {
	values, err := d.list(machinesBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.ArcadeMachine, 0, len(values))
	onlyPublished := len(publishedOnly) > 0 && publishedOnly[0]
	for _, data := range values {
		machine, decodeErr := model.DecodeMachine(data)
		if decodeErr != nil {
			return nil, decodeErr
		}
		if !onlyPublished || machine.Published {
			result = append(result, machine)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (d *DB) DeleteMachine(id string) error { return d.remove(machinesBucket, []byte(id)) }

func (d *DB) MachineCount() (int, error) { return d.count(machinesBucket) }
