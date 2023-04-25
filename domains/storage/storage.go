package storage

import "key-value-ts/domains/entities"

// Storer interface that defines the basic logic for any storage implementation
type Storer interface {
	Save(e entities.Sequence) error
	Get(key string, timestamp int64) (string, error)
}
