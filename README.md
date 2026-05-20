# dorkhound — OSINT Investigation Toolkit for TraceLabs CTFs

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gl0bal01/dorkhound)](https://go.dev/)
[![CI](https://github.com/gl0bal01/dorkhound/actions/workflows/ci.yml/badge.svg)](https://github.com/gl0bal01/dorkhound/actions/workflows/ci.yml)
[![OSINT Tool](https://img.shields.io/badge/OSINT-Tool-blue)](https://github.com/gl0bal01/dorkhound)
[![TraceLab CTF](https://img.shields.io/badge/TraceLab-CTF-red)](https://www.tracelabs.org/)
[![Google Dorks](https://img.shields.io/badge/Google-Dorks-green)](https://github.com/gl0bal01/dorkhound)
[![Go Report Card](https://goreportcard.com/badge/github.com/gl0bal01/dorkhound)](https://goreportcard.com/report/github.com/gl0bal01/dorkhound)
[![Go Reference](https://pkg.go.dev/badge/github.com/gl0bal01/dorkhound.svg)](https://pkg.go.dev/github.com/gl0bal01/dorkhound)

dorkhound generates 340+ OSINT dork queries, direct-profile URLs, and reverse-image-search links across 25 categories to accelerate missing-person investigations and TraceLabs CTF competitions. You provide a name (and optionally emails, phones, usernames, a photo, or a full case file); dorkhound produces a ranked list of search URLs and direct endpoints you open and triage in an interactive dashboard. Single binary, zero runtime dependencies, cross-platform.

dorkhound is an **investigation accelerator**, not an automated identification system. It does not verify identity, prove that a lead belongs to a subject, bypass access controls, or replace human review. Treat every result as an investigative lead that must be corroborated with source context, timestamps, screenshots, and independent evidence.

> **Responsible use:** This tool is intended for authorized OSINT investigations, CTF competitions, and educational use only. Always comply with applicable laws and platform terms of service.

## Status and Scope

dorkhound is suitable for public CLI distribution and local operator workflows. It is not a hosted multi-user service, and it is not designed to store, process, or share cases on a server. Case files, notes, exports, browser history, clipboard contents, and downloaded evidence remain the operator's responsibility.

Use dorkhound when you need to quickly build a structured search plan from known identifiers. Do not use it to harass, stalk, dox, make automated eligibility decisions, or publish unverified allegations.

## Safety Model

- **Local-first:** The dashboard binds to `127.0.0.1` on a random port and an unguessable path. It serves only the local browser session.
- **Private by default:** Dashboard notes use session storage by default. Use "Persist state" only on a trusted machine and browser profile.
- **Explicit exports:** Export files are created with owner-only permissions and are not overwritten unless `--force-output` is passed.
- **Safer paste output:** Markdown exports neutralize Discord mentions and escape user-controlled text. CSV exports protect against spreadsheet formula execution.
- **Network guardrails:** Preflight probes skip search-query dorks and block private, loopback, link-local, and metadata-style targets by default.
- **Supply-chain hygiene:** CI and release workflows pin GitHub Actions by commit SHA, use reduced token permissions, and release checksums are generated for binary artifacts.

## Install

Three options:

```bash
# Option 1: go install
go install github.com/gl0bal01/dorkhound/cmd/dorkhound@latest

# Option 2: download a pre-built binary
# See: https://github.com/gl0bal01/dorkhound/releases

# Option 3: build from source
git clone https://github.com/gl0bal01/dorkhound.git
cd dorkhound
make build
```

Requirements: Go 1.25+ (for options 1 and 3 only).

## Quick Start

```bash
# Simple: name only
dorkhound --name "Jane Doe"

# Rich case from YAML
dorkhound --case examples/full.yaml --dashboard

# TraceLabs submission draft
dorkhound --case case.yaml --export tracelabs -o submission.md
```

If the output file already exists, dorkhound exits instead of overwriting it. Pass `--force-output` only when replacing the file is intentional.

## TraceLabs CTF Workflow

1. **Receive the brief** — organizers provide name, approximate age, last-known location, and a photo URL.
2. **Populate a case file** — copy `examples/tracelabs.yaml`, fill in the four fields, save as `case.yaml`.
3. **Launch the dashboard** — `dorkhound --case case.yaml --dashboard`. A local web UI opens in your browser.
4. **Triage** — use the filter bar to focus on high-signal categories (image, username, direct-profile). The dashboard keeps notes in session storage by default; enable "Persist state" only on a trusted browser.
5. **Pivot** — as you discover new handles, emails, or phone numbers, add them to `case.yaml` and re-run. Each new identifier significantly expands the dork corpus.
6. **Export** — `dorkhound --case case.yaml --export tracelabs -o submission.md` produces a checklist formatted for TraceLabs flag submission.

The dashboard's "Open batch" button opens URLs in rate-limited groups to avoid search-engine CAPTCHAs (configurable with `--batch` and `--batch-pause`).

## Operational Guidance

Good OSINT work is slower than link collection. Use dorkhound to find places to look, then record why each lead is relevant.

Recommended workflow:

1. Start with the minimum known identifiers from an authorized brief.
2. Open high-priority categories first: image, username, direct-profile, email, phone, and official registries.
3. Mark false positives and dead ends aggressively; common names produce noise.
4. Capture source URLs, screenshots, timestamps, and the exact observed facts before drawing conclusions.
5. Add newly corroborated identifiers back into the case file and re-run.
6. Keep case files, exports, screenshots, and browser profiles secured when working with sensitive personal data.

Do not submit or share a lead solely because dorkhound generated it. A generated query is not evidence.

## Limitations

- Search engine results vary by region, account state, personalization, rate limits, and time.
- Direct-profile URLs can exist for unrelated people using the same handle.
- People-search sites may be inaccurate, stale, paywalled, jurisdiction-specific, or legally restricted.
- Reverse-image search results can be weak for cropped, compressed, edited, or AI-generated images.
- Nuclei results depend on the installed binary, templates, network conditions, and platform changes.
- Preflight only checks whether direct URLs respond; it does not validate that a page contains relevant evidence.

## Features

- 340+ dorks across 25 categories: social, records, financial, location, forums, people-db, email, phone, username, cache, documents, dating, marketplace, image, gravatar, github, academic, direct-profile, twitter, reddit, fundraiser, telegram, vehicle, crypto, nuclei
- Twitter/X audit: every Twitter dork emits `site:twitter.com OR site:x.com` plus nitter mirrors, Wayback Machine captures of both domains, the X syndication endpoint, and X advanced-search URLs for login-less reconnaissance
- Reddit audit: `old.reddit.com` alongside `reddit.com`, RSS feeds, `about.json` profile endpoints, and third-party search indexers
- Direct-URL dorks for 20+ platforms (Telegram, Keybase, Twitter/X, Mastodon, Bluesky, GitHub, Reddit old/new, Steam, Twitch, Last.fm, SoundCloud, Medium, Dev.to, Dribbble, Flickr, About.me, Linktree, and more) — no search engine required
- Reverse image search across Google Lens, Yandex (best for faces), TinEye, Bing Visual Search, PimEyes, SauceNAO, IQDB, KarmaDecay — pass `--photo-url` to activate
- GitHub OSINT: commits-by-email search, `.keys` / `.gpg` leak probes, public-events API, gists
- Gravatar lookups: MD5-hashed email → avatar, profile JSON, profile page
- Missing-person registries: NamUs, CharleyProject, DoeNetwork, NCMEC (US), MissingPeople (UK), Interpol notices
- Nuclei v2 integration: username enumeration across 600+ sites via `-tags osint`
- Preflight HTTP checker: HEAD-probes direct-URL dorks and drops dead links before operator triage
- Interactive localhost dashboard with per-session notes by default, optional trusted-browser persistence, per-result evidence fields, filter bar, rate-limited "Open batch", keyboard shortcuts
- TraceLabs submission format export
- Export formats: discord, json, csv, clipboard, tracelabs
- Region filters for US, CA, UK, AU, RU, FR, DE, AT, NL
- Interactive guided-prompt mode (`-i`)
- Dork deduplication: collisions across generators are collapsed, highest-priority wins

## Flags

Run `dorkhound --help` for the full flag reference. The most commonly used flags:

| Flag | Description |
|------|-------------|
| `--name` / `-n` | Full name ("First Last") |
| `--case` | Path to YAML or JSON case file |
| `--emails` | Email addresses, comma-separated |
| `--phones` | Phone numbers, comma-separated |
| `--usernames` | Usernames/handles, comma-separated |
| `--photo-url` | Photo URL for reverse image search |
| `--location` / `-l` | Last known location |
| `--dashboard` | Serve local web dashboard with notes & filters |
| `--export` | Output format: discord, json, csv, clipboard, tracelabs |
| `--force-output` | Overwrite an existing `--output` file |
| `--category` | Category filter (see `--list-categories`) |
| `--region` | Region filter (see `--list-regions`) |
| `--open` | Open all URLs in default browser |
| `--batch` / `--batch-pause` | Rate-limited batch opening (default 10 per 30s to avoid CAPTCHAs) |
| `--delay` | Delay between tabs when `--open` is set (default 2000ms) |
| `--noise-filter` | Append noise-suppression operators (`-site:pinterest.com` etc.) |
| `--preflight` | Drop dead direct-URL dorks before output |
| `--nuclei` | Run nuclei OSINT templates against usernames |
| `--stats` | Print dork count breakdown and exit |
| `-i` / `--interactive` | Guided prompt mode |

## Dashboard Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `/` | Focus filter bar |
| `Esc` | Clear filter |
| `j` / `k` | Next / previous result |
| `x` | Toggle current result between pending and useful |
| `d` | Mark current result as a dead end |
| `n` | Toggle the notes panel for the current result |
| `?` | Show help overlay |

## Nuclei Integration

Install the nuclei binary and OSINT templates:

```bash
go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest
nuclei -update-templates
```

Run dorkhound with nuclei enabled:

```bash
dorkhound --name "Jane Doe" --usernames "jdoe42" --nuclei
# Custom tags and timeout
dorkhound --name "Jane Doe" --usernames "jdoe42" --nuclei --nuclei-tags osint --nuclei-timeout 5m
```

Nuclei probes 600+ sites for each username and appends results under category `nuclei`. If the binary is not found, a warning is printed and the rest of the dorks are returned normally.

## Preflight Dead-Link Filter

Probes each dork whose query is a direct public HTTP(S) URL (people-db lookups, Wayback, HIBP, nuclei matched-at URLs, etc.) with an HTTP HEAD request and drops any that return 4xx/5xx or fail to connect. Search-engine-wrapped dorks are passed through untouched — probing `google.com/search?q=...` always returns 200 regardless of result count.

Preflight blocks private, loopback, link-local, and metadata-style targets by default. This prevents accidental probing of local services or cloud metadata endpoints if a future direct URL source produces an unsafe target.

```bash
dorkhound --name "Jane Doe" --emails "jane@example.com" --preflight
```

| Flag | Default | Description |
|------|---------|-------------|
| `--preflight` | false | Enable preflight probing |
| `--preflight-timeout` | 5s | Per-request timeout |
| `--preflight-concurrency` | 8 | Max parallel probes |
| `--preflight-rate` | 250ms | Per-worker cooldown |

## Reverse Image Search

```bash
dorkhound --name "Jane Doe" --photo-url "https://example.com/jane.jpg" --category image
```

Emits direct-URL dorks for Google Lens, Yandex (best for faces), TinEye, Bing Visual Search, PimEyes, SauceNAO, IQDB, and KarmaDecay. No credentials required — the operator opens each link.

## TraceLabs Submission Format

```bash
dorkhound --name "Jane Doe" --emails "jane@example.com" --usernames "jdoe42" \
  --photo-url "https://example.com/jane.jpg" --export tracelabs -o submission.md
```

Produces a `# TraceLabs Submission — <name>` Markdown document with a case header, leads grouped by category, and a submission workflow checklist.

## Regions

| Code | Sites |
|------|-------|
| `us` | Spokeo, Whitepages, TruePeopleSearch, FastPeopleSearch, BeenVerified |
| `ca` | Canada411, CanadaPeopleSearch, WhitePages.ca |
| `uk` | 192.com, FindMyPast, BT Phone Book, UKElectoralRoll |
| `au` | WhitePages AU, PeopleFinder AU, ReverseAustralia |
| `ru` | VK, OK.ru, Yandex People, NumBuster |
| `fr` | PagesBlanches |
| `de` | DasTelefonbuch, Telefonbuch.de |
| `at` | Herold.at, DasTelefonbuch.at |
| `nl` | DeTelefoongids, WhitePages.nl, Numberway.nl |

## Export Formats

Exports may contain sensitive personal data. Store them in an appropriate case folder, avoid pasting raw exports into public channels, and review all generated content before submission.

| Format | Description |
|--------|-------------|
| `discord` | Markdown grouped by category, ready to paste into Discord |
| `json` | Full case metadata + results with region info |
| `csv` | Columns: label, category, region, priority, query, url |
| `clipboard` | Discord format copied to system clipboard |
| `tracelabs` | TraceLabs CTF submission checklist |

Markdown exports escape user-controlled text and neutralize Discord mentions. CSV exports prefix formula-like cells so spreadsheet software does not execute them as formulas.

## Shell Completion

```bash
# Bash (current session)
source <(dorkhound completion bash)

# Bash (persist)
dorkhound completion bash > /etc/bash_completion.d/dorkhound

# Zsh
dorkhound completion zsh > "${fpath[1]}/_dorkhound"

# Fish
dorkhound completion fish > ~/.config/fish/completions/dorkhound.fish

# PowerShell
dorkhound completion powershell > dorkhound.ps1
```

Or use `make completion-install` to install for your current shell automatically.

## Example Case Files

The `examples/` directory contains ready-to-use templates:

| File | When to use |
|------|-------------|
| `examples/minimal.yaml` | Bare minimum — name only. Simplest possible case. |
| `examples/full.yaml` | Every field documented with inline comments. Copy and replace values. |
| `examples/tracelabs.yaml` | TraceLabs CTF starting point. Fill in four fields from the brief and go. |
| `examples/advanced.yaml` | Targeted investigation with filtered categories when you already have leads. |

## Case File Reference

```yaml
name: "John Doe"
aliases: ["JD", "Johnny"]
dob: "1990-01-15"
age: 34
location: "Seattle, WA"
description: "Red hair, tattoo on left arm"
associates: ["Jane Smith", "Bob Johnson"]
emails: ["john@example.com"]
phones: ["+1-555-000-0000"]
usernames: ["jdoe42"]
photo_url: "https://example.com/john.jpg"
photo_path: "/local/path/to/photo.jpg"
region: "us"
categories: ["social", "records"]
engine: "google"
```

## Contributing

Pull requests welcome. Keep new dork templates focused on OSINT use cases. Run `make lint` before submitting.

## License

[MIT](LICENSE)

[Releases]: https://github.com/gl0bal01/dorkhound/releases
