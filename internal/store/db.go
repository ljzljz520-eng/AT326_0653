package store

import (
	"errors"
	"fmt"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

var (
	playersBucket  = []byte("players")
	machinesBucket = []byte("machines")
	scoresBucket   = []byte("scores")
	matchesBucket  = []byte("matches")
	prizesBucket   = []byte("prizes")
	receiptsBucket = []byte("receipts")
	ErrNotFound    = errors.New("record not found")
)

type DB struct {
	path string
	db   *bolt.DB
}

func Open(path string) (*DB, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	clean := filepath.Clean(path)
	handle, err := bolt.Open(clean, 0o600, nil)
	if err != nil {
		return nil, err
	}
	result := &DB{path: clean, db: handle}
	if err := result.initialize(); err != nil {
		_ = handle.Close()
		return nil, err
	}
	return result, nil
}

func (d *DB) initialize() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{playersBucket, machinesBucket, scoresBucket, matchesBucket, prizesBucket, receiptsBucket}
		for _, name := range buckets {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *DB) Close() error {
	if d == nil || d.db == nil {
		return nil
	}
	err := d.db.Close()
	d.db = nil
	return err
}

func (d *DB) Reopen() error {
	if d == nil {
		return fmt.Errorf("nil database")
	}
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			return err
		}
	}
	handle, err := bolt.Open(d.path, 0o600, nil)
	if err != nil {
		return err
	}
	d.db = handle
	return d.initialize()
}

func (d *DB) Path() string {
	if d == nil {
		return ""
	}
	return d.path
}

func (d *DB) read(bucket, key []byte) ([]byte, error) {
	if d == nil || d.db == nil {
		return nil, fmt.Errorf("database is closed")
	}
	var value []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrNotFound
		}
		data := b.Get(key)
		if data == nil {
			return ErrNotFound
		}
		value = append([]byte(nil), data...)
		return nil
	})
	return value, err
}

func (d *DB) write(bucket, key, value []byte) error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database is closed")
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrNotFound
		}
		return b.Put(key, value)
	})
}

func (d *DB) remove(bucket, key []byte) error {
	if d == nil || d.db == nil {
		return fmt.Errorf("database is closed")
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrNotFound
		}
		if b.Get(key) == nil {
			return ErrNotFound
		}
		return b.Delete(key)
	})
}

func (d *DB) list(bucket []byte) (map[string][]byte, error) {
	if d == nil || d.db == nil {
		return nil, fmt.Errorf("database is closed")
	}
	result := make(map[string][]byte)
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrNotFound
		}
		return b.ForEach(func(k, v []byte) error { result[string(k)] = append([]byte(nil), v...); return nil })
	})
	return result, err
}

func (d *DB) count(bucket []byte) (int, error) {
	values, err := d.list(bucket)
	if err != nil {
		return 0, err
	}
	return len(values), nil
}
