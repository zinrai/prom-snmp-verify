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
