package pkg

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	bucketName  string = "main"
	lastHashKey string = "l"
)

type boltStorage struct {
	db *bolt.DB
}

// newBoltStorage initializes the storage and ensures the bucket exists
func newBoltStorage(dbpath string) *boltStorage {
	var err error

	boltDB := &boltStorage{
		db: dbConnection(dbpath),
	}

	if !boltDB.dbExistsBucket() {
		err = boltDB.db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte(bucketName))
			return err
		})
		if err != nil {
			log.Fatalf("Failed to create the bucket: %v", err)
			boltDB.db.Close()
		}
	}

	return boltDB

}

// dbGetLastHash returns the hash of the last block; nil if empty
func (boltDB *boltStorage) dbGetLastHash() [32]byte {
	var lastHash []byte

	if !boltDB.dbExistsBucket() {
		log.Fatalf("storage.go/dbGetLastHash. Can not get the last hash because does not exist the bucket with name '%s'", bucketName)
		boltDB.db.Close()
	}

	err := boltDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})

	if err != nil {
		log.Fatalf("storage.go/dbGetLastHash. %v", err)
		boltDB.db.Close()
	}

	if lastHash == nil {
		return [32]byte{}
	}

	return [32]byte(lastHash)
}

// dbAddBlock stores a block and updates the last hash pointer
func (boltDB *boltStorage) dbAddBlock(hash [32]byte, encodedBlock []byte) {
	err := boltDB.db.Update(func(tx *bolt.Tx) error {
		var err error

		bucket := tx.Bucket([]byte(bucketName))
		err = bucket.Put(hash[:], encodedBlock)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(lastHashKey), hash[:])
	})

	if err != nil {
		log.Fatalf("storage.go/dbAddBlock. %v", err)
		boltDB.db.Close()
	}
}

// dbGetEncodedBlock retrieves a block by its hash
func (boltDB *boltStorage) dbGetEncodedBlock(hash [32]byte) []byte {
	var encodedBlock []byte

	err := boltDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		encodedBlock = bucket.Get(hash[:])
		return nil
	})

	if err != nil {
		log.Fatalf("storage.go/dbGetEncodedBlock. %v", err)
		boltDB.db.Close()
	}

	return encodedBlock
}

// Private

// dbExistsBucket checks if the bucket exists
func (boltDB *boltStorage) dbExistsBucket() bool {
	var exists bool

	boltDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		exists = bucket != nil
		return nil
	})

	return exists
}

// dbConnection opens the database connection
func dbConnection(dbpath string) *bolt.DB {
	db, err := bolt.Open(dbpath, 0600, nil)

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		db.Close()
	}

	return db
}
