package blockchain

import "bytes"

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
func newBlockchain(dbPath string) *Blockchain {
	boltDB := newBoltStorage(dbPath)
	tip := boltDB.dbGetLastHash()

	if tip == nil {
		GenesisBlock := newGenesisBlock()
		boltDB.dbAddBlock(GenesisBlock.Hash, GenesisBlock.serialize())
		tip = GenesisBlock.Hash
	}

	blockchain := &Blockchain{
		tip:    tip,
		boltDB: boltDB,
	}

	return blockchain
}

// addBlock stores a new block and updates the chain tip
func (bc *Blockchain) addBlock(block *Block) {
	bc.boltDB.dbAddBlock(block.Hash, block.serialize())
	bc.tip = block.Hash
}

// searchBlock looks for a block by its hash
func (bc *Blockchain) searchBlock(hash []byte) *Block {
	var block *Block

	bc.forEachBlock(func(b *Block) bool {
		if bytes.Equal(b.Hash, hash) {
			block = b
			return true
		}
		return false
	})

	return block
}

// forEachBlock iterates over all blocks and applies the given function
func (bc *Blockchain) forEachBlock(fn func(*Block) bool) {
	it := bc.newIterator()

	for it.validHash() {
		block := it.getCurrentBlock()
		if fn(block) {
			break
		}
		it.nextBlock(block)
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
func (it *BlockchainIterator) getCurrentBlock() *Block {
	encodedBlock := it.boltDB.dbGetEncodedBlock(it.currentHash)
	return deserializeBlock(encodedBlock)
}

// nextBlock moves the iterator to the previous block
func (it *BlockchainIterator) nextBlock(block *Block) {
	it.currentHash = block.PrevBlockHash
}

// validHash checks if the iterator has a valid current position
func (it *BlockchainIterator) validHash() bool {
	return len(it.currentHash) != 0
}

// func (it *BlockchainIterator) Reset() {
// 	it.currentHash = it.boltDB.dbGetLastHash()
// }
