package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_DefaultOutput(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-n", "John Doe")
	cmd.Dir = "." // run from cmd/dork/
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Error("no output")
	}
	if !strings.Contains(string(out), "social") {
		t.Error("missing social category")
	}
}

func TestCLI_JSONExport(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-n", "John Doe", "--export", "json")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, out)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Errorf("not valid JSON: %v", err)
	}
}

func TestCLI_NoNameError(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "."
	err := cmd.Run()
	if err == nil {
		t.Error("should fail without --name")
	}
}

func TestCLI_CaseFile(t *testing.T) {
	dir := t.TempDir()
	caseFile := filepath.Join(dir, "test.yaml")
	os.WriteFile(caseFile, []byte("name: \"Test User\"\nlocation: \"NYC\"\n"), 0644)

	cmd := exec.Command("go", "run", ".", "--case", caseFile, "--export", "json")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Test User") {
		t.Error("missing case file name")
	}
}

func TestCLI_RegionFilter(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-n", "John Doe", "--region", "us", "--export", "json")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Spokeo") {
		t.Error("us region should include Spokeo")
	}
}
