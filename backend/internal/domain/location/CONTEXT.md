# Location Management

Handles physical business locations and their Google Business Profile connections.

## Language

**Standort** (Location):
A physical business location that belongs to an Account. Has a name, address, phone, website, opening hours, and description. In German, the canonical term is _Standort_.
_Avoid_: Place, Store, Site, Filiale

**Google-Profil** (Google Profile):
The link between a Standort and its Google Business Profile entry. Contains the Google Place ID and the encrypted OAuth refresh token used to access the GBP API.
_Avoid_: Connection, Link, Binding

**Profil** (Profile):
The public-facing data of a Standort: opening hours, phone number, website, description, photos. This is what gets synced to/pulled from Google Business Profile.
_Avoid_: Settings, Config, Metadata

**Verknüpfen** (Connect):
The OAuth flow that grants Lokalsicht access to a user's Google Business Profile account. The user approves the `business.manage` scope.
_Avoid_: Link, Authorize, Register
