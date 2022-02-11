package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LevelDBOptions configuration for LevelDB.
type LevelDBOptions struct {
	DataDirectoryPath string `yaml:"DataDirectoryPath"`
}

// LevelDBStore is the official storage implementation for storing and retrieving
// blockchain data.
type LevelDBStore struct {
	db   *leveldb.DB
	path string
}

// NewLevelDBStore returns a new LevelDBStore object that will
// initialize the database found at the given path.
func NewLevelDBStore(cfg LevelDBOptions) (*LevelDBStore, error) {
	var opts = new(opt.Options) // should be exposed via LevelDBOptions if anything needed

	opts.Filter = filter.NewBloomFilter(10)
	db, err := leveldb.OpenFile(cfg.DataDirectoryPath, opts)
	if err != nil {
		return nil, err
	}

	return &LevelDBStore{
		path: cfg.DataDirectoryPath,
		db:   db,
	}, nil
}

// Put implements the Store interface.
func (s *LevelDBStore) Put(key, value []byte) error {
	return s.db.Put(key, value, nil)
}

// Get implements the Store interface.
func (s *LevelDBStore) Get(key []byte) ([]byte, error) {
	value, err := s.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		err = ErrKeyNotFound
	}
	return value, err
}

// Delete implements the Store interface.
func (s *LevelDBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

// PutBatch implements the Store interface.
func (s *LevelDBStore) PutBatch(batch Batch) error {
	lvldbBatch := batch.(*leveldb.Batch)
	return s.db.Write(lvldbBatch, nil)
}

// PutChangeSet implements the Store interface.
func (s *LevelDBStore) PutChangeSet(puts map[string][]byte) error {
	tx, err := s.db.OpenTransaction()
	if err != nil {
		return err
	}
	for k := range puts {
		if puts[k] != nil {
			err = tx.Put([]byte(k), puts[k], nil)
		} else {
			err = tx.Delete([]byte(k), nil)
		}
		if err != nil {
			tx.Discard()
			return err
		}
	}
	return tx.Commit()
}

// Seek implements the Store interface.
func (s *LevelDBStore) Seek(rng SeekRange, f func(k, v []byte) bool) {
	var (
		next func() bool
		ok   bool
		iter = s.db.NewIterator(seekRangeToPrefixes(rng), nil)
	)

	if !rng.Backwards {
		ok = iter.Next()
		next = iter.Next
	} else {
		ok = iter.Last()
		next = iter.Prev
	}

	for ; ok; ok = next() {
		if !f(iter.Key(), iter.Value()) {
			break
		}
	}
	iter.Release()
}

// Batch implements the Batch interface and returns a leveldb
// compatible Batch.
func (s *LevelDBStore) Batch() Batch {
	return new(leveldb.Batch)
}

// Close implements the Store interface.
func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
