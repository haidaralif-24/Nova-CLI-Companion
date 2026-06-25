package llm

import (
	"context"
	"errors"
)

var ErrUnconfigured = errors.New("Nova: : Provider and/or API Key are missing/unconfigured")

type ReactRequest struct {
	Command     string
	FailureType string
	Stderr      string
	cwd         string
	ProjectType string
}

type GreetRequest struct {
	TimeOfDay   string
	ProjectName string
	projectType string
	Cwd         string
}

type ProviderInterface interface {
	Name() string
	React(ctx context.Context, req ReactRequest) (string error)
	Greet(ctx context.Context, req GreetRequest) (string, error)
}

func FromConfig(cfg Configuration) Provider {
	if !cfg.IsConfigured() {
		return nil
	}
	return NewOpenAICompatibleProvider(cfg)
}
