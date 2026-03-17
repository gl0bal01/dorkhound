# Dork Go Rebuild — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild the OSINT missing person finder as a single Go binary optimized for TraceLab CTF speed.

**Architecture:** CLI parses input (flags or YAML case file) into a `Case` struct. The dork engine generates `[]Dork` from the case. Output handlers consume the dork list to open tabs, serve a dashboard, or export formatted text. Each layer depends only on the one before it.

**Tech Stack:** Go 1.23+ (use whatever is installed; `go mod init` will set the version automatically), Cobra (CLI), gopkg.in/yaml.v3 (case files), charmbracelet/bubbletea+huh (interactive mode), atotto/clipboard, Go embed (dashboard HTML).

**Spec:** `docs/superpowers/specs/2026-03-17-dork-rebuild-design.md`

---

## File Structure

All new files live under a `dork/` subdirectory in the repo root. The existing Python code stays untouched.

```
dork/
├── cmd/dork/main.go                    # Cobra root command, flag parsing, orchestration
├── internal/
│   ├── caseinfo/case.go                # Case struct, name parsing, YAML/JSON loading, CLI merge
│   ├── caseinfo/case_test.go           # Table-driven tests for parsing and merging
│   ├── dork/engine.go                  # Generate() orchestrator, filtering, sorting
│   ├── dork/engine_test.go             # Tests for generation, filtering, priority ordering
│   ├── dork/dork.go                    # Dork struct definition
│   ├── dork/dorks_social.go            # Social media dork definitions
│   ├── dork/dorks_records.go           # Public records dork definitions
│   ├── dork/dorks_financial.go         # Financial dork definitions
│   ├── dork/dorks_location.go          # Location dork definitions
│   ├── dork/dorks_forums.go            # Forum/community dork definitions
│   ├── dork/dorks_peopledb.go          # People search database URLs (region-tagged)
│   ├── output/browser.go              # Tab blaster (open URLs in default browser)
│   ├── output/discord.go              # Discord markdown formatter
│   ├── output/discord_test.go         # Tests for Discord output format
│   ├── output/jsonexport.go           # JSON export
│   ├── output/csvexport.go            # CSV export
│   ├── output/clipboard.go            # Clipboard copy
│   ├── output/stdout.go               # Default stdout output (label + URL per line)
│   ├── output/dashboard.go            # Local HTTP server serving embedded HTML
│   └── interactive/prompt.go          # Bubbletea interactive mode
├── web/
│   ├── dashboard.html                  # Embedded HTML/CSS/JS dashboard template
│   └── embed.go                        # Go embed directive for dashboard.html
├── go.mod
├── Makefile
└── README.md
```

NOTE: Package name is `caseinfo` (not `case`) since `case` is a Go reserved word.

---

## Task 1: Go Module & Skeleton

**Files:**
- Create: `dork/go.mod`
- Create: `dork/cmd/dork/main.go`
- Create: `dork/Makefile`

- [ ] **Step 1: Initialize Go module**

```bash
cd /home/dev/projects/dork-missing-person-finder
mkdir -p dork/cmd/dork
cd dork
go mod init github.com/gl0bal01/dork
```

- [ ] **Step 2: Create minimal main.go with Cobra**

Create `dork/cmd/dork/main.go` with:
- All input flags: `--name/-n`, `--location/-l`, `--age`, `--dob`, `--aka`, `--associates`, `--description`, `--case`
- All output flags: `--open`, `--dashboard`, `--export`, `--output/-o`
- All filter flags: `--category`, `--region`, `--engine`, `--delay`, `--version`
- Interactive flag: `-i`
- Placeholder `run` function that prints "not yet implemented"
- Version variable with ldflags injection point

```go
var version = "dev"
```

- [ ] **Step 3: Install Cobra dependency**

```bash
cd /home/dev/projects/dork-missing-person-finder/dork
go get github.com/spf13/cobra
go mod tidy
```

- [ ] **Step 4: Create Makefile**

Create `dork/Makefile` with targets:
- `build` — builds for current platform with version ldflags
- `install` — go install with ldflags
- `test` — go test ./... -v
- `release` — cross-compile linux-amd64, darwin-amd64, darwin-arm64, windows-amd64
- `clean` — remove binaries

- [ ] **Step 5: Verify it builds and runs**

```bash
cd /home/dev/projects/dork-missing-person-finder/dork
make build
./dork --version    # prints "dev"
./dork --help       # shows all flags
./dork -n "John Doe"  # prints "not yet implemented"
```

- [ ] **Step 6: Commit**

```bash
git add dork/
git commit -m "feat: scaffold Go project with Cobra CLI and all flags"
```

---

## Task 2: Case Struct & Name Parsing

**Files:**
- Create: `dork/internal/caseinfo/case.go`
- Create: `dork/internal/caseinfo/case_test.go`

- [ ] **Step 1: Write failing tests for name parsing**

Create `dork/internal/caseinfo/case_test.go` with table-driven tests for `ParseName`:

| Input | Expected First | Expected Last |
|-------|---------------|--------------|
| `"John Doe"` | `"John"` | `"Doe"` |
| `"Mary Jane Watson"` | `"Mary Jane"` | `"Watson"` |
| `"Anna Smith-Jones"` | `"Anna"` | `"Smith-Jones"` |
| `"Cher"` | `""` | `"Cher"` |
| `"  John   Doe  "` | `"John"` | `"Doe"` |
| `""` | `""` | `""` |

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/caseinfo/ -v
```

Expected: FAIL — package doesn't exist yet.

- [ ] **Step 3: Implement Case struct and ParseName**

Create `dork/internal/caseinfo/case.go` with:
- `Case` struct with all fields (Name, FirstName, LastName, Aliases, DOB, Age, Location, Description, Associates, Region, Categories, Engine)
- `ParseName(fullName string) (first, last string)` — split on last space, collapse multiple spaces, handle edge cases
- `LoadFromFile(path string) (*Case, error)` — YAML or JSON based on file extension
- `New(name string) *Case` — create from name string
- `(c *Case) Merge(overrides *Case)` — non-zero CLI values override file values

- [ ] **Step 4: Run tests to verify they pass**

```bash
go get gopkg.in/yaml.v3
go mod tidy
go test ./internal/caseinfo/ -v
```

- [ ] **Step 5: Add tests for LoadFromFile and Merge**

Add test cases:
- `TestLoadFromFile_YAML` — write temp YAML file, load it, verify all fields
- `TestLoadFromFile_JSON` — same with JSON
- `TestMerge` — verify non-empty overrides replace, empty overrides don't

- [ ] **Step 6: Run all case tests**

```bash
go test ./internal/caseinfo/ -v
```

- [ ] **Step 7: Commit**

```bash
git add dork/internal/caseinfo/
git commit -m "feat: add Case struct with name parsing, YAML/JSON loading, and CLI merge"
```

---

## Task 3: Dork Struct & Engine Core

**Files:**
- Create: `dork/internal/dork/dork.go`
- Create: `dork/internal/dork/engine.go`
- Create: `dork/internal/dork/engine_test.go`

- [ ] **Step 1: Create Dork struct**

Create `dork/internal/dork/dork.go` with:
- `Dork` struct: Query, Category, Region, Priority, Label
- `Engines` map: google, bing, duckduckgo, yandex URL templates
- `(d Dork) URL(engine string) string` — **the single canonical URL resolver.** If Query starts with `http://` or `https://`, returns it as-is (for people-db direct links). Otherwise wraps in the engine's search URL. All output handlers MUST use this method — no separate `resolveURL` function.

- [ ] **Step 2: Write failing tests for engine**

Create `dork/internal/dork/engine_test.go` with:
- `TestGenerate_BasicNameOnly` — generates at least some social dorks
- `TestFilter_ByCategory` — filtering by category returns correct subset
- `TestFilter_ByRegion` — global/us/all filtering works correctly:
  - `global` returns only global-tagged dorks
  - `us` returns global + us dorks
  - `all` returns everything
- `TestSort_ByPriority` — priority 3 before 2 before 1

- [ ] **Step 3: Run tests to verify they fail**

```bash
go test ./internal/dork/ -v
```

- [ ] **Step 4: Implement engine.go**

Create `dork/internal/dork/engine.go` with:
- `DorkGenerator` type: `func(c *caseinfo.Case) []Dork`
- `registry []DorkGenerator` — populated by init() in dork definition files
- `Register(fn DorkGenerator)` — adds to registry
- `Generate(c *caseinfo.Case) []Dork` — runs all registered generators
- `Filter(dorks []Dork, categories, regions []string) []Dork` — category and region filtering
  - `categories=["all"]` means no category filter
  - `regions=["global"]` means only global dorks
  - `regions=["us","ca"]` means global + us + ca
  - `regions=["all"]` means all regions
- `Sort(dorks []Dork) []Dork` — returns copy sorted by priority descending

- [ ] **Step 5: Run tests — Filter and Sort should pass, Generate may fail (no generators yet)**

```bash
go test ./internal/dork/ -v
```

- [ ] **Step 6: Commit**

```bash
git add dork/internal/dork/
git commit -m "feat: add Dork struct, engine core with Generate/Filter/Sort"
```

---

## Task 4: Dork Definitions (All Categories)

**Files:**
- Create: `dork/internal/dork/dorks_social.go`
- Create: `dork/internal/dork/dorks_records.go`
- Create: `dork/internal/dork/dorks_financial.go`
- Create: `dork/internal/dork/dorks_location.go`
- Create: `dork/internal/dork/dorks_forums.go`
- Create: `dork/internal/dork/dorks_peopledb.go`

Each file registers its generator via `init()`. All generators receive `*caseinfo.Case` and return `[]Dork`.

- [ ] **Step 1: Create social dorks**

Create `dork/internal/dork/dorks_social.go`:
- Global social: Facebook, LinkedIn, Instagram, Twitter/X, TikTok, YouTube, GitHub, Medium, Reddit (all Priority 3)
- Generic profile URL search (Priority 2)
- Alias cross-searches on social platforms (Priority 2)
- Region-specific: VK.com + OK.ru (`ru`), Copainsdavant (`fr`), StayFriends (`de`), Hyves archives (`nl`)
- All queries use location narrowing when available

- [ ] **Step 2: Create records dorks**

Create `dork/internal/dork/dorks_records.go`:
- Court, property, marriage records (Priority 2)
- PDF public records (Priority 2)
- Obituaries via legacy.com/findagrave.com (Priority 2)
- News articles (Priority 2)
- Newspaper archives, Archive.org (Priority 1)
- Inmate/corrections records (Priority 2)
- US-specific: mugshots.com, arrests.org, vinelink.com (tagged `us`)
- DOB and age narrowing when available

- [ ] **Step 3: Create financial dorks**

Create `dork/internal/dork/dorks_financial.go`:
- PayPal/Venmo/CashApp mentions (Priority 1)
- Bank/account mentions (Priority 1)
- Loan/mortgage/credit (Priority 1)
- Contact spreadsheets xlsx/csv (Priority 1)

- [ ] **Step 4: Create location dorks**

Create `dork/internal/dork/dorks_location.go`:
- Location mentions with city (Priority 2, when location available)
- Google Maps reviews (Priority 1)
- City/state/moved mentions (Priority 1)
- Travel/booking/reservation (Priority 1)
- Yelp/Foursquare (Priority 1)

- [ ] **Step 5: Create forums dorks**

Create `dork/internal/dork/dorks_forums.go`:
- Reddit, Quora (Priority 1)
- Forum profiles (Priority 1)
- Education .edu sites (Priority 1)
- Resume/CV PDFs (Priority 1)
- Contact cards vcf/vcard (Priority 1)
- Associates cross-reference (Priority 2)

- [ ] **Step 6: Create people-db dorks**

Create `dork/internal/dork/dorks_peopledb.go`:
- People-DB dorks are special: `Query` field holds a direct URL (not a search query)
- US: Spokeo, Whitepages, TruePeopleSearch, FastPeopleSearch, BeenVerified
- CA: Canada411, CanadaPeopleSearch, WhitePages.ca
- UK: 192.com, FindMyPast, BT Phone Book, UKElectoralRoll
- AU: WhitePages.com.au, PeopleFinder.com.au, ReverseAustralia
- RU: VK search, OK.ru search, Yandex People, NumBuster
- FR: PagesBlanches
- DE: DasTelefonbuch, Telefonbuch.de
- AT: Herold.at, DasTelefonbuch.at
- NL: DeTelefoongids, WhitePages.nl, Numberway.nl
- All Priority 3, each tagged with its country code

- [ ] **Step 7: Run all tests**

```bash
go test ./... -v
```

Expected: all tests pass including TestGenerate_BasicNameOnly.

- [ ] **Step 8: Commit**

```bash
git add dork/internal/dork/dorks_*.go
git commit -m "feat: add dork definitions for all 6 categories with multi-region support"
```

---

## Task 5: Output — Stdout Default & Discord Export

**Files:**
- Create: `dork/internal/output/stdout.go`
- Create: `dork/internal/output/discord.go`
- Create: `dork/internal/output/discord_test.go`

- [ ] **Step 1: Write failing test for Discord output**

Create `dork/internal/output/discord_test.go`:
- Pass a Case and list of dorks to `Discord(w, c, dorks, engine)`
- Verify output contains: name, location, category headers ("### Social", "### Records"), labels, Google search URLs

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/output/ -v
```

- [ ] **Step 3: Implement stdout.go and discord.go**

`stdout.go`:
- `Stdout(w io.Writer, dorks []Dork, engine string)` — one line per dork: `Label: URL`, grouped by category
- Use `d.URL(engine)` for all URL resolution (the method on Dork handles both direct URLs and search queries)

`discord.go`:
- `Discord(w io.Writer, c *Case, dorks []Dork, engine string)` — Discord markdown format
- Header with name + metadata (location, age, DOB)
- Groups by category in fixed order: social, records, financial, location, forums, people-db
- Each group: `### Category (N links)` then `- Label: URL` per dork

Use `groupByCategory()` helper to bucket dorks.

- [ ] **Step 4: Run tests**

```bash
go test ./internal/output/ -v
```

- [ ] **Step 5: Commit**

```bash
git add dork/internal/output/stdout.go dork/internal/output/discord.go dork/internal/output/discord_test.go
git commit -m "feat: add stdout and Discord output formatters"
```

---

## Task 6: Output — JSON, CSV, Clipboard

**Files:**
- Create: `dork/internal/output/jsonexport.go`
- Create: `dork/internal/output/csvexport.go`
- Create: `dork/internal/output/clipboard.go`

- [ ] **Step 1: Write failing tests for JSON and CSV**

Create `dork/internal/output/export_test.go`:
- `TestJSONFormat` — pass Case + dorks, verify output is valid JSON, contains case name, has correct number of results, each result has label/url/category/priority fields
- `TestCSVFormat` — pass dorks, verify output has header row, correct number of data rows, correct column order (label,category,priority,url)

```bash
go test ./internal/output/ -v -run TestJSON
go test ./internal/output/ -v -run TestCSV
```

Expected: FAIL — functions don't exist yet.

- [ ] **Step 2: Implement JSON export**

`jsonexport.go`:
- `JSON(w io.Writer, c *Case, dorks []Dork, engine string) error`
- Output struct: `{ "case": { name, location, age, dob }, "results": [{ label, url, category, priority }] }`
- Pretty-printed with 2-space indent

- [ ] **Step 2: Implement CSV export**

`csvexport.go`:
- `CSV(w io.Writer, dorks []Dork, engine string) error`
- Header: `label,category,priority,url`
- One row per dork

- [ ] **Step 4: Implement clipboard**

`clipboard.go`:
- `Clipboard(c *Case, dorks []Dork, engine string) error`
- Renders Discord format into a buffer, copies to system clipboard via `atotto/clipboard`

- [ ] **Step 5: Get dependencies and run tests**

```bash
go get github.com/atotto/clipboard
go mod tidy
go test ./... -v
```

- [ ] **Step 6: Commit**

```bash
git add dork/internal/output/jsonexport.go dork/internal/output/csvexport.go dork/internal/output/clipboard.go dork/internal/output/export_test.go
git commit -m "feat: add JSON, CSV, and clipboard output formatters with tests"
```

---

## Task 7: Browser Tab Blaster

**Files:**
- Create: `dork/internal/output/browser.go`

- [ ] **Step 1: Implement browser opener**

`browser.go`:
- `OpenInBrowser(dorks []Dork, engine string, delay time.Duration)` — iterates dorks, opens each URL, sleeps `delay` between tabs, logs warnings to stderr on failure
- `openURL(url string) error` — platform-specific: `xdg-open` (linux), `open` (darwin), `cmd /c start` (windows)

- [ ] **Step 2: Verify it compiles**

```bash
go build ./internal/output/
```

- [ ] **Step 3: Commit**

```bash
git add dork/internal/output/browser.go
git commit -m "feat: add cross-platform browser tab blaster"
```

---

## Task 8: Wire CLI to Engine & Outputs

**Files:**
- Modify: `dork/cmd/dork/main.go`

- [ ] **Step 1: Implement the run function**

Replace the placeholder `run` function with the full orchestration:

1. If `--case` provided, load case from file
2. Build CLI overrides from flags (split `--aka` and `--associates` on comma)
3. If case loaded from file, merge CLI overrides; otherwise create case from flags
4. Validate: name must be set
5. Generate all dorks via `dork.Generate(c)`
6. Parse `--category` and `--region` flags (split on comma)
7. Filter and sort dorks
8. Determine output writer (stdout or `--output` file)
9. Dispatch to output handler based on flags:
   - `--dashboard` → `output.ServeDashboard()` (standalone, returns)
   - `--open` → `output.OpenInBrowser()` with `--delay`
   - `--export discord/json/csv/clipboard` → appropriate formatter
   - No flags → `output.Stdout()` (default)
10. `--open` and `--export` are combinable
11. Exit codes per spec: 0 (success), 1 (input error — bad flags, malformed case file), 2 (runtime error — clipboard fail, browser launch fail). Use `os.Exit(1)` or `os.Exit(2)` explicitly, not `log.Fatal`.

Import all internal packages. Use blank import `_ "github.com/gl0bal01/dork/internal/dork"` to trigger init() registration of dork generators.

- [ ] **Step 2: Build and smoke test**

```bash
make build
./dork -n "John Doe"
./dork -n "John Doe" --region us --category social
./dork -n "John Doe" --export discord
./dork -n "John Doe" --export json
./dork -n "John Doe" --export csv
```

- [ ] **Step 3: Run all tests**

```bash
go test ./... -v
```

- [ ] **Step 4: Commit**

```bash
git add dork/cmd/dork/main.go
git commit -m "feat: wire CLI to dork engine and all output handlers"
```

---

## Task 9: Dashboard

**Files:**
- Create: `dork/web/dashboard.html`
- Create: `dork/web/embed.go`
- Create: `dork/internal/output/dashboard.go`
- Modify: `dork/cmd/dork/main.go` (import web package)

- [ ] **Step 1: Create dashboard HTML**

Create `dork/web/dashboard.html` — self-contained HTML/CSS/JS:
- Dark theme (#0d1117 background, #c9d1d9 text, #58a6ff accents)
- Case info summary at top
- Links grouped by category with headers showing count
- Each link has: status dropdown (pending/useful/dead end/reviewed), label, clickable URL, copy button
- "Copy all in category" button per category
- "Export to Discord" button — generates markdown of useful/pending items, copies to clipboard
- "Open All in Browser" button
- Toast notification for clipboard actions
- Data injected as `const data = /*DATA_PLACEHOLDER*/{};` — Go server replaces the placeholder

**IMPORTANT:** Use safe DOM methods (document.createElement, textContent, setAttribute) instead of string-based HTML injection for all dynamic content. Only use static HTML for the page structure.

- [ ] **Step 2: Create embed.go**

Create `dork/web/embed.go`:

```go
package web

import _ "embed"

//go:embed dashboard.html
var DashboardHTML string
```

- [ ] **Step 3: Implement dashboard.go**

`output/dashboard.go`:
- `ServeDashboard(c *Case, dorks []Dork, engine string, htmlTemplate string) error`
- Build JSON data blob with case info and resolved URLs
- Replace `/*DATA_PLACEHOLDER*/{}` in HTML template with JSON data
- Listen on `127.0.0.1:0` (random available port)
- Register single HTTP handler serving the HTML
- Print dashboard URL to stderr
- Open URL in browser
- Serve until Ctrl+C

- [ ] **Step 4: Update main.go to import web package and pass HTML to ServeDashboard**

```go
import "github.com/gl0bal01/dork/web"
// ...
if flagDashboard {
    return output.ServeDashboard(c, sorted, engine, web.DashboardHTML)
}
```

- [ ] **Step 5: Build and test dashboard**

```bash
make build
./dork -n "John Doe" --region us --dashboard
```

Expected: opens browser with dashboard.

- [ ] **Step 6: Commit**

```bash
git add dork/web/ dork/internal/output/dashboard.go dork/cmd/dork/main.go
git commit -m "feat: add local web dashboard with embedded HTML"
```

---

## Task 10: Interactive Mode

**Files:**
- Create: `dork/internal/interactive/prompt.go`
- Modify: `dork/cmd/dork/main.go`

- [ ] **Step 1: Install bubbletea and huh**

```bash
go get github.com/charmbracelet/huh
go mod tidy
```

- [ ] **Step 2: Implement interactive prompt**

Create `dork/internal/interactive/prompt.go`:
- `Run() (*caseinfo.Case, engine string, region string, category string, openBrowser bool, err error)`
- Step 1 form: name (required), location, age/DOB
- Step 2 form: aliases, associates, description
- Step 3 form: engine (select), regions (multi-select), category (select), open in browser? (confirm)
- Returns populated Case and user selections

Use `charmbracelet/huh` forms with `huh.NewInput`, `huh.NewSelect`, `huh.NewMultiSelect`, `huh.NewConfirm`.

- [ ] **Step 3: Wire interactive mode into main.go**

In `run()`, before case building logic:
- If `flagInteractive`, call `interactive.Run()`
- Use returned case and selections, override the relevant flag variables
- Continue with normal flow

- [ ] **Step 4: Build and test interactive mode**

```bash
make build
./dork -i
```

Expected: interactive prompts appear, selections work, results display.

- [ ] **Step 5: Commit**

```bash
git add dork/internal/interactive/ dork/cmd/dork/main.go
git commit -m "feat: add interactive mode with bubbletea/huh forms"
```

---

## Task 11: Shell Completions

**Files:**
- Modify: `dork/cmd/dork/main.go`

- [ ] **Step 1: Add completion subcommand**

Cobra v1.2+ includes `completion` automatically. If not available, add manually with `GenBashCompletion`, `GenZshCompletion`, `GenFishCompletion`, `GenPowerShellCompletion`.

- [ ] **Step 2: Register custom completions for enum flags**

In `init()`, use `rootCmd.RegisterFlagCompletionFunc` for:
- `--engine`: google, bing, duckduckgo, yandex
- `--region`: global, all, us, ca, uk, au, ru, fr, de, at, nl
- `--category`: all, social, records, financial, location, forums, people-db
- `--export`: discord, json, csv, clipboard

- [ ] **Step 3: Test completions**

```bash
./dork completion bash > /tmp/dork-completion.bash
source /tmp/dork-completion.bash
./dork --engine <TAB>
./dork --region <TAB>
```

- [ ] **Step 4: Commit**

```bash
git add dork/cmd/dork/main.go
git commit -m "feat: add shell completions for all enum flags"
```

---

## Task 12: README & Final Polish

**Files:**
- Create: `dork/README.md`
- Create: `dork/.gitignore`

- [ ] **Step 1: Create README**

Cover: what it is, install options (download binary / `go install` / build from source), quick start examples, all flags, case file format, regions table, shell completions setup.

- [ ] **Step 2: Add .gitignore**

```
dork
dork-*
```

- [ ] **Step 3: Full end-to-end test**

```bash
cd /home/dev/projects/dork-missing-person-finder/dork
go test ./... -v
make build
./dork --version
./dork -n "Test User" --export json
./dork -n "Test User" --export discord --region us,ca
./dork -n "Test User" --export csv --category social
./dork completion bash > /dev/null
```

- [ ] **Step 4: Commit**

```bash
git add dork/README.md dork/.gitignore
git commit -m "docs: add README and .gitignore for dork Go binary"
```

---

## Task 13: Integration Tests

**Files:**
- Create: `dork/internal/dork/integration_test.go`
- Create: `dork/cmd/dork/main_test.go`

- [ ] **Step 1: Write engine integration tests**

Create `dork/internal/dork/integration_test.go` — tests that exercise the full pipeline: Case → Generate → Filter → Sort.

```go
func TestFullPipeline_NameOnly(t *testing.T) {
    c := caseinfo.New("John Doe")
    dorks := Generate(c)
    filtered := Filter(dorks, []string{"all"}, []string{"global"})
    sorted := Sort(filtered)

    if len(sorted) == 0 {
        t.Fatal("pipeline produced no dorks")
    }
    // Verify priority ordering
    for i := 1; i < len(sorted); i++ {
        if sorted[i].Priority > sorted[i-1].Priority {
            t.Errorf("dorks not sorted by priority: index %d has %d > index %d has %d",
                i, sorted[i].Priority, i-1, sorted[i-1].Priority)
        }
    }
    // Verify no region-specific dorks leaked through global filter
    for _, d := range sorted {
        if d.Region != "global" {
            t.Errorf("global filter let through region %q dork: %s", d.Region, d.Label)
        }
    }
}

func TestFullPipeline_WithRegions(t *testing.T) {
    c := caseinfo.New("John Doe")
    dorks := Generate(c)
    filtered := Filter(dorks, []string{"all"}, []string{"us", "ca"})

    // Should have global + us + ca dorks
    regions := map[string]bool{}
    for _, d := range filtered {
        regions[d.Region] = true
    }
    if !regions["global"] {
        t.Error("missing global dorks")
    }
    if !regions["us"] {
        t.Error("missing us dorks")
    }
    if !regions["ca"] {
        t.Error("missing ca dorks")
    }
    // Should NOT have uk, ru, etc.
    for r := range regions {
        if r != "global" && r != "us" && r != "ca" {
            t.Errorf("unexpected region %q leaked through filter", r)
        }
    }
}

func TestFullPipeline_CategoryFilter(t *testing.T) {
    c := caseinfo.New("John Doe")
    dorks := Generate(c)
    filtered := Filter(dorks, []string{"social"}, []string{"all"})

    for _, d := range filtered {
        if d.Category != "social" {
            t.Errorf("category filter let through %q", d.Category)
        }
    }
}

func TestFullPipeline_WithAllCaseFields(t *testing.T) {
    c := &caseinfo.Case{
        Name:        "John Doe",
        Location:    "Seattle, WA",
        Age:         34,
        DOB:         "1990-01-15",
        Aliases:     []string{"JD", "Johnny"},
        Associates:  []string{"Jane Smith"},
        Description: "Red hair, tattoo",
    }
    c.FirstName, c.LastName = caseinfo.ParseName(c.Name)

    dorks := Generate(c)
    if len(dorks) == 0 {
        t.Fatal("no dorks generated with full case")
    }

    // Should have more dorks than name-only (aliases, associates add extra)
    nameOnly := Generate(caseinfo.New("John Doe"))
    if len(dorks) <= len(nameOnly) {
        t.Errorf("full case (%d dorks) should produce more than name-only (%d dorks)",
            len(dorks), len(nameOnly))
    }
}
```

- [ ] **Step 2: Write CLI integration tests**

Create `dork/cmd/dork/main_test.go` — tests that run the binary and check output:

```go
func TestCLI_DefaultOutput(t *testing.T) {
    cmd := exec.Command("go", "run", ".", "-n", "John Doe")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("CLI failed: %v\n%s", err, out)
    }
    if len(out) == 0 {
        t.Error("CLI produced no output")
    }
    if !strings.Contains(string(out), "Social") {
        t.Error("output missing Social category")
    }
}

func TestCLI_JSONExport(t *testing.T) {
    cmd := exec.Command("go", "run", ".", "-n", "John Doe", "--export", "json")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("CLI failed: %v\n%s", err, out)
    }
    var result map[string]interface{}
    if err := json.Unmarshal(out, &result); err != nil {
        t.Errorf("output is not valid JSON: %v", err)
    }
}

func TestCLI_NoNameError(t *testing.T) {
    cmd := exec.Command("go", "run", ".")
    err := cmd.Run()
    if err == nil {
        t.Error("CLI should fail without --name or --case")
    }
}

func TestCLI_CaseFile(t *testing.T) {
    dir := t.TempDir()
    caseFile := filepath.Join(dir, "test.yaml")
    os.WriteFile(caseFile, []byte("name: \"Test User\"\nlocation: \"NYC\"\n"), 0644)

    cmd := exec.Command("go", "run", ".", "--case", caseFile, "--export", "json")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("CLI failed: %v\n%s", err, out)
    }
    if !strings.Contains(string(out), "Test User") {
        t.Error("output missing case file name")
    }
}

func TestCLI_RegionFilter(t *testing.T) {
    cmd := exec.Command("go", "run", ".", "-n", "John Doe", "--region", "us", "--export", "json")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("CLI failed: %v\n%s", err, out)
    }
    if !strings.Contains(string(out), "Spokeo") {
        t.Error("us region should include Spokeo")
    }
}
```

- [ ] **Step 3: Run all tests**

```bash
cd /home/dev/projects/dork-missing-person-finder/dork
go test ./... -v
```

- [ ] **Step 4: Commit**

```bash
git add dork/internal/dork/integration_test.go dork/cmd/dork/main_test.go
git commit -m "test: add integration tests for engine pipeline and CLI"
```

---

## Task 14: GitHub Actions CI

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create CI workflow**

Create `.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.23', '1.24']
    runs-on: ${{ matrix.os }}
    defaults:
      run:
        working-directory: dork

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test ./... -v -race

      - name: Build
        run: go build -o dork ./cmd/dork

  lint:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: dork

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Run go vet
        run: go vet ./...

      - name: Check formatting
        run: |
          unformatted=$(gofmt -l .)
          if [ -n "$unformatted" ]; then
            echo "Files not formatted:"
            echo "$unformatted"
            exit 1
          fi

  build-release:
    needs: [test, lint]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    defaults:
      run:
        working-directory: dork

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Cross-compile
        run: make release

      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: dork-binaries
          path: |
            dork/dork-linux-amd64
            dork/dork-darwin-amd64
            dork/dork-darwin-arm64
            dork/dork-windows-amd64.exe
```

- [ ] **Step 2: Verify workflow syntax**

```bash
# Quick syntax check — ensure YAML is valid
python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))"
```

- [ ] **Step 3: Add Makefile lint target**

Add to `dork/Makefile`:

```makefile
lint:
	go vet ./...
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Files not formatted:"; echo "$$unformatted"; exit 1; \
	fi

fmt:
	gofmt -w .
```

- [ ] **Step 4: Run lint locally**

```bash
cd /home/dev/projects/dork-missing-person-finder/dork
make lint
make fmt
```

- [ ] **Step 5: Commit**

```bash
git add .github/workflows/ci.yml dork/Makefile
git commit -m "ci: add GitHub Actions workflow with cross-platform tests, lint, and release build"
```

---

## Task Summary

| Task | Description | Dependencies |
|------|-------------|-------------|
| 1 | Go module, Cobra CLI skeleton, Makefile | None |
| 2 | Case struct, name parsing, YAML/JSON loading | None |
| 3 | Dork struct, engine core (Generate/Filter/Sort) | Task 2 |
| 4 | Dork definitions (all 6 categories) | Task 3 |
| 5 | Stdout + Discord output | Tasks 3, 4 |
| 6 | JSON, CSV, Clipboard output | Task 5 |
| 7 | Browser tab blaster | Task 3 |
| 8 | Wire CLI to engine + outputs | Tasks 1-7 |
| 9 | Dashboard (HTML + Go server) | Task 8 |
| 10 | Interactive mode (bubbletea/huh) | Task 8 |
| 11 | Shell completions | Task 8 |
| 12 | README + final polish | All |
| 13 | Integration tests (pipeline + CLI) | Task 8 |
| 14 | GitHub Actions CI (test + lint + build) | Task 13 |

**Parallelizable:** Tasks 1+2 can run in parallel. Tasks 5+6+7 can run in parallel. Tasks 9+10+11+13 can run in parallel. Task 14 after 13.
