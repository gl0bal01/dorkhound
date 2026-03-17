# dorkhound — OSINT Missing Person Finder

Fast Google dork URL generator for TraceLab CTF competitions. Single binary, zero dependencies, cross-platform.

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
| Flag | Description |
|------|-------------|
| `--open` | Open all URLs in default browser |
| `--dashboard` | Serve local web dashboard |
| `--export` | Format: `discord`, `json`, `csv`, `clipboard` |
| `--output` | Write to file instead of stdout |

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

MIT
