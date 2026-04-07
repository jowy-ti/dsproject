package blockchain

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	dbPath      string = "blockchain.db"
	bucketName  string = "main"
	lastHashKey string = "l"
)

type boltStorage struct {
	db *bolt.DB
}

func newBoltStorage() *boltStorage {
	var err error

	boltDB := &boltStorage{
		db: dbConnection(),
	}

	if !boltDB.dbExistsBucket() {
		err = boltDB.db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte(bucketName))
			return err
		})
		if err != nil {
			log.Fatalf("Failed to create the bucket: %v", err)
		}
	}

	return boltDB

}

// if hash is nil there is no blocks in the bucket
func (boltDB *boltStorage) dbGetLastHash() (hash []byte, e error) {
	var lastHash []byte

	if !boltDB.dbExistsBucket() {
		return nil, fmt.Errorf("The bucket does not exist")
	}

	err := boltDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})

	return lastHash, err
}

func (boltDB *boltStorage) dbAddBlock(hash []byte, serializedBlock []byte) error {
	return boltDB.db.Update(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket([]byte(bucketName))
		err = bucket.Put(hash, serializedBlock)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(lastHashKey), hash)
	})
}

// Private

func (boltDB *boltStorage) dbExistsBucket() bool {
	var exists bool

	boltDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		exists = bucket != nil
		return nil
	})

	return exists
}

func dbConnection() *bolt.DB {
	db, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db
}
