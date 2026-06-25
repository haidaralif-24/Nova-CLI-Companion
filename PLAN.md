# Nova — AI CLI Companion

## Overview

Nova is an AI-powered terminal companion that lives alongside the user's shell.

Unlike traditional terminal chatbots, Nova's primary purpose is not conversation. Its purpose is to understand the user's workflow, observe terminal activity, provide intelligent assistance, remember project context, and improve the command-line experience.

Nova should remain useful even when no AI service is available.

---

# Goals

## Primary Goals

* Make terminal usage more pleasant and productive
* Detect and help recover from common command mistakes
* Understand project context automatically
* Provide session memory across terminal launches
* Act as a lightweight terminal companion with personality

## Secondary Goals

* Explain command failures
* Provide terminal chat capabilities
* Generate contextual greetings
* Track development progress over time

## Non-Goals

* Replace existing shells
* Replace Claude Code, Codex, or Copilot
* Fully automate development workflows
* Become an autonomous coding agent

---

# Core Philosophy

Nova follows:

```text
Terminal Companion
    >
Terminal Chatbot
```

The terminal integration is the product.

AI is an enhancement layer.

If every AI provider disappeared tomorrow, Nova should still be useful.

---

# High-Level Architecture

```text
┌─────────────────────┐
│       ZSH Shell     │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│    Shell Hooks      │
│ preexec / precmd    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│      Nova Core      │
└─────┬─────┬─────┬───┘
      │     │     │
      ▼     ▼     ▼

 Context  Memory  Suggestions
 Engine   Engine   Engine
      │
      ▼
 AI Provider Layer
      │
      ▼
 OpenAI / Anthropic / Local LLM
```

---

# Major Components

## 1. Shell Integration

### Responsibilities

* Observe commands before execution
* Observe command results
* Capture exit codes
* Detect working directory changes
* Trigger greetings

### ZSH Hooks

Before execution:

```zsh
preexec()
```

After execution:

```zsh
precmd()
```

Shell startup:

```zsh
nova greet
```

---

## 2. Context Engine

Responsible for understanding the current environment.

### Detect Project Types

Rust

```text
Cargo.toml
```

Astro

```text
astro.config.mjs
```

Node.js

```text
package.json
```

Python

```text
pyproject.toml
requirements.txt
```

Go

```text
go.mod
```

### Context Data

```json
{
  "project_type": "rust",
  "project_name": "compiler",
  "git_branch": "parser",
  "modified_files": 4
}
```

---

## 3. Memory Engine

Persistent storage using SQLite.

### Purpose

Remember:

* sessions
* projects
* commands
* failures
* notes

### Tables

Sessions

```sql
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    started_at DATETIME,
    ended_at DATETIME
);
```

Commands

```sql
CREATE TABLE commands (
    id INTEGER PRIMARY KEY,
    command TEXT,
    cwd TEXT,
    exit_code INTEGER,
    timestamp DATETIME
);
```

Projects

```sql
CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    name TEXT,
    path TEXT,
    project_type TEXT
);
```

---

## 4. Suggestion Engine

Provides deterministic assistance.

### Command Correction

Input

```bash
git pul
```

Output

```text
Did you mean:

git pull
```

### Techniques

* Levenshtein distance
* Fuzzy matching
* Command dictionaries

### Sources

* Git commands
* Cargo commands
* NPM commands
* Docker commands
* User command history

---

## 5. Error Analysis Engine

Maps known failures to suggestions.

### Example

Input

```text
git: not a git repository
```

Output

```text
You are not inside a Git repository.

Possible fixes:
- git init
- move to repository root
```

### Architecture

```text
stderr
    ↓
Pattern Matcher
    ↓
Known Rule
    ↓
Suggestion
```

AI should only be used if no rule matches.

---

## 6. Greeting System

Responsible for startup messages.

### Online Mode

Uses AI provider.

Input Context

```json
{
  "time": "20:00",
  "project": "compiler",
  "branch": "parser",
  "last_session_duration": "2h"
}
```

Output

```text
Welcome back.

Looks like you're still working on the parser branch.
```

### Offline Mode

Template-based.

Example

```text
Good evening.

Rust project detected.

Ready to continue?
```

---

## 7. AI Layer

Optional enhancement layer.

### Responsibilities

* Greetings
* Error explanations
* Chat mode
* Session summaries

### Responsibilities It Does NOT Own

* Command correction
* Project detection
* Memory storage
* Shell integration

Those should work without AI.

---

## 8. Chat Mode

Command

```bash
nova chat
```

Example

```text
You:
How does cargo workspaces work?

Nova:
...
```

### Context Injection

Include:

* current directory
* project type
* recent commands
* git status

to improve relevance.

---

# Offline Strategy

Nova uses a three-layer fallback strategy.

```text
Cloud AI
    ↓
Local LLM
    ↓
Templates & Rules
```

## Layer 1

Cloud AI

Examples:

* OpenAI
* Anthropic

## Layer 2

Local Models

Examples:

* Ollama
* MLX

## Layer 3

Offline Core

Always available.

Includes:

* greetings
* memory
* suggestions
* project detection
* error analysis

---

# CLI Commands

## Greeting

```bash
nova greet
```

---

## Chat

```bash
nova chat
```

---

## Explain Last Error

```bash
nova explain
```

---

## Recap Session

```bash
nova recap
```

---

## Project Summary

```bash
nova project
```

---

## Configuration

```bash
nova config
```

---

# Proposed Tech Stack

## Language

Go

### Reasons

* Single binary distribution
* Fast startup
* Excellent CLI ecosystem
* Easy shell integration
* Cross-platform support

---

## Database

SQLite

---

## CLI Framework

Cobra

---

## TUI

Bubble Tea

Lip Gloss

---

## AI Providers

Any OpenAI compatible providers

Local Models

* Ollama
* MLX

---


# Future Features

## Project Timeline

```bash
nova timeline
```

Shows development history.

---

## Smart Project Notes

```bash
nova note "parser completed"
```

Stored automatically.

---

## Daily Recap

```bash
nova daily
```

Generates summary of work completed.

---

## Team Mode

Shared project memory between developers.

---

# Success Criteria

Nova is successful when:

* It feels useful without AI.
* It understands project context automatically.
* It catches common mistakes.
* It remembers previous work.
* It provides assistance without interrupting workflow.

The terminal should feel like it has a helpful companion, not a chatbot attached to a shell.
