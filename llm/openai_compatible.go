// Package llm provides LLM provider abstraction for Nova.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAICompatibleProvider calls any OpenAI-compatible chat completions API.
type OpenAICompatibleProvider struct {
	cfg    Configuration
	client *http.Client
}

// NewOpenAICompatibleProvider returns a provider that calls the configured API.
func NewOpenAICompatibleProvider(cfg Configuration) *OpenAICompatibleProvider {
	return &OpenAICompatibleProvider{
		cfg:    cfg,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// Name returns the provider name from config, or a default.
func (p *OpenAICompatibleProvider) Name() string {
	if p.cfg.Provider != "" {
		return p.cfg.Provider
	}
	return "openai-compatible"
}

func (p *OpenAICompatibleProvider) React(ctx context.Context, req ReactRequest) (string, error) {
	if !p.cfg.IsConfigured() {
		return "", ErrUnconfigured
	}

	body, err := json.Marshal(chatCompletionsRequest{
		Model: p.cfg.Model,
		Messages: []chatMessage{
			{Role: "system", Content: buildSystemPrompt(p.cfg.Personality)},
			{Role: "user", Content: buildUserPrompt(req)},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := p.endpointURL()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if p.cfg.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s returned HTTP %d: %s",
			url, resp.StatusCode, truncate(string(respBody), 200))
	}

	var chatResp chatCompletionsResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return "", errors.New("provider returned no choices")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

func (p *OpenAICompatibleProvider) endpointURL() string {
	base := p.cfg.BaseURL
	if base == "" {
		return "https://api.openai.com/v1/chat/completions"
	}
	base = strings.TrimRight(base, "/")
	if !strings.HasSuffix(base, "/v1") {
		base += "/v1"
	}
	return base + "/chat/completions"
}

func buildSystemPrompt(personality string) string {
	if personality != "" {
		return personality
	}
	return "You are Nova, a concise terminal companion. " +
		"Explain command failures briefly and suggest fixes. " +
		"Keep responses under 5 sentences unless asked otherwise."
}

func buildUserPrompt(req ReactRequest) string {
	var b strings.Builder
	b.WriteString("A command failed in the terminal. Help me understand why and suggest a fix.\n\n")
	if req.ProjectType != "" {
		b.WriteString("Project type: " + req.ProjectType + "\n")
	}
	if req.Cwd != "" {
		b.WriteString("Working directory: " + req.Cwd + "\n")
	}
	if req.Command != "" {
		b.WriteString("Command: " + req.Command + "\n")
	}
	if req.FailureType != "" {
		b.WriteString("Failure type: " + req.FailureType + "\n")
	}
	if req.Stderr != "" {
		b.WriteString("Stderr:\n```\n" + req.Stderr + "\n```\n")
	}
	b.WriteString("\nRespond with: (1) one-sentence cause, (2) suggested fix command.")
	return b.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionsRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatCompletionsResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}
