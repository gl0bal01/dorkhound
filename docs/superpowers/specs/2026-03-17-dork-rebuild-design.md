# Dork — OSINT Missing Person Finder Rebuild

## Overview

Rebuild the existing Python OSINT missing person finder as a Go CLI tool optimized for TraceLab CTF competitions. The tool generates Google dork URLs from case information and gets them into a browser as fast as possible. No scraping, no HTTP requests to search engines — just smart URL generation and multiple output modes.

Unlike the original Python tool, this rebuild does not make HTTP requests to search engines. It generates and opens search URLs directly, avoiding rate-limiting and blocking issues.

## Context

- **Users:** TraceLab CTF participants working in teams, sharing results via Discord
- **Platforms:** Windows, Linux, macOS — must be cross-platform
- **Key constraint:** Speed. In a CTF, seconds matter. Zero setup, zero dependencies, single binary.

## Why Go

- Single binary, no runtime or dependencies needed on target machine
- Cross-compiles to Windows/Linux/macOS trivially
- Great stdlib for CLI, HTTP server, browser opening, embed
- Cobra gives shell completions for free

## CLI Interface

Binary name: `dork`

### Input Flags

| Flag | Short | Required | Description |
|------|-------|----------|-------------|
| `--name` | `-n` | Yes (unless `--case`) | Full name as "First Last" |
| `--location` | `-l` | No | Last known location |
| `--age` | | No | Approximate age |
| `--dob` | | No | Date of birth |
| `--aka` | | No | Aliases/nicknames, comma-separated |
| `--associates` | | No | Known associates, comma-separated |
| `--description` | | No | Physical description |
| `--case` | | No | Path to YAML/JSON case file. If both `--case` and `--name` are provided, CLI flags override case file values. |

### Output Flags

| Flag | Description |
|------|-------------|
| `--open` | Open all URLs in default browser immediately |
| `--dashboard` | Serve local web dashboard |
| `--export discord` | Discord-formatted output to stdout |
| `--export json` | JSON output to stdout |
| `--export csv` | CSV output to stdout |
| `--export clipboard` | Discord format copied to clipboard |
| `--output <file>` | Write export to file instead of stdout |

### Filter Flags

| Flag | Values | Default | Description |
|------|--------|---------|-------------|
| `--category` | `social`, `records`, `financial`, `location`, `forums`, `people-db`, `all` | `all` | Filter dork categories, comma-separated |
| `--region` | Country codes, `all`, `global` | `global` | Comma-separated. `global` = only non-region-specific dorks. `all` = include every region. Country codes: `us`, `ca`, `uk`, `au`, `ru`, `fr`, `de`, `at`, `nl`. |
| `--engine` | `google`, `bing`, `duckduckgo`, `yandex` | `google` | Search engine for URL generation |
| `--delay` | milliseconds | `100` | Delay between opening browser tabs |
| `--version` | | | Print version and exit |

### Modes

| Flag | Description |
|------|-------------|
| `-i` | Interactive mode with tab completion |
| `completion <shell>` | Generate shell completion script (bash/zsh/fish/powershell) |

### Example Usage

```bash
# Minimum: just a name, open tabs
dork -n "John Doe" --open

# Full case info from file, dashboard mode
dork --case case.yaml --dashboard

# Quick social media search, export for Discord
dork -n "John Doe" -l "Seattle" --category social --export discord

# US + Canada people databases
dork -n "John Doe" --region us,ca --open

# UK-specific search
dork -n "John Doe" --region uk --open

# Everything, all countries
dork -n "John Doe" --region all --open

# Interactive mode
dork -i
```

## Case File Format

YAML (also supports JSON):

```yaml
name: "John Doe"
aliases:
  - "JD"
  - "Johnny"
dob: "1990-01-15"
age: 34
location: "Seattle, WA"
description: "Red hair, tattoo on left arm"
associates:
  - "Jane Smith"
  - "Bob Johnson"
# Optional output preferences (CLI flags override these)
region: "us,ca"
categories:
  - "social"
  - "records"
engine: "google"
```

## Dork Generation Engine

### Core Data Structure

```go
type Dork struct {
    Query    string   // the search query string
    Category string   // social, records, financial, location, forums, people-db
    Region   string   // country code (us, ca, uk, au, ru, fr, de, at, nl) or "global" for non-region-specific
    Priority int      // 1-3, higher = opens first
    Label    string   // human-readable label (e.g., "Facebook profile")
}
```

### Categories

**Social** — Global: Facebook, LinkedIn, Instagram, Twitter/X, Reddit, GitHub, Medium, TikTok, YouTube. Region-specific: VK.com and OK.ru (`ru`), Copainsdavant (`fr`), StayFriends (`de`), Hyves archives (`nl`).

**Records** — Court records, property records, marriage records, PDF public records, obituaries, news articles.

**Financial** — PayPal, Venmo, CashApp mentions, banking, loan, mortgage, credit references. Higher priority dorks only.

**Location** — City/state mentions, Google Maps reviews, travel/booking mentions, Foursquare/Yelp reviews.

**Forums** — Reddit, Quora, forum profiles, username pattern searches, community mentions.

**People-DB** — Region-specific people search databases:

| Region | Sites |
|--------|-------|
| `us` | Spokeo, Whitepages, TruePeopleSearch, FastPeopleSearch, BeenVerified |
| `ca` | Canada411, CanadaPeopleSearch, WhitePages.ca |
| `uk` | 192.com, FindMyPast, BT Phone Book, UKElectoralRoll |
| `au` | WhitePages.com.au, PeopleFinder.com.au, ReverseAustralia |
| `ru` | VK.com, OK.ru, Yandex People, NumBuster |
| `fr` | PagesBlanches.fr, Copainsdavant.linternaute.com, AnnuaireMairie |
| `de` | DasTelefonbuch.de, Telefonbuch.de, StayFriends.de |
| `at` | Herold.at, DasTelefonbuch.at |
| `nl` | DeTelefoongids.nl, WhitePages.nl, Numberway.nl |

### Name Parsing

`--name` accepts a full name string. Parsing strategy: split on the last space. `"John Doe"` → first="John", last="Doe". `"Mary Jane Watson"` → first="Mary Jane", last="Watson". Both components are available to dork generators that need them (e.g., people-db URL patterns like `spokeo.com/John-Doe`). The full unsplit name is also available for quoted-phrase dorks.

### Smart Query Building

The engine uses all available case info to narrow queries. More info = more specific queries.

- **Name only:** `"John Doe" site:facebook.com`
- **Name + location:** `"John Doe" "Seattle" site:facebook.com`
- **Name + DOB/age:** `"John Doe" "1990" OR "age 34"`
- **Name + alias:** `"John Doe" OR "JD" site:facebook.com`
- **Associates cross-ref:** `"John Doe" "Jane Smith"`
- **Description:** `"John Doe" "red hair" OR "tattoo"`

Each dork definition function receives the full case struct and builds the best query it can from available fields.

### Search Engine URL Templates

| Engine | URL Template | Notes |
|--------|-------------|-------|
| Google | `https://www.google.com/search?q={query}` | Full dork operator support (`site:`, `inurl:`, `filetype:`, etc.) |
| Bing | `https://www.bing.com/search?q={query}` | Supports `site:` and `filetype:`, limited `inurl:` |
| DuckDuckGo | `https://duckduckgo.com/?q={query}` | Supports `site:`, limited other operators |
| Yandex | `https://yandex.com/search/?text={query}` | Supports `site:`, `mime:` instead of `filetype:` |

Dork queries are passed verbatim to all engines (same as the original tool). Non-Google engines may not support all operators, but the queries still produce useful results as the operators are treated as search terms.

### Priority Assignments

- **Priority 3** — Direct profile URLs: social media `site:` searches, people-db direct links. Most likely to yield a hit.
- **Priority 2** — Record and document searches: court records, PDFs, property, news. Useful but noisier.
- **Priority 1** — Broad/speculative searches: forums, description-based, financial, travel. Long shots.

### Filtering

- `--category social` only generates social dorks
- `--region global` (default) includes only dorks tagged `region: global` (non-region-specific)
- `--region us,ca` includes global dorks PLUS dorks tagged `us` or `ca`
- `--region all` includes every dork regardless of region tag
- Results ordered by priority (highest first) within each category

## Output Modes

### Tab Blaster (`--open`)

Opens URLs in default browser using `os/exec` (platform-appropriate: `xdg-open`, `open`, `start`). Default 100ms delay between tabs (configurable via `--delay`). Ordered by priority, grouped by category so related tabs are adjacent.

### Dashboard (`--dashboard`)

Starts a local HTTP server on a random available port. Serves a single embedded HTML page:

- Case info summary at top
- Links grouped by category
- Each link has: checkbox (reviewed/useful/dead end), copy button
- "Copy all in category" button
- "Export to Discord" button generates formatted summary of checked items
- State lives in memory only — no persistence, no database

**Dashboard architecture:** State is purely client-side JavaScript (in-memory objects). The Go server serves only the initial HTML page with dork data embedded as a JSON blob in a `<script>` tag. No API endpoints needed beyond serving the page. The "Export to Discord" button copies formatted text to clipboard via the Clipboard API.

HTML/CSS/JS embedded in binary via Go's `embed` package. No external dependencies. Single file.

### Discord Export (`--export discord`)

Formatted markdown to stdout:

```
## OSINT Results: John Doe
**Location:** Seattle, WA | **Age:** ~34

### Social Media (12 links)
- Facebook: <url>
- LinkedIn: <url>
...

### Public Records (8 links)
- Court records: <url>
...
```

### JSON Export (`--export json`)

```json
{
  "case": { "name": "John Doe", "location": "Seattle, WA" },
  "results": [
    { "label": "Facebook profile", "url": "...", "category": "social", "priority": 3 }
  ]
}
```

### CSV Export (`--export csv`)

`label,category,priority,url` — one row per dork.

### Clipboard (`--export clipboard`)

Same as Discord format, copied to system clipboard via `atotto/clipboard`.

All exports write to stdout by default. `--output file.txt` redirects to file.

## Interactive Mode (`-i`)

Drops into an interactive prompt (using `charmbracelet/bubbletea` or `charmbracelet/huh` for cross-platform support including Windows). Tab-completes:

- Flag names and values
- Category names
- Region values
- Engine names

Allows building a search incrementally — add fields one by one, then fire. Useful when receiving case info piecemeal.

## Shell Completions

`dork completion bash/zsh/fish/powershell` outputs a shell completion script. Provided by Cobra for free. Completes all flags, enum values (`--category`, `--region`, `--engine`, `--export`).

## Project Structure

```
dork/
├── cmd/
│   └── dork/
│       └── main.go              # entry point, CLI setup via Cobra
├── internal/
│   ├── case/
│   │   └── case.go              # Case struct, YAML/JSON parsing
│   ├── dork/
│   │   ├── engine.go            # dork generation orchestration
│   │   ├── dorks_social.go      # social media dork definitions
│   │   ├── dorks_records.go     # public records dorks
│   │   ├── dorks_financial.go   # financial dorks
│   │   ├── dorks_location.go    # location dorks
│   │   ├── dorks_forums.go      # forum/community dorks
│   │   └── dorks_peopledb.go    # people search database URLs
│   ├── output/
│   │   ├── browser.go           # tab blaster
│   │   ├── dashboard.go         # local web server + embedded HTML
│   │   ├── discord.go           # Discord-formatted export
│   │   ├── json.go              # JSON export
│   │   ├── csv.go               # CSV export
│   │   └── clipboard.go         # clipboard copy
│   └── interactive/
│       └── prompt.go            # interactive mode with tab completion
├── web/
│   └── dashboard.html           # HTML template (embedded in binary)
├── go.mod
├── go.sum
├── Makefile                     # build, cross-compile, install targets
└── README.md
```

## Dependencies

- `github.com/spf13/cobra` — CLI framework, shell completions
- `github.com/charmbracelet/bubbletea` + `github.com/charmbracelet/huh` — interactive mode (cross-platform, good Windows support)
- `gopkg.in/yaml.v3` — YAML case file parsing
- `github.com/atotto/clipboard` — clipboard support
- Standard library for everything else (net/http, embed, encoding/json, encoding/csv, os/exec)

## Build & Distribution

```makefile
# Build for current platform
make build

# Cross-compile all platforms
make release  # produces dork-linux-amd64, dork-darwin-amd64, dork-darwin-arm64, dork-windows-amd64.exe

# Install locally
make install  # copies to $GOPATH/bin
```

Version info injected at build time via `-ldflags "-X main.version=..."`. `dork --version` prints version string.

## Default Behavior

If no output flag is provided (`--open`, `--dashboard`, `--export`), the tool prints all generated URLs to stdout grouped by category — one URL per line with its label. This lets you pipe output or just see what would be opened. At least one output flag is recommended but not required.

## CLI Override Rules

CLI flags always override case file values. If `--case case.yaml` provides `name: "John Doe"` and `--name "Jane Doe"` is also passed, "Jane Doe" wins. Same for all input and preference fields. This applies to both person data (`--name`, `--location`, etc.) and output preferences (`--region`, `--category`, `--engine`).

## Output Flag Interactions

Output flags are combinable: `--open` can be used alongside `--export` and `--output`. Example: `dork -n "John Doe" --open --export discord` opens tabs AND prints Discord format. `--dashboard` is standalone — when used, it serves the dashboard and ignores `--open` and `--export` (the dashboard provides its own export functionality).

## Error Handling

Errors print to stderr. Exit codes:
- **0** — success
- **1** — input error (bad flags, malformed case file, missing required fields)
- **2** — runtime error (clipboard unavailable, browser launch failed)

Partial success is acceptable: if 3 of 50 browser tabs fail to open, log warnings to stderr and continue. Non-critical failures (clipboard, single tab) should not abort the run.

## Testing

Table-driven tests for the core logic:
- **Name parser** — splitting edge cases (single name, three+ part names, hyphenated names)
- **Dork engine** — given a Case struct, verify correct queries generated per category
- **Filtering** — region and category filtering produce correct subsets
- **Output formatters** — Discord, JSON, CSV output matches expected format

No tests for browser opening or dashboard serving — those are manual/integration.

## What This Tool Does NOT Do

- No HTTP requests to search engines (no scraping, no rate-limiting, no blocks)
- No link validation or liveness checking
- No database or persistent storage
- No authentication or API keys
- No proxy support (unnecessary since we don't make external requests)
