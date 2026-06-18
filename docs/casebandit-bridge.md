# CaseBandit Bridge — JSON Wire Format (v1)

`dorkhound --export casebandit -o case.json` writes a JSON document that
CaseBandit (https://casebandit.com) imports as a fully-populated case:
identifiers become Entities, each surviving dork becomes a Capture, and the
result is auditable inside CaseBandit's evidence workspace.

This file is the canonical wire-format contract between the two projects.
Bumping `schema_version` is required for any incompatible change.

## Top-level shape

```json
{
  "schema_version": "dorkhound-casebandit-v1",
  "generator": {
    "tool": "dorkhound",
    "version": "1.5.0",
    "generated_at": "2026-05-27T10:15:30Z"
  },
  "case": { "...": "see Case" },
  "entities": [ { "...": "see Entity" } ],
  "captures":  [ { "...": "see Capture" } ]
}
```

- `schema_version` MUST be `dorkhound-casebandit-v1` for this revision.
- `generated_at` is RFC 3339 / ISO 8601 in UTC.
- **Naming convention:** the envelope (`schema_version`, `generated_at`,
  `generator.generated_at`) uses snake_case; everything inside `case` /
  `entities` / `captures` uses camelCase to mirror CaseBandit's
  `shared/types.ts` interfaces verbatim. Consumers should treat the
  envelope as a dorkhound-specific wrapper and the inner objects as
  CaseBandit-native.
- `entities` and `captures` are flat arrays. `case.entities` /
  `case.captures` are intentionally NOT pre-populated to keep the document
  small; CaseBandit's importer attaches them to the case by `caseId`.

## Case

Mirrors `shared/types.ts` `Case` minus fields that only make sense after the
case lives in CaseBandit (`createdAt`/`updatedAt` server timestamps, sync
markers, audit-ledger fields). Importer fills those in.

```json
{
  "id": "dh-<stable-hash>",
  "name": "Jane Doe",
  "description": "TraceLabs CTF 2026-Q2 — last seen Seattle",
  "tags": ["dorkhound", "tracelabs", "missing-person"],
  "status": "active",
  "notes": "## Scope\n...\n## Objectives\n...",
  "chainOfCustody": ""
}
```

- `id` is `dh-` + 24 hex chars of `sha256(name | dob | location)`. Stable
  across reruns of the same case file so reimport merges instead of dupes.
- `status` is always `active` on first import.

## Entity

One Entity per identifier on the case file. Plus an optional `person`
entity for the case subject.

```json
{
  "id": "dh-ent-<stable-hash>",
  "label": "jdoe42",
  "type": "username",
  "notes": "From dorkhound case file (alias of \"Jane Doe\")",
  "source": "dorkhound:case-file",
  "tags": ["dorkhound", "imported"],
  "captureIds": [],
  "status": "unconfirmed",
  "important": false
}
```

Mapping table:

| Case field        | Entity.type   | Notes |
|-------------------|---------------|-------|
| `name`            | `person`      | one entity per case |
| `emails[]`        | `email`       | lower-cased |
| `phones[]`        | `phone`       | E.164 if recognizable, else verbatim |
| `usernames[]`     | `username`    | leading `@` stripped |
| `aliases[]`       | `username`    | tagged `alias` |
| `associates[]`    | `person`      | tagged `associate` |
| `photo_url`       | `url`         | tagged `photo` |
| `location`        | `location`    | only when present |

`important` MAY be set true when the source dork was preflight-classified
`alive` AND priority ≥ 3. Otherwise false.

`captureIds` is populated by the writer where the entity appears in a
generated dork (e.g. an email Entity links to all captures that contain
that email in the query).

## Capture

One Capture per dork in the output set (post-filter, post-preflight,
post-noise-filter). Search-engine-wrapped dorks become Captures whose
`url` is the rendered engine URL (Google/Bing/DuckDuckGo/Yandex). Direct-
URL dorks use the direct URL as-is.

```json
{
  "id": "dh-cap-<stable-hash>",
  "timestamp": "2026-05-27T10:15:30Z",
  "url": "https://www.google.com/search?q=%22Jane+Doe%22+site%3Alinkedin.com",
  "title": "linkedin: Exact name search",
  "type": "page",
  "tags": ["dorkhound", "social", "category:social", "region:global", "engine:google"],
  "content": {
    "text": "Exact name search — site:linkedin.com \"Jane Doe\""
  }
}
```

- `id` is `dh-cap-` + 24 hex chars of `sha256(category | label | url)`.
  Matches the dashboard row ID so dorkhound→CaseBandit re-imports are
  idempotent for unchanged dorks.
- `timestamp` is the generation time, NOT a probe time.
- `tags` ALWAYS include `dorkhound`, the source category, region, and engine.
- `content.text` carries the human-readable dork label.
- `content.screenshot`/`html` are intentionally absent — dorkhound generates
  search plans, not captured pages. CaseBandit + the browser extension
  attach those when the operator visits the URL.

Preflight signals when present:

- `status` = `"green"` for `alive`, `"red"` for `dead`, `"yellow"` for
  `blocked`, omitted for `unchecked` and for search-engine dorks.

## Importer expectations (CaseBandit side)

The CaseBandit importer (`POST /api/import/dorkhound`) MUST:

1. Reject documents whose `schema_version` is not in its allowlist.
2. Reuse `case.id` if it already exists, otherwise create the case.
3. Reuse Entity/Capture IDs to make reruns idempotent.
4. Apply free-tier limits (1 case for free, unlimited for pro/teams).
5. Refuse to import while the user has an unresolved sync conflict on the
   target case.

## Versioning

Schema changes follow a strict additive rule for the v1 lifetime:

- New optional fields are allowed without a version bump.
- Renamed, removed, or semantically-changed fields require a new
  `schema_version` slug (`dorkhound-casebandit-v2`).
- CaseBandit's importer SHOULD accept multiple versions concurrently while
  the rolling-window deprecation lasts (≥ 90 days).
