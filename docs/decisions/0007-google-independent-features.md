---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Prioritize Google-independent features (Phase 1.5)

## Context and Problem Statement

The product depends critically on the Google Business Profile API for its core value proposition. Google has a history of deprecating APIs (Google+, Reader, My Maps) and regularly changes GBP API behavior. A Google API shutdown or access restriction would render the core product non-functional without alternative features.

## Decision Drivers

- Mitigate existential Google API dependency
- Build features that work without GBP API
- Keep investment low (2–3 weeks total for Phase 1.5)
- Align with the product name "Lokalsicht" (local visibility, not just Google)

## Considered Options

- Phase 1.5 with 5 Google-independent features (chosen)
- Google-only core, defer alternatives to Phase 2 (rejected)
- Full multi-platform from day 1 (rejected — too expensive)

## Decision Outcome

**Chosen: Phase 1.5** immediately following the Core MVP (Month 4–5), prioritized by customer value and independence:

| Priority | Feature | Value | Google-Free |
|----------|---------|-------|-------------|
| 1 | Review-Link-Generator + QR-Code | High | Yes |
| 2 | Schema.org LocalBusiness Generator | Medium | Yes |
| 3 | Multi-Platform Review Inbox (Facebook, Apple) | Very High | Partial |
| 4 | Response/Post Template Library | Medium | Yes |
| 5 | WhatsApp Review Request | High (CH) | Yes |

### Consequences

- Good, because the product remains useful even if GBP API access is restricted
- Good, because features 1, 2, 4, and 5 need ZERO Google API calls
- Good, because review link generator + QR code is a self-service feature — no integration needed
- Bad, because Phase 1.5 delays Phase 2 (Apple Maps, Social Media) by 4–6 weeks

## Implementation Plan

- **Affected paths**: `backend/internal/application/optimization/`, `frontend/src/app/[locale]/(dashboard)/links/`
- **Dependencies**: `go-qrcode` (QR code generation), Schema.org JSON-LD structs (pure Go)
- **Patterns to follow**: Each optimization feature is a standalone Go service in `application/optimization/`
- **Patterns to avoid**: Mixing optimization features into existing location/review handlers

### Verification

- [ ] Review-Link-Generator produces a valid Google review URL + downloadable QR code
- [ ] Schema.org generator outputs valid JSON-LD for LocalBusiness
- [ ] Multi-platform inbox shows placeholder for Facebook reviews (prior to Facebook API integration)
- [ ] All Phase 1.5 features work without active GBP API connection
