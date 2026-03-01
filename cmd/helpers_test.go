package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestOutputJSON(t *testing.T) {
	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := map[string]string{"key": "value"}
	if err := outputJSON(data); err != nil {
		t.Fatalf("outputJSON error: %v", err)
	}

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("key = %q, want %q", result["key"], "value")
	}
}

func TestOutputError(t *testing.T) {
	// Redirect stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	outputError(os.ErrNotExist)

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result["error"] == "" {
		t.Error("error field should not be empty")
	}
}
