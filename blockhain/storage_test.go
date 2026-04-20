package blockchain

import (
	"bytes"
	global "dsproject/internal"
	"os"
	"testing"
)

// TestNewBoltStorage ensures the bucket is created correctly upon initialization
func TestNewBoltStorage(t *testing.T) {
	defer os.Remove(global.TestDBPath)

	storage := newBoltStorage(global.TestDBPath)
	defer storage.db.Close()

	if !storage.dbExistsBucket() {
		t.Errorf("Expected bucket '%s' to exist after initialization", bucketName)
	}
}

// TestAddAndGetBlock verifies that a block can be stored and retrieved by its hash
func TestAddAndGetBlock(t *testing.T) {
	storage, cleanup := newTestDB()
	defer cleanup()

	mockHash := []byte("block-hash-001")
	mockData := []byte("serialized-block-data")

	storage.dbAddBlock(mockHash, mockData)

	// Test retrieval
	retrieved := storage.dbGetEncodedBlock(mockHash)
	if !bytes.Equal(retrieved, mockData) {
		t.Errorf("Expected %x, got %x", mockData, retrieved)
	}

	// Test last hash update
	lastHash := storage.dbGetLastHash()
	if !bytes.Equal(lastHash, mockHash) {
		t.Errorf("Expected last hash to be %x, got %x", mockHash, lastHash)
	}
}

// TestLastHashUpdate verifies the 'l' key updates when a second block is added
func TestLastHashUpdate(t *testing.T) {
	storage, cleanup := newTestDB()
	defer cleanup()

	hash1 := []byte("hash1")
	hash2 := []byte("hash2")
	data := []byte("data")

	storage.dbAddBlock(hash1, data)
	storage.dbAddBlock(hash2, data)

	lastHash := storage.dbGetLastHash()
	if !bytes.Equal(lastHash, hash2) {
		t.Errorf("Expected last hash to update to %x, but it is %x", hash2, lastHash)
	}
}

// TestGetNonExistentBlock ensures the system handles missing keys gracefully
func TestGetNonExistentBlock(t *testing.T) {
	storage, cleanup := newTestDB()
	defer cleanup()

	result := storage.dbGetEncodedBlock([]byte("non-existent"))
	if result != nil {
		t.Errorf("Expected nil for non-existent block, got %v", result)
	}
}

// newTestDB used as a helper function to initialize, close and remove db
func newTestDB() (*boltStorage, func()) {
	storage := newBoltStorage(global.TestDBPath)

	cleanup := func() {
		storage.db.Close()
		os.Remove(global.TestDBPath)
	}

	return storage, cleanup
}
