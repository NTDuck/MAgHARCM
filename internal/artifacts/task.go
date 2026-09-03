// Package artifacts holds the structured output types that flow through
// the translation pipeline. The package has zero upstream dependencies so
// any other package can import it without creating an import cycle; this
// is the key property that lets us share State fragments across producers
// (internal/agents) and consumers (internal/types, internal/graph).
package artifacts

// DocumentWrapper keeps both structured data and markdown representation.
type DocumentWrapper[T any] struct {
	Data        T      `json:"data"`
	RawMarkdown string `json:"raw_markdown"`
}
