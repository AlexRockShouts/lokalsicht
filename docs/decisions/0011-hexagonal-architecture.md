---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Use hexagonal (ports & adapters) architecture for the Go backend

## Context and Problem Statement

The Go backend handles multiple bounded contexts (Account, Location, Review, Analytics, Notification, Billing). Traditional layered architecture (handler → service → repository) tends to couple business logic to infrastructure concerns. Domain-Driven Design with hexagonal architecture ensures the domain layer has zero external dependencies.

## Decision Drivers

- Domain logic must be testable without database, HTTP, or external APIs
- Multiple bounded contexts with clear boundaries
- Infrastructure adapters (PostgreSQL, Google API, DeepSeek, Resend, Stripe) should be swappable
- Consistency with the domain model defined in `CONTEXT-MAP.md`

## Considered Options

- Hexagonal architecture: `domain/` → `application/` → `infrastructure/` → `interfaces/` (chosen)
- Flat layered: `handler/` → `service/` → `repository/` (rejected)
- Clean Architecture (Use Cases as central concern) (rejected — hex is simpler for Go)

## Decision Outcome

**Chosen: Hexagonal architecture** with four layers:

```
internal/
├── domain/       # Pure Go — aggregates, value objects, port interfaces. Zero imports from other layers.
├── application/  # Use cases — orchestrates domain objects. Depends on domain (not infrastructure).
├── infrastructure/# Adapters — implements domain ports (GORM repos, GBP client, AI client, email, Stripe).
└── interfaces/   # HTTP handlers, middleware. Depends on application layer.
```

### Consequences

- Good, because domain is fully testable with mocks
- Good, because infrastructure adapters are swappable (e.g., DeepSeek → Grok via config)
- Good, because each bounded context has its own directory in `domain/`
- Bad, because more files and directories than a flat structure
- Bad, because developers unfamiliar with DDD may find the indirection confusing

## Implementation Plan

- **Affected paths**: `backend/internal/domain/`, `backend/internal/application/`, `backend/internal/infrastructure/`, `backend/internal/interfaces/`
- **Dependencies**: None — architecture convention, not a library
- **Patterns to follow**: Each domain package defines its `Repository` interface (port). Infrastructure packages implement them. Application services accept interfaces, not concrete types.
- **Patterns to avoid**: Domain imports from `gorm`, `net/http`, `github.com/sashabaranov/go-openai`, or any external package. Domain imports only Go stdlib.

### Verification

- [ ] No `gorm` import in any file under `internal/domain/`
- [ ] No `net/http` import in any file under `internal/domain/`
- [ ] Each domain package has a `ports.go` file defining repository interfaces
- [ ] Application services are testable with mock repositories
