# Nova — AI CLI Companion (In Development)

An AI-powered terminal companion that lives alongside your shell.

Nova observes your workflow, understands project context, catches common command mistakes, and remembers your work across sessions — with or without an AI provider.

> **Terminal Companion > Terminal Chatbot.** AI is an enhancement layer. If every AI provider disappeared tomorrow, Nova should still be useful.

## Features

- **Project Detection** — Automatically identifies your project type (Go, Node.js, Python, Rust, and more) from marker files in your working directory.
- **Contextual Greetings** — Time-aware startup messages that include your current project context.
- **Command Logging** — Persists shell command history with exit codes, working directories, and timestamps to a local SQLite database.
- **Offline-First** — Core features work without any AI service. Cloud AI and local LLMs are optional enhancement layers.

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.26+
- GCC (required by [go-sqlite3](https://github.com/mattn/go-sqlite3))

### Install

```bash
git clone https://github.com/haidaralif-24/Nova-CLI-Companion.git
cd Nova-CLI-Companion
go build -o nova .
```

### Usage

```bash
# Contextual greeting with project detection
nova greet

# Log a command to the local database
nova log "your command here"
```

## Planned Features

- **Shell Integration** — ZSH `preexec`/`precmd` hooks for automatic command observation
- **Command Correction** — Fuzzy matching to catch typos like `git pul` → `git pull`
- **Error Analysis** — Deterministic pattern matching on common failures (e.g., "not a git repository")
- **Session Recap** — Summarize what you worked on in your last terminal session
- **Chat Mode** — Ask questions with full context injection (project type, recent commands, git status)
- **AI Greetings** — Personalized startup messages powered by cloud or local LLMs
- **Multi-Provider AI** — OpenAI, Anthropic, Ollama, and more with automatic fallback

## Architecture

```
┌─────────────────────┐
│       ZSH Shell     │
└──────────┬──────────┘
           │ preexec / precmd
           ▼
┌─────────────────────┐
│      Nova Core      │
└─────┬─────┬─────┬───┘
      │     │     │
 Context  Memory  Suggestions
 Engine   Engine   Engine
      │
      ▼
 AI Provider (optional)
```

Nova is built in Go using [Cobra](https://github.com/spf13/cobra) for CLI commands and [SQLite](https://github.com/mattn/go-sqlite3) for persistent storage.
