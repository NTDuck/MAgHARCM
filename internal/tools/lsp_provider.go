package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"MAgHARCM/internal/logger"
)

// LSPProvider defines the unified interface for language server intelligence providers.
type LSPProvider interface {
	Name() string
	GetDefinition(ctx context.Context, input *DefinitionInput) (*DefinitionOutput, error)
	GetDiagnostics(ctx context.Context, input *DiagnosticsInput) (*DiagnosticsOutput, error)
	EditFile(ctx context.Context, input *EditFileInput) (*EditFileOutput, error)
	GetHover(ctx context.Context, input *HoverInput) (*HoverOutput, error)
	GetReferences(ctx context.Context, input *ReferencesInput) (*ReferencesOutput, error)
	RenameSymbol(ctx context.Context, input *RenameSymbolInput) (*RenameSymbolOutput, error)
}

// NativeLSPProvider uses local Tree-Sitter and direct static analysis.
type NativeLSPProvider struct{}

func NewNativeLSPProvider() *NativeLSPProvider {
	return &NativeLSPProvider{}
}

func (n *NativeLSPProvider) Name() string {
	return "native"
}

func (n *NativeLSPProvider) GetDefinition(ctx context.Context, input *DefinitionInput) (*DefinitionOutput, error) {
	return ExecuteDefinition(ctx, input)
}

func (n *NativeLSPProvider) GetDiagnostics(ctx context.Context, input *DiagnosticsInput) (*DiagnosticsOutput, error) {
	return ExecuteDiagnostics(ctx, input)
}

func (n *NativeLSPProvider) EditFile(ctx context.Context, input *EditFileInput) (*EditFileOutput, error) {
	return ExecuteEditFile(ctx, input)
}

func (n *NativeLSPProvider) GetHover(ctx context.Context, input *HoverInput) (*HoverOutput, error) {
	return ExecuteHover(ctx, input)
}

func (n *NativeLSPProvider) GetReferences(ctx context.Context, input *ReferencesInput) (*ReferencesOutput, error) {
	return ExecuteReferences(ctx, input)
}

func (n *NativeLSPProvider) RenameSymbol(ctx context.Context, input *RenameSymbolInput) (*RenameSymbolOutput, error) {
	return ExecuteRenameSymbol(ctx, input)
}

// ABCoderMcpProvider integrates with CloudWeGo abcoder (https://github.com/cloudwego/abcoder)
// via Model Context Protocol (MCP) tool invocation over JSON-RPC stdio.
type ABCoderMcpProvider struct {
	nativeFallback LSPProvider
	serverCmd      string
	serverArgs     []string
	mu             sync.Mutex
}

// NewABCoderMcpProvider creates an ABCoder MCP provider with native fallback.
func NewABCoderMcpProvider(cmd string, args ...string) *ABCoderMcpProvider {
	if cmd == "" {
		cmd = "abcoder"
	}
	return &ABCoderMcpProvider{
		nativeFallback: NewNativeLSPProvider(),
		serverCmd:      cmd,
		serverArgs:     args,
	}
}

func (a *ABCoderMcpProvider) Name() string {
	return "abcoder-mcp"
}

// callMcp executes an MCP JSON-RPC call against the abcoder MCP server.
func (a *ABCoderMcpProvider) callMcp(ctx context.Context, toolName string, arguments map[string]any) ([]byte, error) {
	if _, err := exec.LookPath(a.serverCmd); err != nil {
		return nil, fmt.Errorf("abcoder binary not found in PATH: %w", err)
	}

	reqPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      time.Now().UnixNano(),
		"method":  "tools/call",
		"params": map[string]any{
			"name":      toolName,
			"arguments": arguments,
		},
	}
	data, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, a.serverCmd, a.serverArgs...)
	cmd.Stdin = bytes.NewReader(data)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("abcoder mcp execution failed: %w", err)
	}

	return outBuf.Bytes(), nil
}

func (a *ABCoderMcpProvider) GetDefinition(ctx context.Context, input *DefinitionInput) (*DefinitionOutput, error) {
	args := map[string]any{"symbol": input.Symbol, "file_path": input.FilePath, "project_dir": input.Project}
	if out, err := a.callMcp(ctx, "definition", args); err == nil && len(out) > 0 {
		var res DefinitionOutput
		if err := json.Unmarshal(out, &res); err == nil && len(res.Definitions) > 0 {
			logger.LogTool("abcoder_mcp", "Resolved definition for `%s` via `abcoder`", input.Symbol)
			return &res, nil
		}
	}
	// Fallback to Native LSP
	return a.nativeFallback.GetDefinition(ctx, input)
}

func (a *ABCoderMcpProvider) GetDiagnostics(ctx context.Context, input *DiagnosticsInput) (*DiagnosticsOutput, error) {
	args := map[string]any{"file_path": input.FilePath, "project_dir": input.ProjectDir}
	if out, err := a.callMcp(ctx, "diagnostics", args); err == nil && len(out) > 0 {
		var res DiagnosticsOutput
		if err := json.Unmarshal(out, &res); err == nil {
			logger.LogTool("abcoder_mcp", "Retrieved diagnostics for `%s` via `abcoder`", input.FilePath)
			return &res, nil
		}
	}
	return a.nativeFallback.GetDiagnostics(ctx, input)
}

func (a *ABCoderMcpProvider) EditFile(ctx context.Context, input *EditFileInput) (*EditFileOutput, error) {
	args := map[string]any{"file_path": input.FilePath, "edits": input.Edits, "project_dir": input.ProjectDir}
	if out, err := a.callMcp(ctx, "edit_file", args); err == nil && len(out) > 0 {
		var res EditFileOutput
		if err := json.Unmarshal(out, &res); err == nil && res.Success {
			return &res, nil
		}
	}
	return a.nativeFallback.EditFile(ctx, input)
}

func (a *ABCoderMcpProvider) GetHover(ctx context.Context, input *HoverInput) (*HoverOutput, error) {
	args := map[string]any{"symbol": input.Symbol, "file_path": input.FilePath, "project_dir": input.ProjectDir}
	if out, err := a.callMcp(ctx, "hover", args); err == nil && len(out) > 0 {
		var res HoverOutput
		if err := json.Unmarshal(out, &res); err == nil && res.Found {
			return &res, nil
		}
	}
	return a.nativeFallback.GetHover(ctx, input)
}

func (a *ABCoderMcpProvider) GetReferences(ctx context.Context, input *ReferencesInput) (*ReferencesOutput, error) {
	args := map[string]any{"symbol": input.Symbol, "project_dir": input.ProjectDir}
	if out, err := a.callMcp(ctx, "references", args); err == nil && len(out) > 0 {
		var res ReferencesOutput
		if err := json.Unmarshal(out, &res); err == nil && len(res.References) > 0 {
			return &res, nil
		}
	}
	return a.nativeFallback.GetReferences(ctx, input)
}

func (a *ABCoderMcpProvider) RenameSymbol(ctx context.Context, input *RenameSymbolInput) (*RenameSymbolOutput, error) {
	args := map[string]any{"old_name": input.OldName, "new_name": input.NewName, "project_dir": input.ProjectDir}
	if out, err := a.callMcp(ctx, "rename_symbol", args); err == nil && len(out) > 0 {
		var res RenameSymbolOutput
		if err := json.Unmarshal(out, &res); err == nil && res.Success {
			return &res, nil
		}
	}
	return a.nativeFallback.RenameSymbol(ctx, input)
}

// GetLSPProvider creates the configured LSP provider ("native" or "abcoder").
func GetLSPProvider(providerName string) LSPProvider {
	switch strings.ToLower(strings.TrimSpace(providerName)) {
	case "abcoder", "abcoder-mcp", "mcp":
		return NewABCoderMcpProvider("abcoder")
	default:
		return NewNativeLSPProvider()
	}
}
