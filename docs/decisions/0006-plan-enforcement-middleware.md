---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Enforce subscription plan limits via Go middleware

## Context and Problem Statement

The app has three subscription plans (Basic/Standard/Pro) with different feature access. Enforcement must happen server-side — client-side checks are not security boundaries. Each feature (AI replies, advanced analytics, multiple locations, etc.) needs declarative, type-safe enforcement that cannot be accidentally bypassed.

## Decision Drivers

- Server-side enforcement only (not client-side)
- Declarative — easy to add new limits per feature
- Single place to define the plan-feature matrix
- Type-safe in Go

## Considered Options

- Go middleware with feature flags (chosen)
- Per-handler inline checks (rejected)
- API gateway / Kong / custom proxy (rejected)

## Decision Outcome

**Chosen: Go middleware** that reads the authenticated User's Account Plan from the request context and checks feature access against a compile-time feature matrix.

### Consequences

- Good, because one place to define what each Plan can access
- Good, because type-safe — `PlanGuard("ai_replies")` is a compile-time string literal
- Good, because middleware can be applied per-route or per-route-group via chi
- Bad, because the feature matrix is hardcoded (acceptable for MVP; later could move to DB)

## Implementation Plan

- **Affected paths**: `backend/internal/interfaces/middleware/plan.go`, `backend/cmd/server/main.go` (route setup)
- **Dependencies**: `github.com/go-chi/chi/v5` (middleware chaining)
- **Patterns to follow**: `r.With(middleware.PlanGuard("ai_replies")).Post(...)` applied per route
- **Patterns to avoid**: Inline `if account.Plan != "pro"` checks in handlers
- **Configuration**: Feature matrix defined as `map[string]map[string]bool` in plan.go

### Verification

- [ ] Basic user receives HTTP 403 on `POST /api/locations/:id/reviews/:rid/generate`
- [ ] Standard user can access AI replies but not multi-location creation
- [ ] Pro user can create up to 5 locations (Basic: 1, Standard: 1)
- [ ] Frontend sees structured error `{"type":"/errors/403","title":"Plan Limit","detail":"..."}`
