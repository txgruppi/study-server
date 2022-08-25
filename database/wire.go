package database

import (
	"github.com/timshannon/badgerhold"
)

func Wire() (*badgerhold.Store, error) {
	options := badgerhold.DefaultOptions
	options.Dir = "./data"
	options.ValueDir = "./data"
	store, err := badgerhold.Open(options)
	if err != nil {
		return nil, err
	}
	return store, nil
}
