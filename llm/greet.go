package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (p *OpenAICompatibleProvider) Greet(ctx context.Context, req GreetRequest) (string, error) {
	if !p.cfg.IsConfigured() {
		return "", ErrUnconfigured
	}

	body, err := json.Marshal(chatCompletionsRequest{
		Model: p.cfg.Model,
		Messages: []chatMessage{
			{Role: "system", Content: greetSystemPrompt(p.cfg.Personality)},
			{Role: "user", Content: buildGreetUserPrompt(req)},
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
		return "", fmt.Errorf("provider returned no choices")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

func greetSystemPrompt(personality string) string {
	if personality != "" {
		return personality
	}
	return "You are Nova, a friendly terminal companion. " +
		"Generate a brief, warm greeting for the user based on their current context. " +
		"Rules:\n" +
		"- 1 to 3 sentences, conversational, like a colleague saying hi\n" +
		"- Acknowledge the project/context naturally if present\n" +
		"- No emojis, no markdown, no greeting-card vibes\n" +
		"- Don't introduce yourself as an AI assistant"
}

func buildGreetUserPrompt(req GreetRequest) string {
	var b strings.Builder
	b.WriteString("The user just opened their terminal. Here is their current context:\n\n")
	if req.TimeOfDay != "" {
		b.WriteString("- Time of day: " + req.TimeOfDay + "\n")
	}
	if req.ProjectName != "" || req.ProjectType != "" {
		name := req.ProjectName
		if name == "" {
			name = "(unnamed)"
		}
		if req.ProjectType != "" {
			b.WriteString("- Project: " + name + " (" + req.ProjectType + ")\n")
		} else {
			b.WriteString("- Project: " + name + "\n")
		}
	}
	if req.Cwd != "" {
		b.WriteString("- Directory: " + req.Cwd + "\n")
	}
	b.WriteString("\nWrite a brief greeting.")
	return b.String()
}
