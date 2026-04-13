package main

import (
	"bufio"
	"io"
	"strings"
)

// parseMetricNames extracts unique metric names from Prometheus text exposition format.
// It skips comment lines (starting with #) and empty lines, then extracts the metric
// name from the beginning of each line up to the first '{' or ' '.
func parseMetricNames(r io.Reader) map[string]bool {
	names := make(map[string]bool)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name := line
		if i := strings.IndexByte(line, '{'); i != -1 {
			name = line[:i]
		} else if i := strings.IndexByte(line, ' '); i != -1 {
			name = line[:i]
		}

		if name != "" {
			names[name] = true
		}
	}

	return names
}
