package pkg

type Blockchain struct {
	tip    [32]byte     // The hash of the last block in the chain
	boltDB *boltStorage // The connection to your database (BoltDB)
	trie   Trie         // Current Trie
}

// Iterator
type BlockchainIterator struct {
	currentHash [32]byte
	boltDB      *boltStorage
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

// newBlockchain initializes the blockchain and creates the genesis block if needed
func NewBlockchain(dbPath string) *Blockchain {
	boltDB := newBoltStorage(dbPath)
	tip := boltDB.dbGetLastHash()

	if tip == [32]byte{} {
		GenesisBlock := newGenesisBlock()
		boltDB.dbAddBlock(GenesisBlock.Hash, GenesisBlock.serialize())
		tip = GenesisBlock.getHash()
	}

	blockchain := &Blockchain{
		tip:    [32]byte(tip),
		boltDB: boltDB,
	}

	return blockchain
}

// InsertTrieValue inserts a new node in the trie
func (bc *Blockchain) InsertTrieValue(value string) {
	bc.trie.insert(value)
}

func (bc *Blockchain) VerifyValueInBlock(pos int, value string) bool {
	var found bool = false
	cont := 0

	bc.forEachBlock(func(b *Block) bool {
		if cont == pos {
			found = b.verifyValueInTrie(bc.boltDB, value)
			return true
		}
		cont++
		return false
	})

	return found
}

// AdBlock stores a new block and updates the chain tip
func (bc *Blockchain) AddBlock() {
	// PrintTrie(bc.trie.node)
	var nodes map[[32]byte][]byte = bc.trie.getKeysValues()
	// for key, _ := range nodes {
	// 	fmt.Printf("%x", key)
	// 	println()
	// }

	bc.storeNodes(nodes)
	prevBlockHash := bc.boltDB.dbGetLastHash()

	if prevBlockHash == [32]byte{} {
		panic("AddBlock. There are not blocks in the chain")
	}

	pow := newProofOfWork()
	block := newBlock([32]byte(prevBlockHash), bc.trie.node.Hash(), pow.difficulty)
	pow.mine(block)
	hash := block.getHash()
	bc.boltDB.dbAddBlock(hash, block.serialize())

	bc.tip = block.getHash()
	bc.trie.node = nil
}

// validateChain validates the correctness of hashes in the chain
func (bc *Blockchain) ValidateChain() bool {
	var validChain bool = false

	bc.forEachBlock(func(b *Block) bool {
		validHash := b.computeHash()

		if b.getHash() != validHash {
			return true
		}

		if b.getPrevBlockHash() == [32]byte{} {
			validChain = true
			return true
		}

		if !b.validateTrie(b.RootHash, bc.boltDB) {
			return true
		}

		return false
	})

	return validChain
}

// GetDataFromBlock gets the data from a Block regarding its position being the position 0 the last block
func (bc *Blockchain) GetDataFromBlock(pos int) map[[32]byte]string {
	var data map[[32]byte]string
	cont := 0

	bc.forEachBlock(func(b *Block) bool {
		if cont == pos {
			data = b.getTrieInfo(bc.boltDB)
			return true
		}

		cont++
		return false
	})
	return data
}

// storeNodes stores all key values from a map in db
func (bc *Blockchain) storeNodes(nodes map[[32]byte][]byte) {
	for key, value := range nodes {
		bc.boltDB.dbStoreTrieNode(key, value)
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
	it.currentHash = block.getPrevBlockHash()
	return block
}

// validHash checks if the iterator has a valid current position
func (it *BlockchainIterator) validHash() bool {
	return it.currentHash != [32]byte{}
}
