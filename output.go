package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// writeJSON encodes v as indented JSON and writes it to the specified file path.
func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	return nil
}

// jsonLinesWriter manages a file for writing JSON Lines (one JSON object per line).
type jsonLinesWriter struct {
	file *os.File
	enc  *json.Encoder
}

// newJSONLinesWriter creates a new file at path (removing any existing file) and
// returns a writer for appending JSON Lines.
func newJSONLinesWriter(path string) (*jsonLinesWriter, error) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("removing existing file: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("creating output file: %w", err)
	}

	return &jsonLinesWriter{
		file: f,
		enc:  json.NewEncoder(f),
	}, nil
}

// Write encodes v as a single JSON line and writes it to the file.
func (w *jsonLinesWriter) Write(v any) error {
	if err := w.enc.Encode(v); err != nil {
		return fmt.Errorf("encoding JSON line: %w", err)
	}
	return nil
}

// Close closes the underlying file.
func (w *jsonLinesWriter) Close() error {
	return w.file.Close()
}
