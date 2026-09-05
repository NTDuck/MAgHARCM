package agents

import "MAgHARCM/internal/compiletime"

// State is the canonical pipeline-wide context type re-exported from
// internal/compiletime/state.go. Locality of Behaviour is preserved via
// type aliases — this file only declares the State alias so cross-file
// references inside package agents resolve without a package prefix.
type State = compiletime.State
