package main

import (
	"testing"
)

func TestParseSnmpConfig(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantModules map[string][]string
		wantErr     bool
	}{
		{
			name: "wrapped format with multiple modules",
			input: `
modules:
  module_a:
    walk:
    - 1.3.6.1.2.1.2
    metrics:
    - name: metric_one
      oid: 1.3.6.1.2.1.2.1
      type: gauge
    - name: metric_two
      oid: 1.3.6.1.2.1.2.2
      type: gauge
  module_b:
    walk:
    - 1.3.6.1.2.1.1
    metrics:
    - name: metric_three
      oid: 1.3.6.1.2.1.1.3
      type: counter
`,
			wantModules: map[string][]string{
				"module_a": {"metric_one", "metric_two"},
				"module_b": {"metric_three"},
			},
		},
		{
			name: "wrapped format with single module",
			input: `
modules:
  module_a:
    get:
    - 1.3.6.1.2.1.1.1
    metrics:
    - name: metric_one
      oid: 1.3.6.1.2.1.1.1
      type: counter
`,
			wantModules: map[string][]string{
				"module_a": {"metric_one"},
			},
		},
		{
			name: "flat format with multiple modules",
			input: `
module_a:
  walk:
  - 1.3.6.1.2.1.2
  metrics:
  - name: metric_one
    oid: 1.3.6.1.2.1.2.1
    type: gauge
  - name: metric_two
    oid: 1.3.6.1.2.1.2.2
    type: gauge
  version: 2
  auth:
    community: public
module_b:
  walk:
  - 1.3.6.1.2.1.1
  metrics:
  - name: metric_three
    oid: 1.3.6.1.2.1.1.3
    type: counter
  version: 2
  auth:
    community: public
`,
			wantModules: map[string][]string{
				"module_a": {"metric_one", "metric_two"},
				"module_b": {"metric_three"},
			},
		},
		{
			name: "flat format with single module",
			input: `
module_a:
  get:
  - 1.3.6.1.2.1.1.1
  metrics:
  - name: metric_one
    oid: 1.3.6.1.2.1.1.1
    type: counter
  version: 2
  auth:
    community: public
`,
			wantModules: map[string][]string{
				"module_a": {"metric_one"},
			},
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "valid yaml but no metrics",
			input:   "key: value\n",
			wantErr: true,
		},
		{
			name: "wrapped format with empty modules",
			input: `
modules: {}
`,
			wantErr: true,
		},
		{
			name: "wrapped format with module without metrics",
			input: `
modules:
  module_a:
    walk:
    - 1.3.6.1.2.1.2
`,
			wantErr: true,
		},
		{
			name: "flat format with module without metrics",
			input: `
module_a:
  walk:
  - 1.3.6.1.2.1.2
  version: 2
`,
			wantErr: true,
		},
		{
			name: "wrapped format with metric without name",
			input: `
modules:
  module_a:
    metrics:
    - oid: 1.3.6.1.2.1.2.1
      type: gauge
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules, err := parseSnmpConfig([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil, modules: %v", modules)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(modules) != len(tt.wantModules) {
				t.Errorf("got %d modules, want %d\ngot:  %v", len(modules), len(tt.wantModules), moduleNames(modules))
			}

			for name, wantMetrics := range tt.wantModules {
				mod, ok := modules[name]
				if !ok {
					t.Errorf("missing expected module %q", name)
					continue
				}

				var gotNames []string
				for _, m := range mod.Metrics {
					if m.Name != "" {
						gotNames = append(gotNames, m.Name)
					}
				}

				if len(gotNames) != len(wantMetrics) {
					t.Errorf("module %q: got %d metrics, want %d\ngot:  %v\nwant: %v", name, len(gotNames), len(wantMetrics), gotNames, wantMetrics)
					continue
				}

				for i, want := range wantMetrics {
					if gotNames[i] != want {
						t.Errorf("module %q: metric[%d] = %q, want %q", name, i, gotNames[i], want)
					}
				}
			}
		})
	}
}

func moduleNames(modules map[string]snmpModule) []string {
	var names []string
	for name := range modules {
		names = append(names, name)
	}
	return names
}
