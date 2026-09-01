package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// EventType categorizes execution events in the translation pipeline.
type EventType string

const (
	EventAgentStart EventType = "agent_start"
	EventAgentEnd   EventType = "agent_end"
	EventStep       EventType = "step"
	EventToolCall   EventType = "tool_call"
	EventValidation EventType = "validation"
	EventLog        EventType = "log"
	EventWarning    EventType = "warning"
	EventError      EventType = "error"
)

// Event is the structured record for any pipeline event.
type Event struct {
	Timestamp time.Time              `json:"timestamp"`
	Type      EventType              `json:"type"`
	Agent     string                 `json:"agent,omitempty"`
	Tool      string                 `json:"tool,omitempty"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

var (
	outMu sync.Mutex
	out   io.Writer = os.Stdout
)

// SetOutput swaps the default sink. Used by tests; production code never
// needs to call this because stdout is the only legitimate target.
func SetOutput(w io.Writer) {
	outMu.Lock()
	defer outMu.Unlock()
	if w == nil {
		w = os.Stdout
	}
	out = w
}

func writeLine(format string, args ...any) {
	outMu.Lock()
	defer outMu.Unlock()
	fmt.Fprintf(out, format+"\n", args...)
}

// Emit sends an event to the configured output.
func Emit(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	ts := e.Timestamp.Format("15:04:05")
	switch e.Type {
	case EventAgentStart, EventAgentEnd:
		writeLine("[%s] [%s] %s", ts, e.Agent, e.Message)
	case EventToolCall:
		writeLine("[%s] [Tool: `%s`] %s", ts, e.Tool, e.Message)
	case EventStep:
		writeLine("[%s] %s", ts, e.Message)
	case EventValidation:
		writeLine("[%s] [Validation] %s", ts, e.Message)
	case EventWarning:
		writeLine("[%s] [WARNING] %s", ts, e.Message)
	case EventError:
		writeLine("[%s] [ERROR] %s", ts, e.Message)
	default:
		if e.Agent != "" {
			writeLine("[%s] [%s] %s", ts, e.Agent, e.Message)
		} else {
			writeLine("[%s] %s", ts, e.Message)
		}
	}
}

func LogAgent(agent, format string, args ...any) {
	Emit(Event{Type: EventAgentStart, Agent: agent, Message: fmt.Sprintf(format, args...)})
}

func LogTool(toolName, format string, args ...any) {
	Emit(Event{Type: EventToolCall, Tool: toolName, Message: fmt.Sprintf(format, args...)})
}

func LogStep(format string, args ...any) {
	Emit(Event{Type: EventStep, Message: fmt.Sprintf(format, args...)})
}

func LogValidation(format string, args ...any) {
	Emit(Event{Type: EventValidation, Message: fmt.Sprintf(format, args...)})
}

func LogWarning(format string, args ...any) {
	Emit(Event{Type: EventWarning, Message: fmt.Sprintf(format, args...)})
}

func LogError(format string, args ...any) {
	Emit(Event{Type: EventError, Message: fmt.Sprintf(format, args...)})
}
