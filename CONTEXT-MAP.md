# Context Map

## Contexts

- [Account Management](./backend/internal/domain/account/CONTEXT.md) — users, accounts, plans, resellers
- [Location Management](./backend/internal/domain/location/CONTEXT.md) — locations, Google Business Profile connection
- [Review Management](./backend/internal/domain/review/CONTEXT.md) — reviews, AI-generated replies, multi-platform inbox
- [Analytics](./backend/internal/domain/insight/CONTEXT.md) — insights, performance snapshots
- [Notification](./backend/internal/domain/notification/CONTEXT.md) — email notifications, preferences
- Billing — Stripe subscriptions, plans, checkout (no separate glossary — terms are Stripe-native)
- Optimization — review links, QR codes, Schema.org, templates (Phase 1.5)

## Relationships

- **Account → Location**: Account owns Locations. Each Location references its AccountID.
- **Account → Billing**: Account upgrades Plan. Billing emits webhook events consumed by Account.
- **Location → Review**: Location owns Reviews. Each Review references its LocationID.
- **Location → Analytics**: Analytics syncs InsightSnapshots per LocationID.
- **Review → Notification**: Review context emits `ReviewReceived` domain events; Notification consumes them.
- **Billing → Account**: Plan changes flow from Stripe webhook → Billing service → Account aggregate.

## Shared Types

- `AccountID`, `LocationID`, `ReviewID` are shared across context boundaries.
- `Plan` (Basic | Standard | Pro) is shared between Account and Billing.
