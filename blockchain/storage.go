package blockchain

import (
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
func (boltDB *boltStorage) dbGetLastHash() []byte {
	var lastHash []byte

	if !boltDB.dbExistsBucket() {
		log.Fatalf("storage.go/dbGetLastHash. Can not get the last hash because does not exist the bucket with name '%s'", bucketName)
	}

	err := boltDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})

	if err != nil {
		log.Fatalf("storage.go/dbGetLastHash. %v", err)
	}

	return lastHash
}

func (boltDB *boltStorage) dbAddBlock(hash []byte, encodedBlock []byte) {
	err := boltDB.db.Update(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket([]byte(bucketName))
		err = bucket.Put(hash, encodedBlock)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(lastHashKey), hash)
	})

	if err != nil {
		log.Fatalf("storage.go/dbAddBlock. %v", err)
	}
}

func (boltDB *boltStorage) dbGetEncodedBlock(hash []byte) []byte {
	var encodedBlock []byte

	err := boltDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		encodedBlock = bucket.Get(hash)
		return nil
	})

	if err != nil {
		log.Fatalf("storage.go/dbGetEncodedBlock. %v", err)
	}

	return encodedBlock
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
