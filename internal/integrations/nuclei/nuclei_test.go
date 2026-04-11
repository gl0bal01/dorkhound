package nuclei

import (
	"testing"
)

func TestParseJSONL_Empty(t *testing.T) {
	dorks := parseJSONL([]byte{}, "testuser")
	if len(dorks) != 0 {
		t.Errorf("parseJSONL(empty) = %d dorks, want 0", len(dorks))
	}

	dorks = parseJSONL([]byte("\n\n"), "testuser")
	if len(dorks) != 0 {
		t.Errorf("parseJSONL(only newlines) = %d dorks, want 0", len(dorks))
	}
}

func TestParseJSONL_Fixture(t *testing.T) {
	fixture := []byte(`{"template-id":"twitter-username","info":{"name":"Twitter Username Found","severity":"medium","tags":["osint","osint-social"]},"matched-at":"https://twitter.com/jdoe42","host":"twitter.com"}
{"template-id":"github-username","info":{"name":"GitHub Username Found","severity":"info","tags":["osint"]},"matched-at":"https://github.com/jdoe42","host":"github.com"}
{"template-id":"instagram-username","info":{"name":"Instagram Username Found","severity":"medium","tags":["osint-social"]},"matched-at":"","matched":"https://instagram.com/jdoe42","host":"instagram.com"}
`)
	dorks := parseJSONL(fixture, "jdoe42")
	if len(dorks) != 3 {
		t.Fatalf("parseJSONL: got %d dorks, want 3", len(dorks))
	}

	// First: medium severity -> priority 3
	if dorks[0].Query != "https://twitter.com/jdoe42" {
		t.Errorf("dorks[0].Query = %q, want https://twitter.com/jdoe42", dorks[0].Query)
	}
	if dorks[0].Priority != 3 {
		t.Errorf("dorks[0].Priority = %d, want 3 (medium severity)", dorks[0].Priority)
	}
	if dorks[0].Category != "nuclei" {
		t.Errorf("dorks[0].Category = %q, want nuclei", dorks[0].Category)
	}

	// Second: info severity -> priority 2
	if dorks[1].Query != "https://github.com/jdoe42" {
		t.Errorf("dorks[1].Query = %q, want https://github.com/jdoe42", dorks[1].Query)
	}
	if dorks[1].Priority != 2 {
		t.Errorf("dorks[1].Priority = %d, want 2 (info severity)", dorks[1].Priority)
	}

	// Third: matched-at empty, falls back to matched
	if dorks[2].Query != "https://instagram.com/jdoe42" {
		t.Errorf("dorks[2].Query = %q, want https://instagram.com/jdoe42 (from matched fallback)", dorks[2].Query)
	}

	// Label format: "Nuclei[username] name → template-id"
	if dorks[0].Label != "Nuclei[jdoe42] Twitter Username Found → twitter-username" {
		t.Errorf("dorks[0].Label = %q, unexpected format", dorks[0].Label)
	}
}

func TestParseJSONL_SkipNoURL(t *testing.T) {
	// Line with no matched-at, matched, or host should be skipped
	fixture := []byte(`{"template-id":"empty-template","info":{"name":"No URL","severity":"info"},"matched-at":"","matched":"","host":""}` + "\n")
	dorks := parseJSONL(fixture, "user")
	if len(dorks) != 0 {
		t.Errorf("parseJSONL: got %d dorks for no-URL line, want 0", len(dorks))
	}
}

func TestParseJSONL_InvalidJSON(t *testing.T) {
	fixture := []byte("not json\n{\"template-id\":\"valid\",\"info\":{\"name\":\"Valid\",\"severity\":\"medium\"},\"matched-at\":\"https://example.com\"}\n")
	dorks := parseJSONL(fixture, "user")
	if len(dorks) != 1 {
		t.Errorf("parseJSONL: got %d dorks, want 1 (invalid JSON line skipped)", len(dorks))
	}
}
