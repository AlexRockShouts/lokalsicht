---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Use Go backend with Next.js frontend as separate services

## Context and Problem Statement

We are building a SaaS tool for local businesses. The backend must handle Google Business Profile API integration, AI-powered review responses, analytics, and notifications. The frontend must provide SSR/SSG for SEO, multi-language support (DE/FR/IT/RM/EN), and a modern component library.

The lead developer has stronger Go skills than TypeScript/Node.js for backend work.

## Decision Drivers

- Developer proficiency in Go
- SSR/SSG for landing pages and SEO
- Type safety across frontend/backend boundary
- Simple deployment for MVP (max 2 services)

## Considered Options

- Go backend + Next.js frontend as separate services (chosen)
- Next.js monolith (rejected)
- Vite SPA frontend (rejected)

## Decision Outcome

**Chosen: Go `chi` backend + Next.js 14 App Router frontend**, communicating via REST/JSON with JWT auth forwarded through a Next.js API proxy.

### Consequences

- Good, because Go provides type safety, fast compilation, and single-binary deployment
- Good, because Next.js provides SSR/SSG out of the box with file-based routing
- Bad, because two codebases require two CI pipelines and two deployments
- Bad, because type safety across the boundary requires OpenAPI generation (ADR-0009)
- Bad, because the API proxy adds one network hop per request

## Implementation Plan

- **Affected paths**: `backend/` (Go), `frontend/` (Next.js)
- **Dependencies**: `github.com/go-chi/chi/v5`, `next@14`, `next-auth@5`
- **Patterns to follow**: REST/JSON, JWT in `Authorization: Bearer` header, Next.js API proxy pattern
- **Patterns to avoid**: Direct browser-to-Go calls (CORS complexity), GraphQL, gRPC
- **Configuration**: `BACKEND_URL` env var in Next.js, `PORT` env var in Go

### Verification

- [ ] Go API returns valid JSON from `GET /api/me`
- [ ] Next.js proxy forwards `Authorization: Bearer <JWT>` header correctly to Go
- [ ] `make dev` starts both backend (:5174) and frontend (:3000) concurrently
