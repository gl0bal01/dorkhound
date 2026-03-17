package caseinfo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantFirst string
		wantLast  string
	}{
		{
			name:      "simple two-part name",
			input:     "John Doe",
			wantFirst: "John",
			wantLast:  "Doe",
		},
		{
			name:      "three-part name splits on last space",
			input:     "Mary Jane Watson",
			wantFirst: "Mary Jane",
			wantLast:  "Watson",
		},
		{
			name:      "hyphenated last name",
			input:     "Anna Smith-Jones",
			wantFirst: "Anna",
			wantLast:  "Smith-Jones",
		},
		{
			name:      "single name goes to last",
			input:     "Cher",
			wantFirst: "",
			wantLast:  "Cher",
		},
		{
			name:      "extra whitespace collapsed",
			input:     "  John   Doe  ",
			wantFirst: "John",
			wantLast:  "Doe",
		},
		{
			name:      "empty string",
			input:     "",
			wantFirst: "",
			wantLast:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first, last := ParseName(tt.input)
			if first != tt.wantFirst {
				t.Errorf("ParseName(%q) first = %q, want %q", tt.input, first, tt.wantFirst)
			}
			if last != tt.wantLast {
				t.Errorf("ParseName(%q) last = %q, want %q", tt.input, last, tt.wantLast)
			}
		})
	}
}

func TestLoadFromFile_YAML(t *testing.T) {
	c := &Case{
		Name:        "Jane Doe",
		Aliases:     []string{"JD", "Janey"},
		DOB:         "1990-05-15",
		Age:         35,
		Location:    "New York, NY",
		Description: "Missing since 2024",
		Associates:  []string{"John Smith", "Bob Brown"},
		Region:      "northeast",
		Categories:  []string{"social_media", "public_records"},
		Engine:      "google",
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal YAML: %v", err)
	}

	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "case.yaml")
	if err := os.WriteFile(yamlPath, data, 0644); err != nil {
		t.Fatalf("failed to write temp YAML: %v", err)
	}

	loaded, err := LoadFromFile(yamlPath)
	if err != nil {
		t.Fatalf("LoadFromFile(%q) error: %v", yamlPath, err)
	}

	if loaded.Name != "Jane Doe" {
		t.Errorf("Name = %q, want %q", loaded.Name, "Jane Doe")
	}
	if loaded.FirstName != "Jane" {
		t.Errorf("FirstName = %q, want %q", loaded.FirstName, "Jane")
	}
	if loaded.LastName != "Doe" {
		t.Errorf("LastName = %q, want %q", loaded.LastName, "Doe")
	}
	if len(loaded.Aliases) != 2 || loaded.Aliases[0] != "JD" || loaded.Aliases[1] != "Janey" {
		t.Errorf("Aliases = %v, want [JD Janey]", loaded.Aliases)
	}
	if loaded.DOB != "1990-05-15" {
		t.Errorf("DOB = %q, want %q", loaded.DOB, "1990-05-15")
	}
	if loaded.Age != 35 {
		t.Errorf("Age = %d, want %d", loaded.Age, 35)
	}
	if loaded.Location != "New York, NY" {
		t.Errorf("Location = %q, want %q", loaded.Location, "New York, NY")
	}
	if loaded.Description != "Missing since 2024" {
		t.Errorf("Description = %q, want %q", loaded.Description, "Missing since 2024")
	}
	if len(loaded.Associates) != 2 {
		t.Errorf("Associates = %v, want 2 items", loaded.Associates)
	}
	if loaded.Region != "northeast" {
		t.Errorf("Region = %q, want %q", loaded.Region, "northeast")
	}
	if len(loaded.Categories) != 2 {
		t.Errorf("Categories = %v, want 2 items", loaded.Categories)
	}
	if loaded.Engine != "google" {
		t.Errorf("Engine = %q, want %q", loaded.Engine, "google")
	}

	// Also test .yml extension
	ymlPath := filepath.Join(dir, "case.yml")
	if err := os.WriteFile(ymlPath, data, 0644); err != nil {
		t.Fatalf("failed to write temp .yml: %v", err)
	}
	loaded2, err := LoadFromFile(ymlPath)
	if err != nil {
		t.Fatalf("LoadFromFile(%q) error: %v", ymlPath, err)
	}
	if loaded2.Name != "Jane Doe" {
		t.Errorf(".yml: Name = %q, want %q", loaded2.Name, "Jane Doe")
	}
}

func TestLoadFromFile_JSON(t *testing.T) {
	c := &Case{
		Name:        "Bob Wilson",
		Aliases:     []string{"Bobby"},
		DOB:         "1985-12-01",
		Age:         40,
		Location:    "Chicago, IL",
		Description: "Last seen downtown",
		Associates:  []string{"Alice"},
		Region:      "midwest",
		Categories:  []string{"news"},
		Engine:      "bing",
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "case.json")
	if err := os.WriteFile(jsonPath, data, 0644); err != nil {
		t.Fatalf("failed to write temp JSON: %v", err)
	}

	loaded, err := LoadFromFile(jsonPath)
	if err != nil {
		t.Fatalf("LoadFromFile(%q) error: %v", jsonPath, err)
	}

	if loaded.Name != "Bob Wilson" {
		t.Errorf("Name = %q, want %q", loaded.Name, "Bob Wilson")
	}
	if loaded.FirstName != "Bob" {
		t.Errorf("FirstName = %q, want %q", loaded.FirstName, "Bob")
	}
	if loaded.LastName != "Wilson" {
		t.Errorf("LastName = %q, want %q", loaded.LastName, "Wilson")
	}
	if loaded.Age != 40 {
		t.Errorf("Age = %d, want %d", loaded.Age, 40)
	}
	if loaded.Engine != "bing" {
		t.Errorf("Engine = %q, want %q", loaded.Engine, "bing")
	}
}

func TestLoadFromFile_UnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	txtPath := filepath.Join(dir, "case.txt")
	if err := os.WriteFile(txtPath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	_, err := LoadFromFile(txtPath)
	if err == nil {
		t.Error("LoadFromFile with .txt extension should return error")
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/case.yaml")
	if err == nil {
		t.Error("LoadFromFile with nonexistent path should return error")
	}
}

func TestNew(t *testing.T) {
	c := New("Mary Jane Watson")
	if c.Name != "Mary Jane Watson" {
		t.Errorf("Name = %q, want %q", c.Name, "Mary Jane Watson")
	}
	if c.FirstName != "Mary Jane" {
		t.Errorf("FirstName = %q, want %q", c.FirstName, "Mary Jane")
	}
	if c.LastName != "Watson" {
		t.Errorf("LastName = %q, want %q", c.LastName, "Watson")
	}
}

func TestMerge(t *testing.T) {
	t.Run("non-empty overrides replace base values", func(t *testing.T) {
		base := &Case{
			Name:        "John Doe",
			FirstName:   "John",
			LastName:    "Doe",
			DOB:         "1990-01-01",
			Age:         35,
			Location:    "NYC",
			Description: "original description",
			Region:      "northeast",
			Engine:      "google",
			Aliases:     []string{"JD"},
			Associates:  []string{"Friend1"},
			Categories:  []string{"social_media"},
		}

		overrides := &Case{
			Location: "Los Angeles, CA",
			Age:      36,
			Engine:   "bing",
		}

		base.Merge(overrides)

		if base.Location != "Los Angeles, CA" {
			t.Errorf("Location = %q, want %q", base.Location, "Los Angeles, CA")
		}
		if base.Age != 36 {
			t.Errorf("Age = %d, want %d", base.Age, 36)
		}
		if base.Engine != "bing" {
			t.Errorf("Engine = %q, want %q", base.Engine, "bing")
		}
		// Unchanged fields
		if base.Name != "John Doe" {
			t.Errorf("Name should remain %q, got %q", "John Doe", base.Name)
		}
		if base.DOB != "1990-01-01" {
			t.Errorf("DOB should remain %q, got %q", "1990-01-01", base.DOB)
		}
		if base.Description != "original description" {
			t.Errorf("Description should remain unchanged")
		}
	})

	t.Run("empty overrides do not replace", func(t *testing.T) {
		base := &Case{
			Name:      "John Doe",
			FirstName: "John",
			LastName:  "Doe",
			DOB:       "1990-01-01",
			Age:       35,
			Location:  "NYC",
			Engine:    "google",
		}

		overrides := &Case{} // all zero values

		base.Merge(overrides)

		if base.Name != "John Doe" {
			t.Errorf("Name should remain %q, got %q", "John Doe", base.Name)
		}
		if base.DOB != "1990-01-01" {
			t.Errorf("DOB should remain %q, got %q", "1990-01-01", base.DOB)
		}
		if base.Age != 35 {
			t.Errorf("Age should remain %d, got %d", 35, base.Age)
		}
		if base.Location != "NYC" {
			t.Errorf("Location should remain %q, got %q", "NYC", base.Location)
		}
	})

	t.Run("name override triggers re-parse", func(t *testing.T) {
		base := &Case{
			Name:      "John Doe",
			FirstName: "John",
			LastName:  "Doe",
		}

		overrides := &Case{
			Name: "Mary Jane Watson",
		}

		base.Merge(overrides)

		if base.Name != "Mary Jane Watson" {
			t.Errorf("Name = %q, want %q", base.Name, "Mary Jane Watson")
		}
		if base.FirstName != "Mary Jane" {
			t.Errorf("FirstName = %q, want %q", base.FirstName, "Mary Jane")
		}
		if base.LastName != "Watson" {
			t.Errorf("LastName = %q, want %q", base.LastName, "Watson")
		}
	})

	t.Run("slice overrides replace when non-empty", func(t *testing.T) {
		base := &Case{
			Aliases:    []string{"JD"},
			Associates: []string{"Friend1"},
			Categories: []string{"social_media"},
		}

		overrides := &Case{
			Aliases: []string{"Johnny", "J"},
		}

		base.Merge(overrides)

		if len(base.Aliases) != 2 || base.Aliases[0] != "Johnny" {
			t.Errorf("Aliases = %v, want [Johnny J]", base.Aliases)
		}
		// Unchanged slices
		if len(base.Associates) != 1 || base.Associates[0] != "Friend1" {
			t.Errorf("Associates should remain unchanged, got %v", base.Associates)
		}
		if len(base.Categories) != 1 || base.Categories[0] != "social_media" {
			t.Errorf("Categories should remain unchanged, got %v", base.Categories)
		}
	})
}

func TestSplitTrim(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty string", "", nil},
		{"single element", "social", []string{"social"}},
		{"multiple elements", "us,ca,uk", []string{"us", "ca", "uk"}},
		{"with whitespace", " us , ca , uk ", []string{"us", "ca", "uk"}},
		{"trailing comma", "us,ca,", []string{"us", "ca"}},
		{"all whitespace elements", " , , ", nil},
		{"single with whitespace", "  social  ", []string{"social"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitTrim(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("SplitTrim(%q) = %v (len %d), want %v (len %d)", tt.input, got, len(got), tt.want, len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("SplitTrim(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}
