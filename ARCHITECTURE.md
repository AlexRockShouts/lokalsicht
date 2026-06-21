# Lokalsicht — Architecture Overview

## Stack

| Layer | Technology | Hosting |
|-------|------------|---------|
| Frontend | Next.js 14 App Router, Tailwind CSS, shadcn/ui, next-intl | Vercel (Hobby, CHF 0) |
| Backend | Go 1.22+, chi router, GORM, golang-jwt, go-openai, go-playground/validator | Railway (512MB, $5) |
| Database | PostgreSQL 16 | Railway Managed ($5) |
| AI | DeepSeek V3 (OpenAI-compatible) | Pay-per-use (~CHF 1/mo) |
| Email | Resend | Free (100/day) |
| Cron | GitHub Actions Scheduled | Free |

**Total MVP cost: ~CHF 10/month**

## Project Layout

```
lokalsicht/
├── CONTEXT-MAP.md              # Bounded context map (DDD)
├── ARCHITECTURE.md             # This file
├── Makefile                    # Root orchestration
├── docker-compose.yml          # Local PostgreSQL
├── .env.example                # Environment variable template
│
├── backend/                    # Go service
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/             # Pure domain: aggregates, value objects, ports
│   │   │   ├── account/        #   Account, User, Plan, Reseller
│   │   │   ├── location/       #   Location, GoogleProfile
│   │   │   ├── review/         #   Review, Reply, AI Suggestion, events
│   │   │   ├── insight/        #   InsightSnapshot
│   │   │   └── notification/   #   Notification, Preference
│   │   ├── application/        # Use cases: orchestrates domain objects
│   │   ├── infrastructure/     # Adapters: implements domain ports
│   │   │   ├── persistence/    #   GORM repositories
│   │   │   ├── gbp/            #   Google Business Profile API
│   │   │   ├── ai/             #   DeepSeek client
│   │   │   ├── email/          #   Resend client
│   │   │   └── stripe/         #   Stripe client
│   │   └── interfaces/         # HTTP handlers, middleware
│   │       ├── http/
│   │       └── middleware/
│   ├── docs/swagger.json       # Generated OpenAPI spec
│   └── go.mod
│
├── frontend/                   # Next.js App Router
│   ├── src/
│   │   ├── app/
│   │   │   ├── [locale]/       # i18n routing (de, fr, it, rm, en)
│   │   │   │   ├── (auth)/     # Login, Onboarding
│   │   │   │   ├── (dashboard)/# Protected: Dashboard, Reviews, Analytics, Settings
│   │   │   │   └── (marketing)/# Landing page
│   │   │   └── api/
│   │   │       ├── auth/[...nextauth]/  # NextAuth.js
│   │   │       └── [...path]/          # API proxy → Go backend
│   │   ├── components/ui/     # shadcn/ui components
│   │   ├── lib/               # API client, utils, validation
│   │   └── i18n/              # Translation JSON files
│   └── package.json
│
├── docs/decisions/             # Architecture Decision Records
│   ├── 0001-go-backend-nextjs-frontend.md
│   ├── 0002-railway-hosting.md
│   ├── 0003-deepseek-ai-provider.md
│   ├── 0004-combined-google-oauth.md
│   ├── 0005-github-actions-cron.md
│   ├── 0006-plan-enforcement-middleware.md
│   ├── 0007-google-independent-features.md
│   ├── 0008-reseller-gtm.md
│   ├── 0009-openapi-typesafety.md
│   ├── 0010-five-language-support.md
│   └── 0011-hexagonal-architecture.md
│
└── .github/workflows/
    ├── backend.yml
    ├── frontend.yml
    └── cron.yml
```

## Key Architecture Decisions

See `docs/decisions/` for full decision records.

| ADR | Decision |
|-----|----------|
| 0001 | Go backend + Next.js frontend as separate services |
| 0002 | Railway for Go + PostgreSQL hosting |
| 0003 | DeepSeek V3 as AI provider |
| 0004 | Combined Google OAuth flow (login + GBP) |
| 0005 | GitHub Actions as cron scheduler |
| 0006 | Plan enforcement via Go middleware |
| 0007 | Google-independent features in Phase 1.5 |
| 0008 | Reseller + Direct GTM |
| 0009 | OpenAPI-driven type safety |
| 0010 | 5 languages, 1 primary per user |
| 0011 | Hexagonal (ports & adapters) architecture |

## Bounded Contexts

See `CONTEXT-MAP.md` for the full context map.

| Context | Domain Package | Glossary |
|---------|---------------|----------|
| Account Management | `domain/account/` | Account, User, Plan, Reseller |
| Location Management | `domain/location/` | Standort, Google-Profil |
| Review Management | `domain/review/` | Bewertung, Antwort, KI-Vorschlag |
| Analytics | `domain/insight/` | Einsicht, Sync |
| Notification | `domain/notification/` | Benachrichtigung, Präferenz |
| Billing | (infrastructure + domain/account) | Stripe-native |

## Data Flow

```
[Browser]
    ↓ Next.js (Vercel)
    ↓ JWT Bearer + API Proxy
    ↓
[Go Backend (Railway)]
    ├── interfaces/http/    ← handles HTTP
    ├── application/        ← use cases
    ├── domain/             ← business rules
    └── infrastructure/     ← adapters
        ├── persistence/    → PostgreSQL (Railway)
        ├── gbp/            → Google Business Profile API
        ├── ai/             → DeepSeek API
        ├── email/          → Resend
        └── stripe/         → Stripe
```

## Auth Flow

```
1. User clicks "Login"
2. NextAuth.js → Google OAuth (openid profile email business.manage)
3. Google redirects to NextAuth callback
4. NextAuth issues JWT (NEXTAUTH_SECRET shared with Go)
5. Frontend calls /api/[...path] → Next.js API proxy
6. Proxy forwards JWT → Go verifies HMAC signature
7. Go middleware looks up/creates User in DB
8. User context attached to request
```

## C4 Container Diagram (MVP)

```
   ┌─────────────────────────────────────────────┐
   │ Web Application (Next.js + Vercel)           │
   │ - SSR/SSG landing pages                      │
   │ - Protected dashboard SPA                    │
   │ - next-intl i18n routing                     │
   │ - NextAuth.js Google OAuth                   │
   └──────────────────┬──────────────────────────┘
                      │ REST/JSON (JWT)
   ┌──────────────────▼──────────────────────────┐
   │ API Server (Go + Railway)                    │
   │ - chi router, GORM, JWT verification         │
   │ - Domain-driven hexagonal architecture       │
   │ - Plan enforcement middleware                │
   │ - GBP sync, AI replies, email notifications  │
   └──┬──────────┬──────────┬──────────┬─────────┘
      │          │          │          │
      ▼          ▼          ▼          ▼
   ┌──────┐ ┌──────┐ ┌────────┐ ┌──────────┐
   │ PG   │ │ GBP  │ │ DeepSk │ │ Resend   │
   │ (RWY)│ │ API  │ │ API    │ │ (Email)  │
   └──────┘ └──────┘ └────────┘ └──────────┘
```

## Roadmap

| Phase | Month | Focus |
|-------|-------|-------|
| Core MVP | 1–4 | Scaffolding, Auth, GBP, Dashboard, Reviews, AI, Analytics |
| Phase 1.5 | 4–5 | Review links, Schema.org, Template library, Multi-platform |
| Stripe | 3–4 | Payment integration (silent, activated on first upgrade) |
| Beta | 5–6 | Alpha feedback, i18n completion, bugfixing |
| Soft Launch | 7 | Open registration, early-adopter pricing |
| Phase 2 | 7+ | Apple Maps, Social Media, Competitor monitoring |

## Dependencies

### Go (backend)
- `github.com/go-chi/chi/v5` — HTTP router
- `github.com/go-chi/cors` — CORS middleware
- `gorm.io/gorm` + `gorm.io/driver/postgres` — ORM
- `github.com/golang-jwt/jwt/v5` — JWT verification
- `github.com/sashabaranov/go-openai` — DeepSeek (OpenAI-compatible)
- `golang.org/x/oauth2` + `google.golang.org/api` — Google OAuth + GBP API
- `github.com/resend/resend-go/v2` — Email
- `github.com/go-playground/validator/v10` — Input validation
- `github.com/swaggo/swag` — OpenAPI generation
- `golang.org/x/time/rate` — Rate limiting

### Next.js (frontend)
- `next@14`, `react@18`, `typescript@5`
- `next-auth@5` — Authentication
- `next-intl@3` — i18n
- `@tanstack/react-query` — Server state
- `tailwindcss@3` + `shadcn/ui` — Styling
- `vitest` + `@testing-library/react` — Testing

## Environment Variables

Shared between both services:
- `NEXTAUTH_SECRET` — JWT signing key (Go + Next.js)
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` — Google OAuth

Go only:
- `DATABASE_URL` — Railway PostgreSQL
- `PORT` — default 5174
- `DEEPSEEK_API_KEY`
- `RESEND_API_KEY`
- `ENCRYPTION_KEY` — AES-256 key for token encryption
- `CRON_API_KEY` — internal API key for cron endpoints
- `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`

Next.js only:
- `NEXTAUTH_URL` — canonical URL
- `BACKEND_URL` — Go backend URL (internal Railway URL)
- `SENTRY_DSN` — error monitoring
