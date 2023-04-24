package storage

import "key-value-ts/domains/entities"

type Storer interface {
	Save(e entities.Sequence) error
	Get(key string, timestamp int64) (string, error)
}
