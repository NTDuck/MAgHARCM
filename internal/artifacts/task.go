package artifacts

import "MAgHARCM/internal/agents"

// DocumentWrapper extracted to internal/agents/analyzer.go for Locality of Behaviour (PRIM-24).
type DocumentWrapper[T any] = agents.DocumentWrapper[T]
