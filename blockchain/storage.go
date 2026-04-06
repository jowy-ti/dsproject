package blockchain

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	dbPath     string = "blockchain.db"
	bucketName string = "main"
)

func dbConnection() *bolt.DB {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	return db
}

func (bc *Blockchain) dbAddBlock(key []byte, val []byte) error {
	return bc.db.Update(func(tx *bolt.Tx) error {
		// 1. Get the bucket by name
		bucket := tx.Bucket([]byte(bucketName))

		// 2. If it's the first time running, create it
		if bucket == nil {
			var err error
			bucket, err = tx.CreateBucket([]byte(bucketName))
			if err != nil {
				return err
			}
		}

		// 3. Use the bucket to store data
		return bucket.Put(key, val)
	})
}
