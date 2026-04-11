package dork

import (
	"fmt"
	"sort"
	"strings"

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

// Dedup removes duplicate dorks based on Query string. When duplicates exist,
// the dork with the highest Priority wins (ties broken by first occurrence).
// Input order is preserved for survivors.
func Dedup(dorks []Dork) []Dork {
	seen := make(map[string]int) // query -> index of best dork so far
	var best []Dork
	for _, d := range dorks {
		if idx, ok := seen[d.Query]; ok {
			if d.Priority > best[idx].Priority {
				best[idx] = d
			}
			continue
		}
		seen[d.Query] = len(best)
		best = append(best, d)
	}
	return best
}

// Stats returns a formatted breakdown of dorks by category, region, and priority.
func Stats(dorks []Dork) string {
	catCount := make(map[string]int)
	regCount := make(map[string]int)
	priCount := make(map[int]int)

	for _, d := range dorks {
		catCount[d.Category]++
		regCount[d.Region]++
		priCount[d.Priority]++
	}

	type kv struct {
		key string
		val int
	}

	sortedKV := func(m map[string]int) []kv {
		var pairs []kv
		for k, v := range m {
			pairs = append(pairs, kv{k, v})
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].val != pairs[j].val {
				return pairs[i].val > pairs[j].val
			}
			return pairs[i].key < pairs[j].key
		})
		return pairs
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Dorks: %d total\n", len(dorks))

	b.WriteString("\nBy category:\n")
	for _, p := range sortedKV(catCount) {
		fmt.Fprintf(&b, "  %-20s %d\n", p.key, p.val)
	}

	b.WriteString("\nBy region:\n")
	for _, p := range sortedKV(regCount) {
		fmt.Fprintf(&b, "  %-20s %d\n", p.key, p.val)
	}

	b.WriteString("\nBy priority:\n")
	priorities := []int{3, 2, 1}
	labels := map[int]string{3: "3 (high)", 2: "2 (medium)", 1: "1 (low)"}
	for _, pri := range priorities {
		if n, ok := priCount[pri]; ok {
			fmt.Fprintf(&b, "  %-20s %d\n", labels[pri], n)
		}
	}

	return b.String()
}

// ApplyNoiseFilter appends negative-site operators to each dork's Query to
// suppress common noise domains (Pinterest recipe spam, Amazon, baby-name
// SEO pages). It only mutates dorks whose Query looks like a search query
// (not direct http(s) URLs).
func ApplyNoiseFilter(dorks []Dork) []Dork {
	const suffix = ` -site:pinterest.com -site:amazon.com -site:etsy.com -"baby name"`
	out := make([]Dork, len(dorks))
	for i, d := range dorks {
		out[i] = d
		if !strings.HasPrefix(d.Query, "http://") && !strings.HasPrefix(d.Query, "https://") {
			out[i].Query = d.Query + suffix
		}
	}
	return out
}
