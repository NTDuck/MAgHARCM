package logger

import (
	"fmt"
	"time"
)

// LogAgent prints an agent milestone with consistent formatting.
func LogAgent(agent, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] [%-10s] %s\n", timestamp, agent, msg)
}

// LogTool prints a tool invocation and result.
func LogTool(toolName, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s]   ↳ ⚙ [Tool: %-18s] %s\n", timestamp, toolName, msg)
}

// LogStep prints a sub-step.
func LogStep(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s]     • %s\n", timestamp, msg)
}
