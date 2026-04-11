# dorkhound — OSINT Google Dork URL Generator

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gl0bal01/dorkhound)](https://go.dev/)
[![CI](https://github.com/gl0bal01/dorkhound/actions/workflows/ci.yml/badge.svg)](https://github.com/gl0bal01/dorkhound/actions/workflows/ci.yml)
[![OSINT Tool](https://img.shields.io/badge/OSINT-Tool-blue)](https://github.com/gl0bal01/dorkhound)
[![TraceLab CTF](https://img.shields.io/badge/TraceLab-CTF-red)](https://www.tracelabs.org/)
[![Google Dorks](https://img.shields.io/badge/Google-Dorks-green)](https://github.com/gl0bal01/dorkhound)
[![Go Report Card](https://goreportcard.com/badge/github.com/gl0bal01/dorkhound)](https://goreportcard.com/report/github.com/gl0bal01/dorkhound)
[![Go Reference](https://pkg.go.dev/badge/github.com/gl0bal01/dorkhound.svg)](https://pkg.go.dev/github.com/gl0bal01/dorkhound)

Google dork URL generator for missing person investigations and TraceLab CTF competitions. Generates search query URLs using advanced Google operators — the user reviews results and determines relevance. Single binary, zero runtime dependencies, cross-platform.

> **Responsible use:** This tool is intended for authorized OSINT investigations, CTF competitions, and educational use only. Always comply with applicable laws and platform terms of service.

## Requirements

- Go 1.25+ (for building from source)

## Install

Download a binary from [Releases], or build from source:

```bash
go install github.com/gl0bal01/dorkhound/cmd/dorkhound@latest
```

Or clone and build:
```bash
git clone https://github.com/gl0bal01/dorkhound.git
cd dorkhound
make build
```

## Quick Start

```bash
# Generate dork URLs for a person
dorkhound -n "John Doe"

# Open all results in browser immediately
dorkhound -n "John Doe" --open

# Include US and Canadian people databases
dorkhound -n "John Doe" --region us,ca --open

# Only social media, export for Discord
dorkhound -n "John Doe" --category social --export discord

# Full case file with dashboard
dorkhound --case case.yaml --dashboard

# Interactive mode
dorkhound -i
```

## Flags

### Input
| Flag | Short | Description |
|------|-------|-------------|
| `--name` | `-n` | Full name ("First Last") |
| `--location` | `-l` | Last known location |
| `--age` | | Approximate age |
| `--dob` | | Date of birth |
| `--aka` | | Aliases, comma-separated |
| `--associates` | | Known associates, comma-separated |
| `--description` | | Physical description |
| `--case` | | Path to YAML/JSON case file |

### Output
| Flag | Short | Description |
|------|-------|-------------|
| `--open` | | Open all URLs in default browser |
| `--dashboard` | | Serve local web dashboard |
| `--export` | | Format: `discord`, `json`, `csv`, `clipboard` |
| `--output` | `-o` | Write to file instead of stdout |

### Filters
| Flag | Default | Description |
|------|---------|-------------|
| `--category` | `all` | `social`, `records`, `financial`, `location`, `forums`, `people-db` |
| `--region` | `global` | Country codes: `us`, `ca`, `uk`, `au`, `ru`, `fr`, `de`, `at`, `nl`, or `all` |
| `--engine` | `google` | `google`, `bing`, `duckduckgo`, `yandex` |
| `--delay` | `100` | Milliseconds between opening tabs |

## Case File

```yaml
name: "John Doe"
aliases: ["JD", "Johnny"]
dob: "1990-01-15"
age: 34
location: "Seattle, WA"
description: "Red hair, tattoo on left arm"
associates: ["Jane Smith", "Bob Johnson"]
region: "us,ca"
categories: ["social", "records"]
engine: "google"
```

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

### CSV

Columns: `label`, `category`, `region`, `priority`, `query`, `url`

```bash
dorkhound -n "John Doe" --region us --export csv -o results.csv
```

### JSON

Includes full case metadata and results with region info. Empty fields are omitted.

```bash
dorkhound -n "John Doe" --export json -o results.json
```

### Discord

Markdown-formatted output grouped by category, ready to paste into Discord.

```bash
dorkhound -n "John Doe" --export discord
```

### Clipboard

Same as Discord format, copied directly to system clipboard.

```bash
dorkhound -n "John Doe" --export clipboard
```

## Categories

| Category | Description |
|----------|-------------|
| `social` | Facebook, LinkedIn, Instagram, Twitter/X, TikTok, YouTube, GitHub, etc. |
| `records` | Court records, property, obituaries, education, resumes, contact cards |
| `financial` | PayPal, Venmo, bank/loan mentions, contact spreadsheets |
| `location` | Google Maps reviews, travel/booking, relocation mentions |
| `forums` | Reddit, Quora, forum profiles, associate cross-references |
| `people-db` | Direct lookups on Spokeo, Whitepages, TruePeopleSearch, etc. |

## Nuclei integration (username enumeration)

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

Nuclei results appear under category `nuclei`. If the binary is not found, a warning is printed and the rest of the dorks are returned normally.

## Preflight dead-link filter

```bash
dorkhound --name "Jane Doe" --emails "jane@example.com" --preflight
```

Probes each dork whose `Query` is a direct HTTP(S) URL (HIBP, Wayback, IntelX, TrueCaller, people-db lookups, nuclei `matched-at` URLs) with an HTTP HEAD request and drops any that return 4xx/5xx or fail to connect. Search-engine-wrapped dorks (the majority) are passed through untouched — `HEAD google.com/search?q=...` always returns 200 regardless of result count, so probing them provides no useful signal.

| Flag | Default | Description |
|------|---------|-------------|
| `--preflight` | `false` | Enable preflight probing |
| `--preflight-timeout` | `5s` | Per-request timeout |
| `--preflight-concurrency` | `8` | Max parallel probes |
| `--preflight-rate` | `250ms` | Per-worker cooldown between requests |

## Reverse image search

Pass `--photo-url` with a publicly accessible image URL to generate direct reverse-image-search dorks across nine engines:

```bash
dorkhound --name "Jane Doe" --photo-url "https://example.com/jane.jpg" --category image
```

Emits direct-URL dorks for Google Lens, Google Images, Yandex (best for faces), TinEye, Bing Visual Search, PimEyes, SauceNAO, IQDB, and KarmaDecay. The operator opens each link — no credentials required. Two name-based photo-platform searches (Flickr/500px/SmugMug and Imgur) are also emitted when `--name` is set.

## TraceLabs submission format

Generate a Markdown checklist formatted for TraceLabs CTF flag submission:

```bash
dorkhound --name "Jane Doe" --emails "jane@example.com" --usernames "jdoe42" \
  --photo-url "https://example.com/jane.jpg" --export tracelabs
```

Produces a `# TraceLabs Submission — <name>` document with a case header, leads grouped by category (image → username → email → phone → social → … → unknowns), and a submission workflow checklist. Pipe to a file with `--output report.md`.

## Shell Completions

```bash
# Bash
source <(dorkhound completion bash)

# Zsh
dorkhound completion zsh > "${fpath[1]}/_dorkhound"

# Fish
dorkhound completion fish | source
```

## License

[MIT](LICENSE)

[Releases]: https://github.com/gl0bal01/dorkhound/releases
