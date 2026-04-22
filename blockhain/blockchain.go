package blockchain

import (
	"bytes"
)

type Blockchain struct {
	tip    []byte       // The hash of the last block in the chain
	boltDB *boltStorage // The connection to your database (BoltDB)
}

// Iterator
type BlockchainIterator struct {
	currentHash []byte
	boltDB      *boltStorage
}

// newBlockchain initializes the blockchain and creates the genesis block if needed
func NewBlockchain(dbPath string) *Blockchain {
	boltDB := newBoltStorage(dbPath)
	tip := boltDB.dbGetLastHash()

	if tip == nil {
		GenesisBlock := newGenesisBlock()
		boltDB.dbAddBlock(GenesisBlock.hash, GenesisBlock.serialize())
		tip = GenesisBlock.hash
	}

	blockchain := &Blockchain{
		tip:    tip,
		boltDB: boltDB,
	}

	return blockchain
}

// addBlock stores a new block and updates the chain tip
func (bc *Blockchain) AddBlock(data string) {
	prevBlockHash := bc.boltDB.dbGetLastHash()
	pow := newProofOfWork()
	block := newBlock(data, prevBlockHash, pow.difficulty)

	hash, nonce := pow.mine()
	block.setInfoAfterMining(hash, nonce)

	bc.boltDB.dbAddBlock(block.hash, block.serialize())

	bc.tip = block.hash
}

// searchBlock looks for a block by its hash
func (bc *Blockchain) SearchBlock(hash []byte) *Block {
	var block *Block

	bc.forEachBlock(func(b *Block) bool {
		if bytes.Equal(b.hash, hash) {
			block = b
			return true
		}
		return false
	})

	return block
}

// validateChain validates the correctness of hashes in the chain
func (bc *Blockchain) ValidateChain() bool {
	var validChain bool = false

	bc.forEachBlock(func(b *Block) bool {
		if !bytes.Equal(b.hash, b.computeHash()) {
			return true
		}

		if len(b.prevBlockHash) == 0 {
			validChain = true
		}

		return false
	})

	return validChain
}

// forEachBlock iterates over all blocks and applies the given function
func (bc *Blockchain) forEachBlock(fn func(*Block) bool) {
	it := bc.newIterator()

	for it.validHash() {
		block := it.getBlockAndAdvance()
		if fn(block) {
			break
		}
	}
}

// newIterator creates an iterator starting from the tip
func (bc *Blockchain) newIterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.tip,
		boltDB:      bc.boltDB,
	}
}

// getCurrentBlock returns the block at the current iterator position
func (it *BlockchainIterator) getBlockAndAdvance() *Block {
	encodedBlock := it.boltDB.dbGetEncodedBlock(it.currentHash)
	block := deserializeBlock(encodedBlock)
	it.currentHash = block.prevBlockHash
	return block
}

// validHash checks if the iterator has a valid current position
func (it *BlockchainIterator) validHash() bool {
	return len(it.currentHash) != 0
}

// func (it *BlockchainIterator) Reset() {
// 	it.currentHash = it.boltDB.dbGetLastHash()
// }
