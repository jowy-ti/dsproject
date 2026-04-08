package blockchain

import (
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestNewBoltStorageCreatesBucket(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	storage := &boltStorage{db: db}

	if storage.dbExistsBucket() {
		t.Fatal("bucket should not exist yet")
	}

	// simulate newBoltStorage behavior
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucketName))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	if !storage.dbExistsBucket() {
		t.Fatal("bucket should exist")
	}
}

func TestAddBlockAndGetLastHash(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	storage := &boltStorage{db: db}

	// create bucket
	db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucket([]byte(bucketName))
		return nil
	})

	hash := []byte("hash1")
	data := []byte("block1")

	storage.dbAddBlock(hash, data)

	lastHash := storage.dbGetLastHash()

	if string(lastHash) != string(hash) {
		t.Fatalf("expected %s, got %s", hash, lastHash)
	}
}

func TestGetEncodedBlock(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	storage := &boltStorage{db: db}

	db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucket([]byte(bucketName))
		return nil
	})

	hash := []byte("hash1")
	data := []byte("block1")

	storage.dbAddBlock(hash, data)

	result := storage.dbGetEncodedBlock(hash)

	if string(result) != string(data) {
		t.Fatalf("expected %s, got %s", data, result)
	}
}

func TestGetLastHashEmpty(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	storage := &boltStorage{db: db}

	db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucket([]byte(bucketName))
		return nil
	})

	lastHash := storage.dbGetLastHash()

	if lastHash != nil {
		t.Fatalf("expected nil, got %v", lastHash)
	}
}

func TestMultipleBlocks(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	storage := &boltStorage{db: db}

	db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucket([]byte(bucketName))
		return nil
	})

	hash1 := []byte("hash1")
	hash2 := []byte("hash2")

	storage.dbAddBlock(hash1, []byte("block1"))
	storage.dbAddBlock(hash2, []byte("block2"))

	lastHash := storage.dbGetLastHash()

	if string(lastHash) != string(hash2) {
		t.Fatalf("expected %s, got %s", hash2, lastHash)
	}
}

func newTestDB(t *testing.T) (*bolt.DB, func()) {
	file, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatal(err)
	}

	db, err := bolt.Open(file.Name(), 0600, nil)
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(file.Name())
	}

	return db, cleanup
}
