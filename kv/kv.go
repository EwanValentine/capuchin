package kv

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"io"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// KV -
type KV struct {
	db *leveldb.DB
}

// Row -
type Row map[string]interface{}

// New -
func New() (*KV, error) {
	db, err := leveldb.OpenFile(".capuchin/db", nil)
	if err != nil {
		return nil, err
	}
	return &KV{db}, nil
}

// Close -
func (k *KV) Close() error {
	return k.db.Close()
}

// Put -
func (k *KV) Put(key, value []byte) error {
	return k.db.Put(key, value, nil)
}

// Get -
func (k *KV) Get(key string) ([]byte, error) {
	return k.db.Get([]byte(key), nil)
}

// Delete -
func (k *KV) Delete(key string) error {
	return k.db.Delete([]byte(key), nil)
}

// LoadCSV -
func (k *KV) LoadCSV(reader *csv.Reader) error {
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	if err == io.EOF {
		return nil
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				break
			}

			if err == io.EOF {
				break
			}

			return err
		}

		// Assumes the date is the first value, we could enforce this
		// to make life easier
		row := Row{}
		key := record[0]
		var encodedRow bytes.Buffer

		// Iterate through each value in the row,
		// grab the header, and use it as the key
		for key, val := range record {
			name := headers[key]
			row[name] = val
		}

		encoder := gob.NewEncoder(&encodedRow)
		if err := encoder.Encode(row); err != nil {
			return err
		}

		// Write to database
		if err := k.Put([]byte(key), encodedRow.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

// Query -
func (k *KV) Query(start, end int) ([]Row, error) {
	rows := []Row{}

	s := strconv.Itoa(start)
	e := strconv.Itoa(end)

	iter := k.db.NewIterator(&util.Range{Start: []byte(s), Limit: []byte(e)}, nil)
	for iter.Next() {
		var row Row
		decoder := gob.NewDecoder(bytes.NewBuffer(iter.Value()))
		if err := decoder.Decode(&row); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	iter.Release()

	if err := iter.Error(); err != nil {
		return rows, err
	}

	return rows, nil
}

// QueryPrefix - query prefix
func (k *KV) QueryPrefix(prefix int) ([]Row, error) {
	rows := []Row{}

	p := strconv.Itoa(prefix)
	iter := k.db.NewIterator(util.BytesPrefix([]byte(p)), nil)
	for iter.Next() {
		var row Row
		decoder := gob.NewDecoder(bytes.NewBuffer(iter.Value()))
		if err := decoder.Decode(&row); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	iter.Release()

	if err := iter.Error(); err != nil {
		return rows, err
	}

	return rows, nil
}
