package artifacts

import "MAgHARCM/internal/agents"

// Artifact types extracted to internal/agents/analyzer.go for Locality of Behaviour (PRIM-24).
// Re-exported here for backward compatibility.
type SourceProjectResearch = agents.SourceProjectResearch
type ThirdPartyLibraryAnalysis = agents.ThirdPartyLibraryAnalysis
type LibraryMapping = agents.LibraryMapping
type TargetProjectDesign = agents.TargetProjectDesign
type AnalyzerOutput = agents.AnalyzerOutput
