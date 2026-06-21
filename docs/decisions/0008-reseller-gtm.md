---
status: accepted
date: 2026-06-21
decision-makers: Alexander Pina
---

# Use reseller channel (web designers) + direct word-of-mouth as go-to-market

## Context and Problem Statement

Selling directly to small local businesses is time-consuming and expensive. Web designers already have trusted relationships with these businesses. A reseller model leverages existing trust and reduces customer acquisition cost.

## Decision Drivers

- Low customer acquisition cost
- Trusted existing relationships (web designer → business)
- Swiss market: personal relationships matter more than online ads
- Two channels for redundancy: reseller + direct

## Considered Options

- Reseller + Direct (chosen)
- Direct sales only (rejected)
- Google Ads only (rejected)

## Decision Outcome

**Chosen: Dual channel** — reseller program (20% commission + White-Label option) for web designers, plus direct customer sign-ups through word-of-mouth.

### Consequences

- Good, because web designers bring pre-qualified leads with existing trust
- Good, because 20% commission aligns designer incentives with our revenue
- Good, because White-Label option makes it the designer's product, not ours
- Neutral, because commission tracking requires Stripe Connect or manual accounting
- Bad, because reseller onboarding requires documentation and support bandwidth

## Implementation Plan

- **Affected paths**: `backend/internal/domain/account/account.go` (ResellerID field), `backend/internal/interfaces/http/billing_handler.go`
- **Dependencies**: Stripe Connect for automatic commission (Phase 2), manual tracking in Phase 1
- **Patterns to follow**: ResellerID on Account model, White-Label branding via `Account.Settings` JSONB
- **Patterns to avoid**: Hardcoded reseller lists, fixed commission rates (keep configurable)
- **Configuration**: `RESELLER_COMMISSION_PCT=20` env var

### Verification

- [ ] New Account can be created with a ResellerID
- [ ] Dashboard shows Reseller branding when White-Label is active
- [ ] Commission is calculated on each subscription payment
