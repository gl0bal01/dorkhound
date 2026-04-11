package interactive

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

// Result holds the gathered information from the interactive prompts.
type Result struct {
	Case        *caseinfo.Case
	Engine      string
	Region      string
	Category    string
	OpenBrowser bool
}

// Run launches the interactive prompt flow and returns the collected result.
func Run() (*Result, error) {
	var name, location, age, dob string
	var aka, associates, description string
	var engine, category string
	var regions []string
	var openBrowser bool

	// Step 1: Required info
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Person's full name").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}).Value(&name),
			huh.NewInput().Title("Last known location (optional)").Value(&location),
			huh.NewInput().Title("Date of birth (optional, YYYY-MM-DD)").Value(&dob),
			huh.NewInput().Title("Approximate age (optional)").Value(&age),
		),
	).Run()
	if err != nil {
		return nil, err
	}

	var emails, phones, usernames, photoURL string

	// Step 2: Additional info
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Aliases (comma-separated, optional)").Value(&aka),
			huh.NewInput().Title("Known associates (comma-separated, optional)").Value(&associates),
			huh.NewInput().Title("Physical description (optional)").Value(&description),
			huh.NewInput().Title("Email addresses (comma-separated, optional)").Value(&emails),
			huh.NewInput().Title("Phone numbers (comma-separated, optional)").Value(&phones),
			huh.NewInput().Title("Usernames/handles (comma-separated, optional)").Value(&usernames),
			huh.NewInput().Title("Photo URL (optional, for reverse image search)").Value(&photoURL),
		),
	).Run()
	if err != nil {
		return nil, err
	}

	// Step 3: Search options
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title("Search engine").
				Options(
					huh.NewOption("Google", "google"),
					huh.NewOption("Bing", "bing"),
					huh.NewOption("DuckDuckGo", "duckduckgo"),
					huh.NewOption("Yandex", "yandex"),
				).Value(&engine),
			huh.NewMultiSelect[string]().Title("Regions to include").
				Options(
					huh.NewOption("Global only", "global"),
					huh.NewOption("US", "us"),
					huh.NewOption("Canada", "ca"),
					huh.NewOption("UK", "uk"),
					huh.NewOption("Australia", "au"),
					huh.NewOption("Russia", "ru"),
					huh.NewOption("France", "fr"),
					huh.NewOption("Germany", "de"),
					huh.NewOption("Austria", "at"),
					huh.NewOption("Netherlands", "nl"),
					huh.NewOption("All regions", "all"),
				).Value(&regions),
			huh.NewSelect[string]().Title("Category filter").
				Options(
					huh.NewOption("All categories", "all"),
					huh.NewOption("Social media", "social"),
					huh.NewOption("Public records", "records"),
					huh.NewOption("Financial", "financial"),
					huh.NewOption("Location", "location"),
					huh.NewOption("Forums", "forums"),
					huh.NewOption("People databases", "people-db"),
					huh.NewOption("Email", "email"),
					huh.NewOption("Phone", "phone"),
					huh.NewOption("Username", "username"),
					huh.NewOption("Cache/Archive", "cache"),
					huh.NewOption("Documents", "documents"),
					huh.NewOption("Dating", "dating"),
					huh.NewOption("Marketplace", "marketplace"),
					huh.NewOption("Nuclei results", "nuclei"),
					huh.NewOption("Image/reverse search", "image"),
					huh.NewOption("Twitter/X", "twitter"),
					huh.NewOption("Reddit", "reddit"),
					huh.NewOption("Fundraiser/Memorial", "fundraiser"),
					huh.NewOption("Telegram", "telegram"),
					huh.NewOption("Vehicle", "vehicle"),
					huh.NewOption("Crypto", "crypto"),
				).Value(&category),
			huh.NewConfirm().Title("Open results in browser?").Value(&openBrowser),
		),
	).Run()
	if err != nil {
		return nil, err
	}

	// Build case
	c := caseinfo.New(name)
	c.Location = location
	c.DOB = dob
	c.Description = description
	if aka != "" {
		c.Aliases = caseinfo.SplitTrim(aka)
	}
	if associates != "" {
		c.Associates = caseinfo.SplitTrim(associates)
	}
	if emails != "" {
		c.Emails = caseinfo.SplitTrim(emails)
	}
	if phones != "" {
		c.Phones = caseinfo.SplitTrim(phones)
	}
	if usernames != "" {
		c.Usernames = caseinfo.SplitTrim(usernames)
	}
	c.PhotoURL = photoURL
	if age != "" {
		parsed, err := strconv.Atoi(strings.TrimSpace(age))
		if err != nil {
			return nil, fmt.Errorf("invalid age %q: must be a number", age)
		}
		c.Age = parsed
	}

	regionStr := "global"
	if len(regions) > 0 {
		regionStr = strings.Join(regions, ",")
	}

	return &Result{
		Case:        c,
		Engine:      engine,
		Region:      regionStr,
		Category:    category,
		OpenBrowser: openBrowser,
	}, nil
}
