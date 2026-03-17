package caseinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Case holds information about a missing person case.
type Case struct {
	Name        string   `yaml:"name" json:"name"`
	FirstName   string   `yaml:"-" json:"-"`
	LastName    string   `yaml:"-" json:"-"`
	Aliases     []string `yaml:"aliases" json:"aliases"`
	DOB         string   `yaml:"dob" json:"dob"`
	Age         int      `yaml:"age" json:"age"`
	Location    string   `yaml:"location" json:"location"`
	Description string   `yaml:"description" json:"description"`
	Associates  []string `yaml:"associates" json:"associates"`
	Region      string   `yaml:"region" json:"region"`
	Categories  []string `yaml:"categories" json:"categories"`
	Engine      string   `yaml:"engine" json:"engine"`
}

// ParseName splits a full name into first and last name components.
// It splits on the last space: "Mary Jane Watson" -> ("Mary Jane", "Watson").
// A single name goes entirely to last: "Cher" -> ("", "Cher").
// Empty input returns ("", ""). Multiple spaces are collapsed.
func ParseName(fullName string) (first, last string) {
	// Collapse multiple spaces and trim
	fields := strings.Fields(fullName)
	if len(fields) == 0 {
		return "", ""
	}
	if len(fields) == 1 {
		return "", fields[0]
	}
	last = fields[len(fields)-1]
	first = strings.Join(fields[:len(fields)-1], " ")
	return first, last
}

// LoadFromFile loads a Case from a YAML (.yaml/.yml) or JSON (.json) file.
// After loading, FirstName and LastName are derived from Name.
func LoadFromFile(path string) (*Case, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading case file: %w", err)
	}

	c := &Case{}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, c); err != nil {
			return nil, fmt.Errorf("parsing YAML: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, c); err != nil {
			return nil, fmt.Errorf("parsing JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension %q: use .yaml, .yml, or .json", ext)
	}

	c.FirstName, c.LastName = ParseName(c.Name)
	return c, nil
}

// New creates a new Case from a name string, setting FirstName and LastName.
func New(name string) *Case {
	c := &Case{Name: name}
	c.FirstName, c.LastName = ParseName(name)
	return c
}

// SplitTrim splits a comma-separated string and trims whitespace from each part,
// discarding empty elements.
func SplitTrim(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// Merge applies non-zero/non-empty override values from overrides onto the
// receiver Case. If Name is overridden, FirstName and LastName are re-parsed.
func (c *Case) Merge(overrides *Case) {
	if overrides.Name != "" {
		c.Name = overrides.Name
		c.FirstName, c.LastName = ParseName(c.Name)
	}
	if len(overrides.Aliases) > 0 {
		c.Aliases = overrides.Aliases
	}
	if overrides.DOB != "" {
		c.DOB = overrides.DOB
	}
	if overrides.Age != 0 {
		c.Age = overrides.Age
	}
	if overrides.Location != "" {
		c.Location = overrides.Location
	}
	if overrides.Description != "" {
		c.Description = overrides.Description
	}
	if len(overrides.Associates) > 0 {
		c.Associates = overrides.Associates
	}
	if overrides.Region != "" {
		c.Region = overrides.Region
	}
	if len(overrides.Categories) > 0 {
		c.Categories = overrides.Categories
	}
	if overrides.Engine != "" {
		c.Engine = overrides.Engine
	}
}
