# dsproject

## Class Diagram

classDiagram
    class Node {
        <<interface>>
        +Hash() bytes32
        +Serialize() bytes
    }

    class Trie {
        -node Node
        +getKeysValues() map_hash_bytes
        +insert(value string)
    }

    class Leaf {
        +Path_unique bytes
        +Value bytes
        +Hash() bytes32
        +Serialize() bytes
    }

    class Branch {
        +Childs_hash array_of_bytes32
        -childs array_of_Node
        +insert(index int, node Node)
        +Hash() bytes32
        +Serialize() bytes
    }

    class Extension {
        +Path_shared bytes
        +Next_hash bytes32
        -next_branch Node
        +Hash() bytes32
        +Serialize() bytes
    }

    class Block {
        +Hash bytes32
        +PrevBlockHash bytes32
        +RootHash bytes32
        +Nonce uint64
        +Difficulty uint64
        +Timestamp int64
        -setHash(hash bytes32)
        -getHash() bytes32
        -getPrevBlockHash() bytes32
        -serialize() bytes
        -nextNonce()
        -computeHash() bytes32
        -validateTrie(hash bytes32, boltDB boltStorage) bool
        -getTrieInfo(boltDB boltStorage) map_hash_string
        -verifyValueInTrie(boltDB boltStorage, value string) bool
    }

    class Blockchain {
        -tip bytes32
        +BoltDB boltStorage
        -trie Trie
        -forEachBlock(fn function)
        +InsertTrieValue(value string)
        +VerifyValueInBlock(pos int, value string) bool
        +AddBlock()
        +ValidateChain() bool
        +GetDataFromBlock(pos int)
        -storeNodes(nodes map_hash_bytes)
        +GetChainLength() int
        -newIterator() BlockchainIterator
        +PrintBlockTrie(posBlock int)
    }

    class BlockchainIterator {
        -currentHash bytes32
        -boltDB boltStorage
        -getBlockAndAdvance() Block
        -validHash() bool
    }

    class ProofOfWork {
        -target BigInt
        -difficulty uint64
        -mine(block Block)
    }

    class boltStorage {
        -db BoltDBConn
        -dbExistsBuckets() bool
        -dbGetLastHash() bytes32
        -dbAddBlock(hash bytes32, encodedBlock bytes)
        -dbGetEncodedBlock(hash bytes32) bytes
        -dbStoreTrieNode(key bytes32, value bytes)
        -dbGetTrieValue(key bytes32) bytes
    }

    %% Relationships
    Trie o-- Node
    Leaf ..|> Node
    Branch ..|> Node
    Extension ..|> Node
    Branch o-- Node
    Extension o-- Node

    Blockchain o-- Trie
    Blockchain o-- boltStorage
    BlockchainIterator o-- boltStorage
    Blockchain ..> BlockchainIterator : Creates
    Blockchain ..> ProofOfWork : Creates
    Blockchain ..> Block : Manages
    BlockchainIterator ..> Block : Decodes

    Block ..> boltStorage : Validates State
    ProofOfWork ..> Block : Mines