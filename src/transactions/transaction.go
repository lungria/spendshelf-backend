package transactions

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"time"

	"go.etcd.io/bbolt"

	"github.com/lungria/spendshelf-backend/src/db"

	"go.uber.org/zap"
)

// Transaction represents a model of transactions in database
type Transaction struct {
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
	CategoryID  uint8     `json:"categoryId,omitempty"`
	Amount      int32     `json:"amount"`
	BankId      Bank      `json:"bankId"`
}

type Bank uint8

const (
	Mono Bank = iota + 1
	Privat
)

// Store implements by methods which define in Repository interface
type Store struct {
	logger *zap.SugaredLogger
	db     *db.Connection
}

var ErrTransactionNotFound = errors.New("transaction not found")

// NewStore creates a new instance of Repository
func NewStore(bolt *db.Connection, logger *zap.SugaredLogger) *Store {
	return &Store{
		logger: logger,
		db:     bolt,
	}
}

// ReadUncategorized returns all uncategorized transactions
func (repo *Store) ReadUncategorized() ([]Transaction, error) {
	var list []Transaction
	err := repo.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(db.UncategorizedTransactionsBucket))
		list = make([]Transaction, b.Stats().KeyN)
		buf := bytes.NewBuffer(make([]byte, 0))
		i := 0
		err := b.ForEach(func(k, v []byte) error {
			var t Transaction
			decoder := gob.NewDecoder(buf)
			err := decoder.Decode(&t)
			if err != nil {
				return err
			}
			list[i] = t
			i++
			buf.Reset()
			return nil
		})
		if err != nil {
			return err
		}
		return err
	})
	return list, err
}

// UpdateCategory changes the category for uncategorized transaction
func (repo *Store) SetCategory(transactionTime time.Time, categoryID uint8) error {
	return repo.db.Update(func(tx *bbolt.Tx) error {
		uncategorized := tx.Bucket([]byte(db.UncategorizedTransactionsBucket))
		key, err := transactionTime.MarshalBinary()
		if err != nil {
			return err
		}
		// pop from uncategorized bucket if exists
		data := uncategorized.Get(key)
		if data == nil {
			return ErrTransactionNotFound
		}
		err = uncategorized.Delete(key)
		if err != nil {
			return err
		}
		// set category
		var tr Transaction
		buf := bytes.NewBuffer(data)
		decoder := gob.NewDecoder(buf)
		err = decoder.Decode(&tr)
		if err != nil {
			return err
		}
		tr.CategoryID = categoryID
		// save to categorized bucket
		buf.Reset()
		encoder := gob.NewEncoder(buf)
		err = encoder.Encode(&tr)
		if err != nil {
			return err
		}
		readBuf, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}
		categorized := tx.Bucket([]byte(db.CategoriesBucket))
		err = categorized.Put(key, readBuf)
		return err
	})
}

// InsertMany inserts slice of transactions to database
func (repo *Store) InsertMany(txns []Transaction) error {
	return repo.db.Update(func(tx *bbolt.Tx) error {
		for _, t := range txns {
			b := tx.Bucket([]byte(db.UncategorizedTransactionsBucket))
			key, err := t.Time.MarshalBinary()
			if err != nil {
				return err
			}

			buf := bytes.NewBuffer(make([]byte, 0))
			encoder := gob.NewEncoder(buf)
			err = encoder.Encode(t)
			if err != nil {
				return err
			}
			readBuf, err := ioutil.ReadAll(buf)
			if err != nil {
				return err
			}
			err = b.Put(key, readBuf)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Insert transaction to database
func (repo *Store) Insert(t *Transaction) error {
	return repo.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(db.UncategorizedTransactionsBucket))
		key, err := t.Time.MarshalBinary()
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(make([]byte, 0))
		encoder := gob.NewEncoder(buf)
		err = encoder.Encode(*t)
		if err != nil {
			return err
		}
		readBuf, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}
		err = b.Put(key, readBuf)
		if err != nil {
			return err
		}

		return nil
	})
}

// Find transaction by time
func (repo *Store) Find(t time.Time) (*Transaction, error) {
	var tr Transaction
	err := repo.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(db.UncategorizedTransactionsBucket))
		key, err := t.MarshalBinary()
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(b.Get(key))
		decoder := gob.NewDecoder(buf)
		err = decoder.Decode(&tr)
		if err != nil {
			return err
		}
		return err
	})
	return &tr, err
}
