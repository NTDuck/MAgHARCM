# MAgHARCM Architecture & Workflow

`MAgHARCM` is a multi-agent repository-level code translation system implemented in Go using CloudWeGo Eino, based on the **ReCodeAgent** architecture ([arXiv:2604.07341](https://arxiv.org/abs/2604.07341)).

---

## 1. High-Level Agent Graph

```mermaid
flowchart TD
    START([Start / CLI Prompt]) --> Analyzer[Analyzer Agent<br/><i>gpt-oss:20b</i>]
    Analyzer --> Planning[Planning Agent<br/><i>gpt-oss:20b</i>]
    Planning --> Translator[Translator Agent<br/><i>Qwen3-4B-Instruct</i>]
    Translator --> Validator[Validator Agent<br/><i>gpt-oss:20b</i>]
    
    Validator --> Branch{Validation Passed?<br/><i>Pass Rate > 0%</i>}
    Branch -- No & Iteration < 10 --> Translator
    Branch -- Yes or Iteration >= 10 --> Complete([Complete / Target Repository])

    style START fill:#f9f9f9,stroke:#333,stroke-width:2px
    style Complete fill:#d4edda,stroke:#28a745,stroke-width:2px
    style Analyzer fill:#e1f5fe,stroke:#0288d1,stroke-width:2px
    style Planning fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    style Translator fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    style Validator fill:#ffebee,stroke:#d32f2f,stroke-width:2px
    style Branch fill:#fffde7,stroke:#fbc02d,stroke-width:2px
```

---

## 2. Detailed Dataflow & Tool Interaction

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant CLI as CLI / Main
    participant Graph as Eino Cyclic Graph
    participant Analyzer as Analyzer Agent
    participant Planning as Planning Agent
    participant Translator as Translator Agent
    participant Validator as Validator Agent
    participant Tools as Tool Layer (FS / LSP / Exec / Tree-Sitter)

    User->>CLI: Run MAgHARCM (--prompt-file prompts/c_to_rust.txt)
    CLI->>Graph: Invoke(initialState)

    rect rgb(225, 245, 254)
        Note over Graph,Analyzer: Phase 1: Codebase Analysis
        Graph->>Analyzer: Run(State)
        Analyzer->>Tools: get_directory_tree(source_dir)
        Analyzer->>Tools: parse_file_structure(ast_nodes)
        Analyzer-->>Graph: AnalyzerOutput (Research, Library Mapping, Target Design)
    end

    rect rgb(255, 243, 224)
        Note over Graph,Planning: Phase 2: Decomposition & Skeleton
        Graph->>Planning: Run(State)
        Planning->>Tools: write_file(Cargo.toml, src/lib.rs skeleton)
        Planning-->>Graph: PlanningOutput (NameMapping, Skeletons, Plan Part A/B)
    end

    loop Repair Loop (Up to 10 iterations)
        rect rgb(243, 229, 245)
            Note over Graph,Translator: Phase 3: Translation / Code Repair
            Graph->>Translator: Run(State)
            Translator->>Tools: write_file() / edit_file() on Target Files
            Translator-->>Graph: TranslatedProject (Source Code & Tests)
        end

        rect rgb(255, 235, 238)
            Note over Graph,Validator: Phase 4: Compilation & Test Validation
            Graph->>Validator: Run(State)
            Validator->>Tools: cargo check --tests (Compile verification)
            Validator->>Tools: cargo test --lib --tests (Test execution)
            Validator-->>Graph: ValidationReport (Passed/Failed counts, Diagnostics)
        end

        alt Validation Success (>0% pass rate) or Max Iterations Reached
            Note over Graph: Terminate Loop
        else Validation Failed
            Note over Graph: Increment Iteration & Loop back to Translator
        end
    end

    Graph-->>CLI: Final State
    CLI-->>User: Output Artifacts in Target Directory
```

---

## 3. Subsystem Breakdown

### 1. Analyzer Agent (`pkg/agents/analyzer.go`)
- **Model**: `gpt-oss:20b` (Reasoning)
- **Role**: Explores the source directory, parses AST structures (via Tree-sitter for C, C++, Go, Rust), and outputs:
  - Source Project Research
  - Third-Party Library Analysis
  - Target Project Architecture Design

### 2. Planning Agent (`pkg/agents/planning.go`)
- **Model**: `gpt-oss:20b` (Reasoning)
- **Role**: Bridges source AST to target semantics and initializes the target workspace:
  - Symbol Name Mapping (source $\rightarrow$ target)
  - Dynamic Target Skeleton Generation (`Cargo.toml`, module declarations)
  - 2-Part Implementation Plan:
    - **Part A**: Core source module translation
    - **Part B**: Test suite and characterization test translation

### 3. Translator Agent (`pkg/agents/translator.go`)
- **Model**: `Qwen3-4B-Instruct` (Coding)
- **Role**: Generates and iteratively refines target source code and test files:
  - **Initial Translation**: Generates idiomatic, safe target code adhering to the generated plan.
  - **Repair Mode**: Consumes validation diagnostics, compiler errors, and test failure logs to patch target files.

### 4. Validator Agent (`pkg/agents/validator.go`)
- **Model**: `gpt-oss:20b` (Reasoning)
- **Role**: Validates translation output via tool feedback:
  - Runs `cargo check --tests` and `cargo test --lib --tests`.
  - Parses compiler diagnostics and assertion failures.
  - Triggers Coverage-Guided Test Generation if functions lack coverage.
  - Evaluates completion: **compilation success with $>0\%$ test pass rate**.
