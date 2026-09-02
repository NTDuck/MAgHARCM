package llm

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	ollamaapi "github.com/eino-contrib/ollama/api"

	"MAgHARCM/internal/config"
)

type Models struct {
	Reasoning      model.BaseChatModel
	Coding         model.BaseChatModel
	BaseURL        string
	ReasoningModel string
	CodingModel    string
}

// PrepareReasoning evicts the coding model from VRAM so reasoning has full memory.
func (m *Models) PrepareReasoning() {
	if m != nil {
		UnloadModel(m.BaseURL, m.CodingModel)
	}
}

// PrepareCoding evicts the reasoning model from VRAM so coding has full memory.
func (m *Models) PrepareCoding() {
	if m != nil {
		UnloadModel(m.BaseURL, m.ReasoningModel)
	}
}

// defaultChatOptions sets hard caps and stop tokens that prevent base
// (non-instruction-tuned) models from rambling past the answer. Without
// these, qwen2.5-coder:7b-base-q5_0 in particular can emit thousands of
// tokens of unrelated text after the response is complete — hanging the
// generation until Ollama's default num_predict=-1 cut-off fires much
// later.
//
// NumCtx is sized to fit the analyzer prompt (DirectoryTree +
// StructureSummary + SourceFilesContent), which routinely exceeds 1k
// tokens on Sample 2/4. Empirically, NumCtx=1024 caused ollama HTTP 500
// ("cannot meet free memory target of 1024 MiB") because the KV-cache
// allocator could not satisfy the prompt. 8192 is the smallest value
// that survives the prompt and the foreign llama-server VRAM contention
// on port 8080. See research/_meta/DRIFT_LOG.md for the v4 → v5 revert.
// Stop tokens match the delimiters our prompts ask the model to emit
// at the end of a structured response. `<|im_end|>` is the chatml
// token used by Qwen-style models. We deliberately do NOT stop on
// `\n### ` because our own prompts emit `### FILE:` / `### Source File:`
// / `### Fragment:` headers inside legitimate structured output, and
// premature termination on those headers would truncate every response
// mid-generation.
var defaultChatOptions = &ollamaapi.Options{
	Runner:     ollamaapi.Runner{NumCtx: 8192},
	NumPredict: 4096,
	Stop: []string{
		"<|im_end|>",
		"<|endoftext|>",
	},
	RepeatLastN: 64,
}

// NewModels initializes the Ollama ChatModels for reasoning and coding.
func NewModels(ctx context.Context, cfg *config.Config) (*Models, error) {
	requestTimeout := cfg.Timeout
	if requestTimeout <= 0 || requestTimeout < 600*time.Second {
		requestTimeout = 600 * time.Second
	}
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	reasoningModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL:    cfg.OllamaBaseURL,
		Model:      cfg.ReasoningModel,
		Timeout:    requestTimeout,
		HTTPClient: httpClient,
		Options:    defaultChatOptions,
	})
	if err != nil {
		return nil, err
	}

	codingModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL:    cfg.OllamaBaseURL,
		Model:      cfg.CodingModel,
		Timeout:    requestTimeout,
		HTTPClient: httpClient,
		Options:    defaultChatOptions,
	})
	if err != nil {
		return nil, err
	}

	return &Models{
		Reasoning:      reasoningModel,
		Coding:         codingModel,
		BaseURL:        cfg.OllamaBaseURL,
		ReasoningModel: cfg.ReasoningModel,
		CodingModel:    cfg.CodingModel,
	}, nil
}

// UnloadModel instructs Ollama to evict a model from VRAM to prevent out-of-memory
// errors on memory-constrained GPUs when switching between reasoning and coding models.
func UnloadModel(baseURL, modelName string) {
	if baseURL == "" || modelName == "" {
		return
	}
	url := strings.TrimRight(baseURL, "/") + "/api/generate"
	payload := fmt.Sprintf(`{"model":"%s","keep_alive":0}`, modelName)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err == nil && resp != nil {
		_ = resp.Body.Close()
	}
}

