# Account Management

Handles user accounts, authentication, subscription plans, and reseller relationships.

## Language

**Account**:
A business entity that owns one or more Locations. Has a subscription Plan.
_Avoid_: Tenant, Organisation, Company

**User**:
A person with login credentials belonging to an Account. Has a Role (Owner | Member).
_Avoid_: Teammitglied, Benutzer, Employee

**Plan**:
The subscription tier of an Account: Basic (free, 1 location, no AI), Standard (CHF 69, AI replies, analytics), Pro (CHF 109, 5 locations, all features).
_Avoid_: Subscription, Tarif, Paket

**Owner**:
The User who created the Account. Initially the only User. Has full permissions.
_Avoid_: Admin, Creator

**Reseller**:
A web designer or agency that sells Lokalsicht to end customers on behalf of the platform. Receives 20% commission on referred Accounts. May have White-Label branding.
_Avoid_: Partner, Affiliate, Agent

**White-Label**:
A Reseller option where the Reseller's logo and branding replace Lokalsicht's in the dashboard and emails for their referred Accounts.
_Avoid_: Branding, Custom CSS
