package store

import (
	"sort"

	"arcadeledger/internal/model"
)

func (d *DB) SaveReceipt(receipt model.SubmissionReceipt) error {
	data, err := receipt.Marshal()
	if err != nil {
		return err
	}
	return d.write(receiptsBucket, []byte(receipt.ID), data)
}

func (d *DB) GetReceipt(id string) (model.SubmissionReceipt, error) {
	data, err := d.read(receiptsBucket, []byte(id))
	if err != nil {
		return model.SubmissionReceipt{}, err
	}
	return model.DecodeReceipt(data)
}

func (d *DB) ListReceipts() ([]model.SubmissionReceipt, error) {
	values, err := d.list(receiptsBucket)
	if err != nil {
		return nil, err
	}
	result := make([]model.SubmissionReceipt, 0, len(values))
	for _, data := range values {
		receipt, decodeErr := model.DecodeReceipt(data)
		if decodeErr != nil {
			return nil, decodeErr
		}
		result = append(result, receipt)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (d *DB) ReceiptCount() (int, error) { return d.count(receiptsBucket) }
