package agents

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed prompts/*.md
var promptsFS embed.FS

// renderPrompt loads an embedded markdown prompt template and executes it with data.
func renderPrompt(name string, data any) (string, error) {
	content, err := promptsFS.ReadFile("prompts/" + name)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded prompt %s: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute prompt template %s: %w", name, err)
	}

	return buf.String(), nil
}
