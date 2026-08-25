package prompts

import (
	"os"

	. "MAgHARCM/internal/patterns"
)

// MustLoad reads a prompt markdown file from the given path using the Must pattern.
func MustLoad(path string) string {
	return string(Must(os.ReadFile(path)))
}
