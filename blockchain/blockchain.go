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

func newBlockchain() *Blockchain {
	boltDB := newBoltStorage()
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

func (bc *Blockchain) addBlock(block *Block) {
	bc.boltDB.dbAddBlock(block.Hash, block.serialize())
	bc.tip = block.Hash
}

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

func (bc *Blockchain) newIterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.tip,
		boltDB:      bc.boltDB,
	}
}

func (it *BlockchainIterator) getCurrentBlock() *Block {
	encodedBlock := it.boltDB.dbGetEncodedBlock(it.currentHash)
	return deserializeBlock(encodedBlock)
}

func (it *BlockchainIterator) nextBlock(block *Block) {
	it.currentHash = block.PrevBlockHash
}

func (it *BlockchainIterator) validHash() bool {
	return len(it.currentHash) != 0
}

// func (it *BlockchainIterator) Reset() {
// 	it.currentHash = it.boltDB.dbGetLastHash()
// }
