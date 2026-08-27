package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// EventType categorizes execution events in the translation pipeline.
type EventType string

const (
	EventAgentStart   EventType = "agent_start"
	EventAgentEnd     EventType = "agent_end"
	EventStep         EventType = "step"
	EventToolCall     EventType = "tool_call"
	EventValidation   EventType = "validation"
	EventArtifactSave EventType = "artifact_save"
	EventLog          EventType = "log"
	EventWarning      EventType = "warning"
	EventError        EventType = "error"
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

// Sink defines an interface for consuming log events.
type Sink interface {
	WriteEvent(event Event)
	Close() error
}

// ConsoleSink formats events into human-readable real-time progress.
type ConsoleSink struct {
	out io.Writer
	mu  sync.Mutex
}

// NewConsoleSink creates a console sink writing to stdout.
func NewConsoleSink(out io.Writer) *ConsoleSink {
	if out == nil {
		out = os.Stdout
	}
	return &ConsoleSink{out: out}
}

// WriteEvent formats and prints the event to the console.
func (c *ConsoleSink) WriteEvent(e Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ts := e.Timestamp.Format("15:04:05")
	switch e.Type {
	case EventAgentStart, EventAgentEnd:
		fmt.Fprintf(c.out, "[%s] [%-10s] %s\n", ts, e.Agent, e.Message)
	case EventToolCall:
		fmt.Fprintf(c.out, "[%s]   -> [Tool: %-18s] %s\n", ts, e.Tool, e.Message)
	case EventStep:
		fmt.Fprintf(c.out, "[%s]     - %s\n", ts, e.Message)
	case EventValidation:
		fmt.Fprintf(c.out, "[%s] [Validation] %s\n", ts, e.Message)
	case EventArtifactSave:
		fmt.Fprintf(c.out, "[%s]   [Artifact] %s\n", ts, e.Message)
	case EventWarning:
		fmt.Fprintf(c.out, "[%s] [WARNING] %s\n", ts, e.Message)
	case EventError:
		fmt.Fprintf(c.out, "[%s] [ERROR] %s\n", ts, e.Message)
	default:
		if e.Agent != "" {
			fmt.Fprintf(c.out, "[%s] [%s] %s\n", ts, e.Agent, e.Message)
		} else {
			fmt.Fprintf(c.out, "[%s] %s\n", ts, e.Message)
		}
	}
}

// Close is a no-op for ConsoleSink.
func (c *ConsoleSink) Close() error {
	return nil
}

// JSONFileSink writes JSON-serialized events to a file.
type JSONFileSink struct {
	file *os.File
	mu   sync.Mutex
}

// NewJSONFileSink creates or appends to a JSONL file.
func NewJSONFileSink(path string) (*JSONFileSink, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &JSONFileSink{file: f}, nil
}

// WriteEvent encodes event as a JSON line.
func (j *JSONFileSink) WriteEvent(e Event) {
	j.mu.Lock()
	defer j.mu.Unlock()

	data, err := json.Marshal(e)
	if err == nil {
		_, _ = j.file.Write(append(data, '\n'))
	}
}

// Close closes the underlying file.
func (j *JSONFileSink) Close() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.file != nil {
		return j.file.Close()
	}
	return nil
}

// Manager orchestrates dispatching events to all registered sinks.
type Manager struct {
	sinks []Sink
	mu    sync.RWMutex
}

var globalManager = &Manager{
	sinks: []Sink{NewConsoleSink(os.Stdout)},
}

// AddSink registers a new sink with the global logger.
func AddSink(s Sink) {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()
	globalManager.sinks = append(globalManager.sinks, s)
}

// SetFileSink configures an event JSONL sink to the specified path.
func SetFileSink(path string) error {
	sink, err := NewJSONFileSink(path)
	if err != nil {
		return err
	}
	AddSink(sink)
	return nil
}

// Emit sends an event to all registered sinks.
func Emit(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()
	for _, s := range globalManager.sinks {
		s.WriteEvent(e)
	}
}

// LogAgent emits an agent-level milestone event.
func LogAgent(agent, format string, args ...any) {
	Emit(Event{
		Type:    EventAgentStart,
		Agent:   agent,
		Message: fmt.Sprintf(format, args...),
	})
}

// LogTool emits a tool call/result event.
func LogTool(toolName, format string, args ...any) {
	Emit(Event{
		Type:    EventToolCall,
		Tool:    toolName,
		Message: fmt.Sprintf(format, args...),
	})
}

// LogStep emits a step progress event.
func LogStep(format string, args ...any) {
	Emit(Event{
		Type:    EventStep,
		Message: fmt.Sprintf(format, args...),
	})
}

// LogValidation emits a validation event.
func LogValidation(format string, args ...any) {
	Emit(Event{
		Type:    EventValidation,
		Message: fmt.Sprintf(format, args...),
	})
}

// LogArtifact emits an artifact persistence event.
func LogArtifact(format string, args ...any) {
	Emit(Event{
		Type:    EventArtifactSave,
		Message: fmt.Sprintf(format, args...),
	})
}

// LogWarning emits a warning event.
func LogWarning(format string, args ...any) {
	Emit(Event{
		Type:    EventWarning,
		Message: fmt.Sprintf(format, args...),
	})
}

// LogError emits an error event.
func LogError(format string, args ...any) {
	Emit(Event{
		Type:    EventError,
		Message: fmt.Sprintf(format, args...),
	})
}
