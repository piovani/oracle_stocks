---
name: update-docs
description: Update this project's documentation (README.md and CLAUDE.md) to match the current state of the codebase. Use when packages, commands, dependencies, env vars, migrations, endpoints, or conventions have changed and the docs may be stale.
---

# Update project documentation

Refresh the two docs so they match the **actual** codebase, not what it used to be:

- `README.md` — human-facing: how to set up, run, and use the project.
- `CLAUDE.md` — AI-agent context: architecture, data flow, conventions, gotchas.

Never invent content. Everything in the docs must be verifiable in the repo right now.

## 1. Gather current state

Read the sources of truth before touching either file:

- `go.mod` — module path, Go version, direct dependencies.
- `Makefile` — every target (these are the canonical commands).
- `cmd/*/main.go` — one binary per dir; note flags and what each does.
- `internal/**` — package layout and responsibilities.
- `internal/database/migrations/` — schema and how migrations are applied.
- `.env.example` + `internal/config/config.go` — environment variables and defaults.
- Route registrations in `cmd/api/main.go` — the real list of HTTP endpoints.

Prefer reading files over guessing. If a command or endpoint isn't in the code, it doesn't go in the docs.

## 2. Update CLAUDE.md (agent context)

Keep it concise and current. Sections to maintain:

- **Project** — one-paragraph purpose, module path, Go version.
- **Architecture** — a `cmd/` + `internal/` tree with a one-line role per package.
- **Data pipeline** — how data flows end to end (e.g. provider → mapper → service → repository).
- **Key dependencies** — table mapping package → role.
- **Database & migrations** — connection wrapper, where the schema lives, how migrations run.
- **Configuration** — env vars.
- **Common commands** — mirror the Makefile.
- **Conventions** — the rules a new agent must follow (see below).
- **Current API endpoints** — only routes actually registered.

## 3. Update README.md (human-facing)

Sections to maintain: short intro, Stack, Requirements, Getting started (the real setup sequence in order), Environment variables, any feature-specific usage (e.g. backfill examples), Make targets table, API table, Docker, and a Project structure tree.

Make sure the "Getting started" steps run in the correct order and actually work (e.g. migrations must run before the app needs the schema).

## 4. Preserve project conventions

Carry the existing conventions forward verbatim unless the code contradicts them. Current ones to respect when describing or showing examples:

- GORM struct tags carry **only** the column name (`gorm:"column:x"`); types/constraints/indexes live in migration SQL.
- Nullable DB columns map to pointer fields in Go (`nilIfZero` helper).
- Functional options for constructors; thin commands with orchestration in a service.
- Errors wrapped with `fmt.Errorf("context: %w", err)`; logging via `log/slog`.
- No comments unless the WHY is non-obvious. Keep docs lean — no filler.

## 5. Verify

After editing, sanity-check the docs against reality:

- Every Make target listed exists in the `Makefile`.
- Every endpoint listed is registered in `cmd/api`.
- The structure tree matches `cmd/` and `internal/`.
- `go build ./...` still passes (you didn't change code, but confirm nothing references removed doc-only assumptions).

Report what changed in each file in 2–3 bullets.
