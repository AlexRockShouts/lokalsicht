---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Generate OpenAPI spec from Go and TypeScript types from OpenAPI

## Context and Problem Statement

Go backend and Next.js frontend communicate via REST/JSON. Without automated type synchronization, API type definitions will drift between the two codebases. Manual synchronization is error-prone and wastes developer time on avoidable bugs.

## Decision Drivers

- Go structs are the single source of truth for API types
- Frontend needs typed API client (fetch wrapper with TypeScript types)
- Must run in CI to prevent drift
- Low setup overhead (minutes, not days)

## Considered Options

- swaggo (Go) → OpenAPI JSON → openapi-typescript (TypeScript) (chosen)
- Manual type maintenance (rejected)
- gRPC/Protobuf (rejected)

## Decision Outcome

**Chosen: `swaggo/swag`** generates OpenAPI spec from Go handler comments. **`openapi-typescript`** generates TypeScript types from that spec. Runs as a CI check: if types are stale, the frontend build fails.

### Consequences

- Good, because single source of truth (Go structs + swaggo comments)
- Good, because CI enforces type consistency — no drift
- Good, because zero runtime overhead — types are compile-time only
- Bad, because swaggo comments add maintenance burden to Go handlers
- Bad, because complex nested types may need manual adjustment

## Implementation Plan

- **Affected paths**: `backend/docs/swagger.json`, `frontend/src/api/types.ts`, `.github/workflows/frontend.yml`
- **Dependencies**: `github.com/swaggo/swag/cmd/swag`, `openapi-typescript` (npm devDependency)
- **Patterns to follow**: swaggo annotations on every handler function, `openapi-typescript` run as npm script
- **Patterns to avoid**: Manual type definitions that mirror Go structs
- **Configuration**: `swag init -g cmd/server/main.go -o docs/`, `openapi-typescript docs/swagger.json -o frontend/src/api/types.ts`

### Verification

- [ ] `swag init` generates `docs/swagger.json` from handler comments
- [ ] `openapi-typescript` generates `frontend/src/api/types.ts`
- [ ] Frontend CI build fails if types are stale (`npm run typegen:check` step)
- [ ] API client functions in `frontend/src/api/` use generated types
