package blockchain

import (
	bolt "go.etcd.io/bbolt"
)

type Blockchain struct {
	tip []byte   // The hash of the last block in the chain
	db  *bolt.DB // The connection to your database (BoltDB)
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func newBlockchain() *Blockchain {
	firstBlock := newGenesisBlock()

	blockchain := &Blockchain{
		tip: firstBlock.Hash,
		db:  dbConnection(),
	}

	return blockchain
}

func (bc *Blockchain) addBlock() {

}
