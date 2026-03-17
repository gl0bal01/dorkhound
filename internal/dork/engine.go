package dork

import (
	"sort"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

// DorkGenerator is a function that produces dorks for a given case.
type DorkGenerator func(c *caseinfo.Case) []Dork

// registry holds all registered dork generators, populated by init() in dork
// definition files.
var registry []DorkGenerator

// register appends a DorkGenerator to the global registry.
// Called only from init() functions within this package.
func register(fn DorkGenerator) {
	registry = append(registry, fn)
}

// Generate runs all registered generators against the given case and collects
// the resulting dorks into a single slice.
func Generate(c *caseinfo.Case) []Dork {
	var all []Dork
	for _, gen := range registry {
		all = append(all, gen(c)...)
	}
	return all
}

// Filter returns dorks that match the given category and region constraints.
//
// Category rules:
//   - categories=["all"] -> no category filter (all categories pass)
//   - otherwise, only dorks whose Category is in the list pass
//
// Region rules:
//   - regions=["all"] -> no region filter (all regions pass)
//   - regions=["global"] -> only dorks with Region=="global"
//   - regions=["us","ca"] -> dorks with Region=="global" OR Region in the list
func Filter(dorks []Dork, categories, regions []string) []Dork {
	catSet := makeSet(categories)
	regSet := makeSet(regions)

	filterCat := !catSet["all"]
	filterReg := !regSet["all"]

	var result []Dork
	for _, d := range dorks {
		if filterCat && !catSet[d.Category] {
			continue
		}
		if filterReg {
			// When filtering by region, "global" dorks only pass if "global"
			// is explicitly in the region set. Specific regions pass if they
			// are in the set OR if the dork is Region=="global" (global dorks
			// are always included when any specific region is requested).
			if !regSet[d.Region] && d.Region != "global" {
				continue
			}
			// If only "global" is in the set, non-global dorks must not pass.
			if onlyGlobal(regSet) && d.Region != "global" {
				continue
			}
		}
		result = append(result, d)
	}
	return result
}

// Sort returns a new slice sorted by Priority descending (highest first).
// The sort is stable, preserving the original order among equal priorities.
func Sort(dorks []Dork) []Dork {
	sorted := make([]Dork, len(dorks))
	copy(sorted, dorks)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})
	return sorted
}

// makeSet converts a string slice into a set (map) for O(1) lookups.
func makeSet(items []string) map[string]bool {
	s := make(map[string]bool, len(items))
	for _, item := range items {
		s[item] = true
	}
	return s
}

// onlyGlobal returns true if the set contains only the key "global".
func onlyGlobal(s map[string]bool) bool {
	return len(s) == 1 && s["global"]
}
