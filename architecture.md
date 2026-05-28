graph TD
    %% Styling
    classDef layerStyle fill:#f9f9f9,stroke:#333,stroke-width:2px;
    classDef componentStyle fill:#e1f5fe,stroke:#0288d1,stroke-width:1px;
    classDef dbStyle fill:#fff3e0,stroke:#f57c00,stroke-width:2px;

    subgraph Interface_Layer ["1. Interface and Diagnostic Layer"]
        A[Tests & Debugging Utilities] --> B[Blockchain Main Loop]
    end
    
    subgraph Core_Engine ["2. Core Logic and Consensus Layer"]
        B --> C[Blockchain Controller]
        C <--> D[ProofOfWork Mining Engine]
    end

    subgraph Cryptographic_Layer ["3. Authenticated Data Layer (MPT)"]
        C --> E[Trie State Manager]
        E --> F{Polymorphic Nodes}
        F --> G[Leaf Node]
        F --> H[Branch Node]
        F --> I[Extension Node]
    end

    subgraph Persistence_Layer ["4. Persistence Storage Layer"]
        J[boltStorage Wrapper]
        K[(B+ Tree BoltDB File)]
        
        subgraph Buckets
            L["'blocks' Bucket"]
            M["'tries' Bucket"]
        end
    end

    %% Cross-Layer Mapping
    C --> J
    G & H & I -- GOB Binary Stream --> J
    J --> K
    K --> L
    K --> M

    %% Class Applications
    style Interface_Layer layerStyle
    style Core_Engine layerStyle
    style Cryptographic_Layer layerStyle
    style Persistence_Layer layerStyle
    style K dbStyle