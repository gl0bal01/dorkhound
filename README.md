# dorkhound — OSINT Google Dork URL Generator

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gl0bal01/dorkhound)](https://go.dev/)
[![CI](https://github.com/gl0bal01/dorkhound/actions/workflows/ci.yml/badge.svg)](https://github.com/gl0bal01/dorkhound/actions/workflows/ci.yml)
[![OSINT Tool](https://img.shields.io/badge/OSINT-Tool-blue)](https://github.com/gl0bal01/dorkhound)
[![TraceLab CTF](https://img.shields.io/badge/TraceLab-CTF-red)](https://www.tracelabs.org/)
[![Google Dorks](https://img.shields.io/badge/Google-Dorks-green)](https://github.com/gl0bal01/dorkhound)
[![Go Report Card](https://goreportcard.com/badge/github.com/gl0bal01/dorkhound)](https://goreportcard.com/report/github.com/gl0bal01/dorkhound)
[![Go Reference](https://pkg.go.dev/badge/github.com/gl0bal01/dorkhound.svg)](https://pkg.go.dev/github.com/gl0bal01/dorkhound)

dorkhound generates 150+ OSINT dork queries and direct-profile URLs across 20 categories to accelerate missing-person investigations and TraceLabs CTF competitions. You provide a name (and optionally emails, usernames, a photo, or a full case file); dorkhound produces a ranked list of search URLs you open and triage. Single binary, zero runtime dependencies, cross-platform.

> **Responsible use:** This tool is intended for authorized OSINT investigations, CTF competitions, and educational use only. Always comply with applicable laws and platform terms of service.

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

## TraceLabs CTF Workflow

1. **Receive the brief** — organizers provide name, approximate age, last-known location, and a photo URL.
2. **Populate a case file** — copy `examples/tracelabs.yaml`, fill in the four fields, save as `case.yaml`.
3. **Launch the dashboard** — `dorkhound --case case.yaml --dashboard`. A local web UI opens in your browser.
4. **Triage** — use the filter bar to focus on high-signal categories (image, username, direct-profile). The dashboard persists notes and evidence fields in localStorage between sessions.
5. **Pivot** — as you discover new handles, emails, or phone numbers, add them to `case.yaml` and re-run. Each new identifier significantly expands the dork corpus.
6. **Export** — `dorkhound --case case.yaml --export tracelabs -o submission.md` produces a checklist formatted for TraceLabs flag submission.

The dashboard's "Open batch" button opens URLs in rate-limited groups to avoid search-engine CAPTCHAs (configurable with `--batch` and `--batch-pause`).

## Features

- 150+ dorks across 20 categories: social, records, financial, location, forums, people-db, email, phone, username, cache, documents, dating, marketplace, image, gravatar, github, academic, direct-profile, nuclei, and more
- Direct-URL dorks for 20+ platforms including Telegram, Keybase, Twitter/X, Mastodon, Bluesky, GitHub, Gravatar, and others — no search engine required
- Reverse image search across Google Lens, Yandex, TinEye, Bing Visual Search, PimEyes, SauceNAO, IQDB, and KarmaDecay — pass `--photo-url` to activate
- Nuclei v2 integration: username enumeration across 600+ sites via `-tags osint`
- Preflight HTTP checker: probes direct-URL dorks with HEAD requests and drops dead links before operator triage
- Interactive dashboard with localStorage state persistence, per-result notes and evidence fields, filter bar, rate-limited "Open batch", and keyboard shortcuts
- TraceLabs submission format export
- Export formats: discord, json, csv, clipboard, tracelabs
- Region filters for US, CA, UK, AU, RU, FR, DE, AT, NL
- Interactive guided-prompt mode (`-i`)

## Flags

Run `dorkhound --help` for the full flag reference. The most commonly used flags:

| Flag | Description |
|------|-------------|
| `--name` / `-n` | Full name ("First Last") |
| `--case` | Path to YAML case file |
| `--dashboard` | Serve local web dashboard |
| `--export` | Output format: discord, json, csv, clipboard, tracelabs |
| `--category` | Category filter (default: all) |
| `--region` | Region filter (default: global) |
| `--open` | Open all URLs in default browser |
| `--preflight` | Drop dead direct-URL dorks before output |
| `--nuclei` | Run nuclei OSINT templates against usernames |
| `--stats` | Print dork count breakdown and exit |

## Dashboard Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `/` | Focus filter bar |
| `Esc` | Clear filter |
| `j` / `k` | Next / previous result |
| `x` | Mark current result reviewed |
| `d` | Open current URL in browser |
| `n` | Edit note for current result |
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

Probes each dork whose query is a direct HTTP(S) URL (people-db lookups, Wayback, HIBP, nuclei matched-at URLs, etc.) with an HTTP HEAD request and drops any that return 4xx/5xx or fail to connect. Search-engine-wrapped dorks are passed through untouched — probing `google.com/search?q=...` always returns 200 regardless of result count.

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

| Format | Description |
|--------|-------------|
| `discord` | Markdown grouped by category, ready to paste into Discord |
| `json` | Full case metadata + results with region info |
| `csv` | Columns: label, category, region, priority, query, url |
| `clipboard` | Discord format copied to system clipboard |
| `tracelabs` | TraceLabs CTF submission checklist |

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
