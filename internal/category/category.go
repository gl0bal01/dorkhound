// Package category is the single source of truth for dork categories.
// CLI flags, interactive prompts, exports, and the dashboard all share this
// catalog so a new category added here flows everywhere automatically.
package category

import "sort"

// Catalog enumerates every category dorkhound can produce, in the order they
// should appear by default (highest-signal categories first).
var Catalog = []Entry{
	{Slug: "image", Title: "Image", Description: "Reverse image search and lookalike pivots"},
	{Slug: "username", Title: "Username", Description: "Username/handle search across 30+ sites"},
	{Slug: "direct-profile", Title: "Direct profile", Description: "Direct profile-URL templates per network"},
	{Slug: "email", Title: "Email", Description: "Email-based search and breach checks"},
	{Slug: "phone", Title: "Phone", Description: "Phone-number variants and lookups"},
	{Slug: "gravatar", Title: "Gravatar", Description: "Gravatar email→avatar lookup"},
	{Slug: "github", Title: "GitHub", Description: "GitHub profile/code/commit search"},
	{Slug: "social", Title: "Social", Description: "Generic social-platform dorks"},
	{Slug: "people-db", Title: "People-DB", Description: "People-search aggregators"},
	{Slug: "nuclei", Title: "Nuclei", Description: "Results imported from nuclei OSINT templates"},
	{Slug: "cache", Title: "Cache", Description: "Web archives and cache lookups"},
	{Slug: "records", Title: "Records", Description: "Public records and government data"},
	{Slug: "academic", Title: "Academic", Description: "Academic publications and university directories"},
	{Slug: "documents", Title: "Documents", Description: "Files exposed on the web (PDF, DOC, XLS)"},
	{Slug: "forums", Title: "Forums", Description: "Forum posts and discussion boards"},
	{Slug: "financial", Title: "Financial", Description: "Crowdfunding, fundraising and financial mentions"},
	{Slug: "location", Title: "Location", Description: "Geographic and check-in pivots"},
	{Slug: "dating", Title: "Dating", Description: "Dating-site profile dorks"},
	{Slug: "marketplace", Title: "Marketplace", Description: "Classifieds and resale-site mentions"},
	{Slug: "twitter", Title: "Twitter/X", Description: "Twitter/X-specific and nitter mirrors"},
	{Slug: "reddit", Title: "Reddit", Description: "Reddit-specific deep search"},
	{Slug: "fundraiser", Title: "Fundraiser", Description: "GoFundMe, memorials, obituaries"},
	{Slug: "telegram", Title: "Telegram", Description: "Telegram channel and handle search"},
	{Slug: "vehicle", Title: "Vehicle", Description: "License plate and VIN lookups"},
	{Slug: "crypto", Title: "Crypto", Description: "Crypto address and wallet pivots"},
}

// Entry describes a single category.
type Entry struct {
	Slug        string
	Title       string
	Description string
}

// Slugs returns every category slug in catalog order.
func Slugs() []string {
	out := make([]string, len(Catalog))
	for i, e := range Catalog {
		out[i] = e.Slug
	}
	return out
}

// AllSlugs returns the "all" sentinel plus every catalog slug.
// Used by CLI completion and the --list-categories command.
func AllSlugs() []string {
	out := make([]string, 0, len(Catalog)+1)
	out = append(out, "all")
	out = append(out, Slugs()...)
	return out
}

// Titles returns a slug→title map for export rendering.
func Titles() map[string]string {
	m := make(map[string]string, len(Catalog))
	for _, e := range Catalog {
		m[e.Slug] = e.Title
	}
	return m
}

// IsKnown reports whether slug is a recognized category (including "all").
func IsKnown(slug string) bool {
	if slug == "all" {
		return true
	}
	for _, e := range Catalog {
		if e.Slug == slug {
			return true
		}
	}
	return false
}

// OrderForExport returns the catalog order for known categories that are
// present, followed by any unknown categories sorted alphabetically.
// Used by Discord/TraceLabs/dashboard so a new category never silently drops.
func OrderForExport(present []string) []string {
	presentSet := make(map[string]bool, len(present))
	for _, slug := range present {
		presentSet[slug] = true
	}

	var ordered []string
	seen := make(map[string]bool)
	for _, e := range Catalog {
		if presentSet[e.Slug] {
			ordered = append(ordered, e.Slug)
			seen[e.Slug] = true
		}
	}
	var extras []string
	for slug := range presentSet {
		if !seen[slug] {
			extras = append(extras, slug)
		}
	}
	sort.Strings(extras)
	return append(ordered, extras...)
}
