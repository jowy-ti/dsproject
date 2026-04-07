package blockchain

import "log"

type Blockchain struct {
	tip    []byte       // The hash of the last block in the chain
	boltDB *boltStorage // The connection to your database (BoltDB)
}

type BlockchainIterator struct {
	currentHash []byte
	boltDB      *boltStorage
}

func newBlockchain() *Blockchain {
	boltDB := newBoltStorage()

	tip, err := boltDB.dbGetLastHash()

	if err != nil {
		log.Fatalf("Initialising blockchain. %v. ", err)
	}

	if tip == nil {
		GenesisBlock := newGenesisBlock()
		err = boltDB.dbAddBlock(GenesisBlock.Hash, GenesisBlock.serialize())

		if err != nil {
			log.Fatal(err)
		}

		tip = GenesisBlock.Hash
	}

	blockchain := &Blockchain{
		tip:    tip,
		boltDB: boltDB,
	}

	return blockchain
}

func (bc *Blockchain) addBlock(block *Block) {
	err := bc.boltDB.dbAddBlock(block.Hash, block.serialize())

	if err != nil {
		log.Fatal(err)
	}

	bc.tip = block.Hash
}
