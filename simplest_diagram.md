

classDiagram
    class Node {
        <<interface>>
    }

    class Trie {
    }

    class Leaf {
    }

    class Branch {
    }

    class Extension {
    }

    class Block {
    }

    class Blockchain {
    }

    class BlockchainIterator {
    }

    class ProofOfWork {
    }

    class boltStorage {
    }

    %% Relationships
    Trie o-- Node
    Leaf ..|> Node
    Branch ..|> Node
    Extension ..|> Node
    Branch o-- Node
    Extension o-- Node

    Blockchain o-- Trie
    Blockchain *-- boltStorage
    BlockchainIterator *-- boltStorage
    Blockchain ..> BlockchainIterator : Creates
    Blockchain ..> Block : Manages
    BlockchainIterator ..> Block : Decodes

    Block ..> boltStorage : Validates State
    ProofOfWork ..> Block : Mines



    
{
  "theme": "neutral",
  "class": {
    "hideEmptyMembersBox": true
  }
}