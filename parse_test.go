package main

import (
	"strings"
	"testing"
)

func TestParseMetricNames(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]bool
	}{
		{
			name:     "empty input",
			input:    "",
			expected: map[string]bool{},
		},
		{
			name:     "metric without labels",
			input:    "sysUpTime 12345",
			expected: map[string]bool{"sysUpTime": true},
		},
		{
			name:     "metric with labels",
			input:    `ifNumber{instance="192.0.2.1"} 42`,
			expected: map[string]bool{"ifNumber": true},
		},
		{
			name: "skip HELP comment",
			input: `# HELP ifNumber The number of network interfaces
ifNumber{instance="192.0.2.1"} 42`,
			expected: map[string]bool{"ifNumber": true},
		},
		{
			name: "skip TYPE comment",
			input: `# TYPE ifNumber gauge
ifNumber{instance="192.0.2.1"} 42`,
			expected: map[string]bool{"ifNumber": true},
		},
		{
			name: "skip empty lines",
			input: `sysUpTime 12345

ifNumber{instance="192.0.2.1"} 42`,
			expected: map[string]bool{"sysUpTime": true, "ifNumber": true},
		},
		{
			name: "multiple metrics",
			input: `# HELP sysUpTime System uptime
# TYPE sysUpTime gauge
sysUpTime 12345
# HELP ifNumber Number of interfaces
# TYPE ifNumber gauge
ifNumber{instance="192.0.2.1"} 42
# HELP ifDescr Interface description
# TYPE ifDescr gauge
ifDescr{ifIndex="1",instance="192.0.2.1"} 1`,
			expected: map[string]bool{
				"sysUpTime": true,
				"ifNumber":  true,
				"ifDescr":   true,
			},
		},
		{
			name: "duplicate metric names",
			input: `ifDescr{ifIndex="1",instance="192.0.2.1"} 1
ifDescr{ifIndex="2",instance="192.0.2.1"} 1`,
			expected: map[string]bool{"ifDescr": true},
		},
		{
			name:     "metric with multiple labels",
			input:    `ifHCInOctets{ifAlias="",ifDescr="ge1",ifIndex="1",ifName="ge1"} 123456`,
			expected: map[string]bool{"ifHCInOctets": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMetricNames(strings.NewReader(tt.input))

			if len(result) != len(tt.expected) {
				t.Errorf("got %d metric names, want %d\ngot:  %v\nwant: %v", len(result), len(tt.expected), result, tt.expected)
				return
			}

			for name := range tt.expected {
				if !result[name] {
					t.Errorf("missing expected metric name %q\ngot: %v", name, result)
				}
			}
		})
	}
}
