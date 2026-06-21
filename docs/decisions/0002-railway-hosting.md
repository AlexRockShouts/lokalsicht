---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Host Go backend and PostgreSQL on Railway

## Context and Problem Statement

We need a hosting platform that runs our Go backend and a managed PostgreSQL database. Requirements: no cold start (sleep), single provider, zero-to-minimal cost for MVP, git-push deployment.

Alternatives evaluated: Fly.io (sleep problem, no managed PostgreSQL), Supabase + Fly.io (two providers, sleep), Hetzner VPS (self-managed DevOps overhead), Render (sleep problem).

## Decision Drivers

- No cold start / sleep on idle — the app must always be responsive
- Single provider for backend + database
- Git-push deployment
- MVP cost ≤ CHF 10/month

## Considered Options

- Railway $10/month (chosen)
- Fly.io $0 + Supabase $0 (rejected)
- Hetzner CX22 ~CHF 4/month (rejected)
- Render $7/month (rejected)

## Decision Outcome

**Chosen: Railway** with one Go service (512MB RAM) + one managed PostgreSQL (1GB).

### Consequences

- Good, because no sleep, responsive on first request
- Good, because git-push deploy like Fly.io
- Good, because managed PostgreSQL with automatic backups
- Bad, because $10/month is the most expensive option evaluated (but still trivial for a SaaS)
- Neutral, because Railway can scale by adjusting the service slider (more RAM/CPU)

## Implementation Plan

- **Affected paths**: `backend/railway.toml` or `backend/Dockerfile`
- **Dependencies**: None (Railway detects Go via `go.mod`)
- **Patterns to follow**: Railway auto-deploy on git push to main branch
- **Patterns to avoid**: `fly.toml`, Supabase connection strings
- **Configuration**: `DATABASE_URL` from Railway PostgreSQL, `PORT=5174` (internal)

### Verification

- [ ] Go binary deployed and reachable via Railway URL
- [ ] PostgreSQL accessible from Go binary (ping on startup)
- [ ] Railway auto-deploys on `git push` to main
