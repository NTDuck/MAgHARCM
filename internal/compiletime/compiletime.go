package compiletime

import (
	"errors"
	"fmt"
)

// Backlink: [[1.0.0 Architecture]] and [[1.0.0 Methodology]].
// Package compiletime provides centralized compile-time configuration,
// constants, sentinels, enums, and initialization helpers using the Must pattern.
// These are strictly distinguished from runtime YAML configurations (package config).

// -------------------------------------------------------------------------
// Compile-time Sentinels & Constants
// -------------------------------------------------------------------------

// CurrentSchemaVersion is the canonical schema version string for role artifacts.
const CurrentSchemaVersion = "1.0.0"
// ErrorUnknown is the sentinel compiler/toolchain string when no language matches.
const ErrorUnknown = "Unknown"

// ToolchainCargo is the canonical Cargo toolchain name.
const ToolchainCargo = "cargo"

// ToolchainGo is the canonical Go toolchain name.
const ToolchainGo = "go"

// LangRust is the lower-case Rust language key.
const LangRust = "rust"

// LangRustCanonical is the display-canonical Rust language name.
const LangRustCanonical = "Rust"

// LangGo is the Go language key.
const LangGo = "go"

// LangC is the C language key.
const LangC = "c"

// LangJava is the Java language name.
const LangJava = "java"

// LangPython is the Python language name.
const LangPython = "python"

// LSPProviderNative is the tree-sitter based LSP provider.
const LSPProviderNative = "native"

// LSPProviderABCoder is the canonical abcoder MCP LSP provider name (default).
const LSPProviderABCoder = "abcoder-mcp"

// SpecMinerToolName is the canonical PRIM-4 SpecMiner tool name.
const SpecMinerToolName = "specminer"

// DefaultArtifactDir is the canonical artifact/checkpoint directory.
const DefaultArtifactDir = ".artifacts"

// SourceReindexed marks a symbol served from the emitted-fragment index (PRIM-31).
const SourceReindexed = "reindexed"

// SourceFresh marks a symbol requiring a fresh Navigator round-trip.
const SourceFresh = "fresh"

// DefaultRequestFile is the canonical YAML request file name.
const DefaultRequestFile = "magharcm-request.yml"

// DefaultSourceTreeDepth is the maximum directory depth walked by the
// analyzer/planner BuildDirectoryTree probes. Centralised here so the
// analyzer and planner agree on the same traversal radius.
const DefaultSourceTreeDepth = 15

// TranslatedPackagePlaceholder is the fallback Rust crate / Go module / Python
// package name used by the translator when neither the target directory nor
// its parent yields a usable sanitised identifier.
const TranslatedPackagePlaceholder = "translated_project"

// LegacySourceSampleDescriptor is the placeholder source-side prompt fed
// to the verdict panel when role-flip / consensus review asks about the
// original source artifact (the panel itself inspects the translated output).
const LegacySourceSampleDescriptor = "legacy source code"

// MaxGraphRunSteps is the per-invocation run-step ceiling Eino imposes on
// the compiled MAgHARCM graph. Tuned to be safe for the bounded repair
// loop while still aborting runaway cycles.
const MaxGraphRunSteps = 50

// ModularityTrapYears is the age threshold beyond which a design rule is
// flagged as potentially obsolete — the "modularity trap" introduced in
// [[1.0.0 P-71]] (Fleming &amp; Baldwin 2024 retrospective on Design Rules).
// Design rules that were historically optimal but lock in a decomposition
// older than this threshold accumulate translation debt at the rate Lehman
// predicted. The Archaeologist agent surfaces design rules older than this
// threshold to the user for explicit re-validation. Default 15 years matches
// the empirical half-life of design-rule stability reported in P-71.
const ModularityTrapYears = 15

// -------------------------------------------------------------------------
// Logger Scopes (PRIM-Pipeline logging)
// -------------------------------------------------------------------------

// LogScope* constants identify the structured-log namespace each agent emits
// under. Centralising them keeps logger.New(scope) callsites consistent and
// makes log filters / dashboards addressable by symbol rather than free text.
const (
	LogScopeAnalyzer          = "Analyzer"
	LogScopeStrategyRegistry  = "StrategyRegistry"
	LogScopeRecruiter         = "Recruiter"
	LogScopeIterRetrieval     = "IterRetrieval"
	LogScopeNavigator         = "Navigator"
	LogScopeVerdictPanel      = "VerdictPanel"
	LogScopeRoleFlip          = "RoleFlip"
)

// -------------------------------------------------------------------------
// Strategy Thresholds (PRIM-21)
// -------------------------------------------------------------------------

// Mueller migration-strategy gating thresholds. Tuned to match the
// PRIM-21 evaluator described in internal/agents/strategy.go.
const (
	BigBangFileMax          = 3
	BigBangLoCMax           = 500
	PilotFileMin            = 50
	PilotLoCMin             = 10000
	ParallelCutoverFileMin  = 10
)

// StrategyRationale* are the human-readable rationales surfaced by the
// analyzer after strategy selection. Kept here so a single edit propagates
// to prompts, audit logs, and downstream consumers.
const (
	StrategyRationaleBigBang         = "Small self-contained project (<500 LoC, <=3 files): single-pass direct translation."
	StrategyRationalePilot           = "Large-scale codebase (>50 files or >10k LoC): chunked subsystem pilot translation."
	StrategyRationaleFrozenLegacy    = "Legacy codebase without test harness: requires test synthesis and boundary freezing."
	StrategyRationaleParallelCutover = "Modular project with comprehensive test suite: multi-stage parallel module cutover."
	StrategyRationaleIncremental     = "Standard multi-module project: reverse-topological incremental translation."
	StrategyRationaleUnknown         = "No matching strategy selected."
)

// -------------------------------------------------------------------------
// Verdict Vocabulary (PRIM-7 / Optional Checks)
// -------------------------------------------------------------------------

// VerdictEquivalent is the canonical token a judge emits when the source
// and target fragments are functionally equivalent.
const VerdictEquivalent = "EQUIVALENT"

// VerdictNotEquivalent is the canonical token a judge emits when the source
// and target fragments diverge.
const VerdictNotEquivalent = "NOT_EQUIVALENT"

// VerdictJudgeIDPrefix is the prefix used when synthesising per-judge IDs
// ("J1", "J2", …) for trace attribution.
const VerdictJudgeIDPrefix = "J"

// VerdictAliasEquivalent is the set of free-text judge responses that
// classify as EQUIVALENT after case-insensitive normalisation.
var VerdictAliasEquivalent = []string{"EQUIV", "YES", "TRUE", "AGREE", "MATCH", "EQUAL"}

// VerdictAliasNotEquivalent is the set of free-text judge responses that
// classify as NOT_EQUIVALENT after case-insensitive normalisation.
var VerdictAliasNotEquivalent = []string{"NOT EQUIVALENT", "NO", "FALSE", "DISAGREE", "MISMATCH", "DIFFERENT"}

// Verdict is the typed result of an optional check (pass/fail/skipped).
type Verdict string

const (
	VerdictPass    Verdict = "pass"
	VerdictFail    Verdict = "fail"
	VerdictSkipped Verdict = "skipped"
)

// OptionalCheckName constants are the stable identifiers returned by
// OptionalCheck.Name() — referenced by dashboards and regression tests.
const (
	OptionalCheckVerdictPanel   = "PRIM-7-verdict-panel"
	OptionalCheckMockValidator  = "PRIM-8-mock-validator"
	OptionalCheckImplAgnostic   = "PRIM-11-impl-agnostic"
	OptionalCheckWasmOracle     = "PRIM-12-wasm-oracle"
	OptionalCheckRoleFlipGate   = "PRIM-25-role-flip-gate"
)

// -------------------------------------------------------------------------
// RoleFlip Gate (PRIM-25)
// -------------------------------------------------------------------------

// RoleFlip prompt / hint strings consumed by the communicative
// de-hallucination role-flip gate (ChatDev §2.4).
const (
	RoleFlipSystemPrompt    = "you are a critical reviewer who must find at least one bug in the code below"
	RoleFlipAcceptedReason  = "reviewer accepted"
	RoleFlipRetryHint       = "re-check for common bug classes"
	RoleFlipNoDefectToken   = "NO_DEFECT"
)

// ErrRoleFlipGateNotConfigured is the sentinel returned when the role-flip
// gate is invoked without a backing chat model. Surfaces a clear panic-able
// diagnostic at startup rather than a nil-deref at runtime.
var ErrRoleFlipGateNotConfigured = errors.New("roleflip: Model is nil")

// -------------------------------------------------------------------------
// Navigator / Iterative Retrieval (PRIM-26 / PRIM-31)
// -------------------------------------------------------------------------

// ProjectDirPlaceholder is the fallback project root used when an empty path is
// supplied to a Navigator lookup.
const ProjectDirPlaceholder = "."

// ErrNavigatorNoProvider is the sentinel returned when the Navigator is
// invoked without a configured LSPProvider.
var ErrNavigatorNoProvider = errors.New("navigator: no LSP provider configured")

// IterativeContextBudgetBytes is the per-symbol body budget used by the
// IterativeNavigator. Symbols whose stored body exceeds this are trimmed
// before being returned (RepoCoder-style 4 KiB local context window).
const IterativeContextBudgetBytes = 4 * 1024

// -------------------------------------------------------------------------
// Checkpoint Persistence (PRIM-28)
// -------------------------------------------------------------------------

// CheckpointDirMode / CheckpointFileMode are the filesystem modes used
// when creating checkpoint directories and files.
const (
	CheckpointDirMode  = 0o755
	CheckpointFileMode = 0o644
)

// CheckpointFilePattern is the canonical per-iteration checkpoint name
// pattern. sprintf it as fmt.Sprintf(compiletime.CheckpointFilePattern, n).
const CheckpointFilePattern = "iter-%04d.json"

// CheckpointExt is the canonical checkpoint file extension.
const CheckpointExt = ".json"

// DefaultRunID is the run identifier used when no source directory is
// available (e.g. tests, ad-hoc invocations).
const DefaultRunID = "default"
// -------------------------------------------------------------------------
// Source Walker — Directory Skip List (PRIM-14 Archaeologist)
// -------------------------------------------------------------------------

// ArchaeologySkipDirs are directory names the Archaeologist's source
// walker must skip. Centralised so the runner, validator, and TUI all
// reference the same vocabulary. ADR-C-005: hardcoded-by-necessity values
// live in compiletime, not in agent modules.
var ArchaeologySkipDirs = []string{".git", ".artifacts", "node_modules", "target", "vendor", "dist", "build", "testdata"}
// -------------------------------------------------------------------------
// SLM Prompt Contracts (PRIM for 4B–30B / 4B-30B-class models)
// -------------------------------------------------------------------------
//
// Backlink: research/P-50-..-P-53 (forthcoming Sprint 2026-09-07) and
// ADR-C-014. These contracts encode the SLM-aware prompt directives so
// every agent's system prompt emits the same preamble regardless of which
// prompt template the agent loads.

// SLMPromptContractPreamble is the directive block prepended to every
// agent prompt template. It targets 4B–30B / 4B-30B-class models (Qwen2.5-Coder,
// Phi-3, StarCoder2, …) and enforces the SLM-aware output contract:
// no conversational preamble, no markdown code-block fences outside the
// requested artifacts, and stop at the closing delimiter. ADR-C-007.
const SLMPromptContractPreamble = `You are a strict code-translation agent running on a SMALL LANGUAGE MODEL (4B-30B parameters).
Follow every rule below. Do not add conversational text.

RULES:
1. NEVER start with polite preamble ("Sure", "Here is", "Of course"). Begin the response with the FIRST required artifact.
2. NEVER include markdown code-block fences around artifacts unless the schema explicitly requests fenced code.
3. Emit ONLY the requested artifacts in the order specified. No trailing commentary.
4. If a field is unknown, emit the literal token ` + "`" + ErrorUnknown + "`" + ` — never omit.
5. Stop generation after the closing delimiter of the LAST requested artifact.
`

// -------------------------------------------------------------------------
// Spec Lifecycle (PRIM-16)
// -------------------------------------------------------------------------

// SpecLifecyclePhase identifies the active stage in the spec-driven lifecycle.
type SpecLifecyclePhase string

const (
	PhaseConstitution SpecLifecyclePhase = "CONSTITUTION"
	PhaseSpecify      SpecLifecyclePhase = "SPECIFY"
	PhaseDraft        SpecLifecyclePhase = "DRAFT"
	PhasePlan         SpecLifecyclePhase = "PLAN"
	PhaseTasks        SpecLifecyclePhase = "TASKS"
	PhaseReview       SpecLifecyclePhase = "REVIEW"
	PhaseImplement    SpecLifecyclePhase = "IMPLEMENT"
	PhaseApproved     SpecLifecyclePhase = "APPROVED"
	PhaseConverged    SpecLifecyclePhase = "CONVERGED"
)

// -------------------------------------------------------------------------
// Concept Assignment (PRIM-20)
// -------------------------------------------------------------------------

// ConceptLabel* are the human-readable concept names surfaced by the
// ConceptAssigner when clustering source identifiers.
const (
	ConceptLabelValidation    = "Validation & Verification"
	ConceptLabelParsing       = "Parsing & Lexical Analysis"
	ConceptLabelMathStats     = "Mathematical & Statistical Computation"
	ConceptLabelEntityStorage = "Entity & Storage Domain"
)

// ConceptTokenMinLength is the minimum token length considered when
// clustering identifiers into concepts. Shorter tokens are skipped.
const ConceptTokenMinLength = 4

// ConceptDescriptionPlaceholder is the fallback description used when a
// concept binding has no caller-provided description.
const ConceptDescriptionPlaceholder = "(no description)"

// ConceptKeywordCluster binds a human-readable concept label to the list of
// substring keywords that, when found in an identifier (case-insensitive),
// route the identifier into that concept.
type ConceptKeywordCluster struct {
	Label    string
	Keywords []string
}

// DefaultConceptClusters is the canonical keyword → concept table used by
// the PRIM-20 ConceptAssigner.
var DefaultConceptClusters = []ConceptKeywordCluster{
	{Label: ConceptLabelValidation, Keywords: []string{"validat", "check", "verif"}},
	{Label: ConceptLabelParsing, Keywords: []string{"parse", "lex", "token", "scan"}},
	{Label: ConceptLabelMathStats, Keywords: []string{"stat", "math", "calc", "mean"}},
	{Label: ConceptLabelEntityStorage, Keywords: []string{"item", "store", "repo", "record"}},
}

// -------------------------------------------------------------------------
// ArchitectureStabilityLayer is one of the L1/L2/L3 design-rule partitions (PRIM-19).
type ArchitectureStabilityLayer string

// ArchitectureStabilityLayerL1/L2/L3 are the canonical layer identifiers.
const (
	ArchitectureStabilityLayerL1 ArchitectureStabilityLayer = "L1 Interface"
	ArchitectureStabilityLayerL2 ArchitectureStabilityLayer = "L2 Subsystem"
	ArchitectureStabilityLayerL3 ArchitectureStabilityLayer = "L3 Leaf"
)

// ArchitectureStabilityDescriptionL1/L2/L3 are short descriptions paired
// with each layer for report rendering.
const (
	ArchitectureStabilityDescriptionL1 = "High-stability design rule contract"
	ArchitectureStabilityDescriptionL2 = "Intermediate subsystem coordinator"
	ArchitectureStabilityDescriptionL3 = "Volatile concrete leaf implementation"
)

// -------------------------------------------------------------------------
// Comprehension Recognition (PRIM-22)
// -------------------------------------------------------------------------

// ComprehensionRecognition* are the labels the Comprehension pipeline
// surfaces when it recognises a library / framework fingerprint.
const (
	ComprehensionRecognitionMath    = "Math/Numerics Library"
	ComprehensionRecognitionTesting = "Testing Framework"
)

// ComprehensionExplanationDefault is the placeholder explanation rendered
// when the Comprehension pipeline has nothing specific to report.
const ComprehensionExplanationDefault = "Decomposed into structural units with breadth-first linear traversal ordering."

// -------------------------------------------------------------------------
// Recruiter Tool / Agent Names (PRIM-29)
// -------------------------------------------------------------------------

// Tool* are the canonical tool names the Recruiter surfaces in a
// RecruitmentPlan. Centralised so log filters and dashboards reference
// the same vocabulary as the producer.
const (
	ToolValidator            = "validator"
	ToolDiagnostics          = "diagnostics"
	ToolChunkedTranslator    = "chunked_translator"
)

// Agent* are the canonical downstream-agent names the Recruiter surfaces
// in a RecruitmentPlan.
const (
	AgentValidator         = "validator"
	AgentVerdictPanel      = "verdict_panel"
	AgentRoleFlip          = "roleflip"
	AgentAdversarialSuite  = "adversarial_suite"
	AgentPlateauBreaker    = "plateau_breaker"
	AgentRecruitTranslatorV2 = "recruit_translator_v2"
)

// -------------------------------------------------------------------------
// Strict Enums
// -------------------------------------------------------------------------

// CompilationStatus represents binary per-project compilation outcomes.
// There is no partial compilation percentage: status is strictly Pass or Fail.
type CompilationStatus string

const (
	CompilationStatusPass CompilationStatus = "PASS"
	CompilationStatusFail CompilationStatus = "FAIL"
)

// StrategyKind names one of Mueller's five migration strategies (PRIM-21).
type StrategyKind string

const (
	StrategyBigBang         StrategyKind = "BIG_BANG"
	StrategyIncremental     StrategyKind = "INCREMENTAL"
	StrategyPilot           StrategyKind = "PILOT"
	StrategyFrozenLegacy    StrategyKind = "FROZEN_LEGACY"
	StrategyParallelCutover StrategyKind = "PARALLEL_CUTOVER"
)

// WorkUnitStatus defines P17 Blackboard scheduler lifecycle states (CAID).
type WorkUnitStatus string

const (
	WorkUnitPending   WorkUnitStatus = "PENDING"
	WorkUnitClaimed   WorkUnitStatus = "CLAIMED"
	WorkUnitCompleted WorkUnitStatus = "COMPLETED"
	WorkUnitFailed    WorkUnitStatus = "FAILED"
	WorkUnitTerminal  WorkUnitStatus = "TERMINAL"
)

// -------------------------------------------------------------------------
// Shared Magic Values (ADR-C-005 Sweep)
// -------------------------------------------------------------------------

// RoleFlipGateToolSuffix is the trailing vocabulary piece appended to
// compiletime.AgentRoleFlip ("roleflip") when the recruiter enumerates
// the PRIM-25 communicative-de-hallucination gate as a standalone tool.
// Centralised here so the recruiter, the runner's status surface, and any
// future log-filter wiring reference the same suffix instead of repeating
// the concatenation at every callsite.
const RoleFlipGateToolSuffix = "_gate"

// MinRealTestsFloor is the absolute lower bound applied to the per-run
// MinRealTests threshold (validator.go). Distinct-source-count * 2 is the
// scale-up rule; this floor guarantees a non-zero target for tiny repos
// where 1 source file would otherwise map to a 2-test minimum.
const MinRealTestsFloor = 5

// MinRealTestsMultiplier is the per-distinct-source multiplier the
// validator applies when computing MinRealTests. Tuned in PR-#54 so a
// single source file with N functions maps to roughly 2N expected real
// tests; keeping the multiplier centralised prevents drift between the
// validator's primary cascade and any downstream remediation loop.
const MinRealTestsMultiplier = 2

// ASTEmptyElementsSizeThreshold is the minimum source-file size (bytes)
// below which the validator's pre-compile AST check skips emitting a
// "0 AST elements" warning. Tiny stub files (mod.rs, mod.go) routinely
// parse to zero elements but are correct as-is; the byte floor filters
// the noise without losing real syntax-error coverage on larger files.
const ASTEmptyElementsSizeThreshold = 200

// TestResultNotParsedSentinel is the RunTestsOutput.RealTests value used
// by tools.RunProjectTests to flag a target language whose test runner
// does not emit a parseable per-test count. Validator treats any
// non-negative RealTests value as authoritative.
const TestResultNotParsedSentinel = -1

// -------------------------------------------------------------------------
// Wave-16 SLM-Era Pattern Constants (Sprint 2026-09-27)
// -------------------------------------------------------------------------
//
// These constants back the opt-in patterns documented in METHODOLOGY.md §11
// and Architecture.md §8. Each is read by its consumer via YAML config; the
// compile-time default is the safe-blind value; the YAML value opts in to the
// wave-16 pattern.

// MaxDiscoveryArms is the upper bound on MAB arms explored by the
// bounded-discovery substrate (BOAD, [[1.0.0 P-123]]). Larger values
// increase combinatorial coverage at the cost of evaluation budget.
const MaxDiscoveryArms = 16

// MaxDiscoverySamplesPerArm is the upper bound on per-arm samples drawn
// during MAB evaluation. Combined with MaxDiscoveryArms this caps the total
// evaluation budget for a single bounded-discovery run.
const MaxDiscoverySamplesPerArm = 3

// MaxMemoryTriples is the upper bound on the persistent-memory substrate
// (ReasoningBank, [[1.0.0 P-122]]). The memory store evicts the lowest-reward
// triple when the bound is exceeded.
const MaxMemoryTriples = 256

// MemoryTripleRewardEMA is the exponential-moving-average decay factor
// applied to per-iteration rewards when ranking memory triples for eviction.
// Larger values give more weight to recent rewards; smaller values smooth.
const MemoryTripleRewardEMA = 0.3
// MaTTSDefaultBudget is the default iteration ceiling for the Memory-aware
// Test-Time Scaling loop ([[1.0.0 P-122]] ReasoningBank anchor) invoked from
// memorystore.ApplyMaTTS. When a caller passes budget <= 0 the package falls
// back to this value so the contract is satisfied without forcing every
// callsite to thread the constant through. Value 10 matches the canonical
// wave-16 ReasoningBank MaTTS budget and the bounded-discover ceiling; a
// human engineer may override by passing an explicit positive budget at
// the callsite rather than mutating this constant.
const MaTTSDefaultBudget = 10

// BanditPolicyDefault is the safe-blind default for the migration-strategy
// registry when bandit_policy is not opted in via YAML. "greedy" matches the
// pre-wave-16 behaviour; "ucb" and "thompson" are wave-16 opt-ins.
const BanditPolicyDefault = "greedy"
const BanditPolicyUCB = "ucb"
const BanditPolicyThompson = "thompson"

// DiscoveryExplorationConstant is the UCB exploration constant c used by the
// bandit-policy registry when BanditPolicyUCB is selected. The classical UCB1
// default is sqrt(2); MAgHARCM uses 1.5 to favour exploitation given the
// bounded evaluation budget.
const DiscoveryExplorationConstant = 1.5

// -------------------------------------------------------------------------
// Must Pattern Helpers (No Fallbacks)
// -------------------------------------------------------------------------

// Must panics if err is non-nil; otherwise it returns v.
// Used for compile-time and startup initializations where failure is fatal.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("compiletime.Must: %v", err))
	}
	return v
}

// MustNotNil panics if ptr is nil.
func MustNotNil[T any](ptr *T, fieldName string) *T {
	if ptr == nil {
		panic(fmt.Sprintf("compiletime.MustNotNil: %s must not be nil", fieldName))
	}
	return ptr
}

// MustNotEmpty panics if s is empty.
func MustNotEmpty(s string, fieldName string) string {
	if s == "" {
		panic(fmt.Sprintf("compiletime.MustNotEmpty: %s must not be empty", fieldName))
	}
	return s
}

