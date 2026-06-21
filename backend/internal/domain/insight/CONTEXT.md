# Analytics

Handles performance metrics for Google Business Profiles.

## Language

**Einsicht** (Insight):
A daily snapshot of a Standort's Google Business Profile performance metrics, fetched from the Google Business Profile Performance API. Contains views, clicks, calls, and direction requests for that day.
_Avoid_: Statistic, Metric, Data Point, Snapshot

**Aufruf** (View):
How many times the Google Business Profile was shown in Google Search or Google Maps results on a given day.
_Avoid_: Impression, Display, Hit

**Klick** (Click):
How many times someone clicked on the profile (e.g., to visit the website, call, or get directions) on a given day.
_Avoid_: Interaction, Engagement

**Anruf** (Call):
How many times someone tapped the call button from the Google Business Profile on a given day.
_Avoid_: Phone Call, Dial

**Routenabfrage** (Direction Request):
How many times someone requested directions to the Standort from Google Maps on a given day.
_Avoid_: Navigation, Map Request

**Sync** (Synchronisation):
The scheduled process that pulls new Einsichten from the Google Performance API and stores them as Insight Snapshots in the database. Runs daily via GitHub Actions cron.
_Avoid_: Fetch, Import, Update
