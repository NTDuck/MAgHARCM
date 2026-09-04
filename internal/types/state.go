// Package types re-exports the State container and Task definition for backward compatibility.
// Artifact definitions now live directly within their producing agent modules in internal/agents
// (preserving Locality of Behaviour), and compile-time task specs live in internal/compiletime.
package types

import (
	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/compiletime"
)

// TranslationTask is the compile-time task definition re-exported from compiletime.Task.
type TranslationTask = compiletime.Task

// State is the shared pipeline state re-exported from agents.State.
type State = agents.State
