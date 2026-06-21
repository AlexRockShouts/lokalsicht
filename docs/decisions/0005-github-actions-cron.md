---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Use GitHub Actions Scheduled as external cron scheduler

## Context and Problem Statement

We need recurring jobs: (1) daily insight sync from Google Business Profile Performance API, and (2) 15-minute review checks for new customer reviews. Railway has no native cron scheduler. In-process Go goroutines were considered but rejected due to instability during deployments and lack of high-availability guarantees.

## Decision Drivers

- Zero additional cost
- Reliable execution (no missed runs)
- Runs outside the Go process (survives deployments and restarts)
- Simple to monitor (CI logs + notification on failure)

## Considered Options

- GitHub Actions Scheduled workflow (chosen)
- In-process Go goroutine with `time.Ticker` (rejected)
- cron-job.org external service (rejected)
- Vercel Cron Jobs (rejected)

## Decision Outcome

**Chosen: GitHub Actions Scheduled** with two cron triggers calling internal Go API endpoints.

### Consequences

- Good, because CHF 0 and runs in the same GitHub account as CI/CD
- Good, because external to the Go process — survives deployments
- Good, because manual trigger available via `workflow_dispatch`
- Bad, because GitHub Actions schedules can be delayed by up to 1 hour during peak (acceptable for insight sync, slightly annoying for review check)

## Implementation Plan

- **Affected paths**: `.github/workflows/cron.yml`, `backend/internal/interfaces/http/cron_handler.go`
- **Dependencies**: None (uses GitHub built-in schedule + curl)
- **Patterns to follow**: Internal cron endpoints protected by `CRON_API_KEY` Bearer token, DB unique constraints prevent duplicate data
- **Patterns to avoid**: In-process goroutine timers for production jobs
- **Configuration**: `BACKEND_URL`, `CRON_API_KEY` as GitHub Secrets

### Verification

- [ ] `POST /api/internal/cron/check-reviews` runs every 15 minutes
- [ ] `POST /api/internal/cron/sync-insights` runs daily at 02:00 UTC
- [ ] `CRON_API_KEY` is required for both endpoints (401 without it)
- [ ] Insight snapshot duplicate prevention via `@@unique([locationId, date])` constraint
- [ ] GitHub Actions log shows success/failure for each run
