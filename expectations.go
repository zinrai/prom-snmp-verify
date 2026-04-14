package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

func runExpectations(args []string) error {
	fs := flag.NewFlagSet("expectations", flag.ExitOnError)
	snmpYml := fs.String("snmp-yml", "", "Path to snmp.yml (required)")
	output := fs.String("output", "expectations.json", "Output file path")
	fs.Parse(args)

	if *snmpYml == "" {
		return fmt.Errorf("--snmp-yml is required")
	}

	expectations, err := loadExpectations(*snmpYml)
	if err != nil {
		return err
	}

	for name, metrics := range expectations {
		slog.Info("module found", "module", name, "metrics", len(metrics))
	}
	slog.Info("loaded expectations", "modules", len(expectations), "output", *output)

	return writeJSON(*output, expectations)
}

// loadExpectations parses snmp.yml and returns a map of module name to metric names.
// It supports two snmp.yml formats:
//   - wrapped: modules are nested under a top-level "modules" key (snmp_exporter >= v0.26.0)
//   - flat: modules are defined directly at the top level (snmp_exporter <= v0.25.0)
func loadExpectations(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snmp.yml: %w", err)
	}

	modules, err := parseSnmpConfig(data)
	if err != nil {
		return nil, fmt.Errorf("parsing snmp.yml: %w", err)
	}

	result := make(map[string][]string)
	for name, mod := range modules {
		var names []string
		for _, m := range mod.Metrics {
			if m.Name != "" {
				names = append(names, m.Name)
			}
		}
		result[name] = names
	}

	return result, nil
}

// parseSnmpConfig tries both snmp.yml formats and returns the parsed modules.
func parseSnmpConfig(data []byte) (map[string]snmpModule, error) {
	if modules, ok := tryParseWrapped(data); ok {
		return modules, nil
	}

	if modules, ok := tryParseFlat(data); ok {
		return modules, nil
	}

	return nil, fmt.Errorf("unrecognized snmp.yml format: expected either a top-level 'modules' key or top-level module definitions with 'metrics' fields")
}

// snmpConfigWrapped represents snmp.yml with a top-level "modules" key (>= v0.26.0).
type snmpConfigWrapped struct {
	Modules map[string]snmpModule `yaml:"modules"`
}

// snmpModule represents a single module definition in snmp.yml.
type snmpModule struct {
	Metrics []snmpMetric `yaml:"metrics"`
	Walk    []string     `yaml:"walk"`
	Get     []string     `yaml:"get"`
	Version int          `yaml:"version"`
}

// snmpMetric represents a single metric definition within a module.
type snmpMetric struct {
	Name string `yaml:"name"`
}

// tryParseWrapped attempts to parse data as the wrapped format (modules under "modules" key).
// Returns the modules and true if the format matches, or nil and false otherwise.
func tryParseWrapped(data []byte) (map[string]snmpModule, bool) {
	var config snmpConfigWrapped
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, false
	}

	if !validateModules(config.Modules) {
		return nil, false
	}

	return config.Modules, true
}

// tryParseFlat attempts to parse data as the flat format (modules at top level).
// Returns the modules and true if the format matches, or nil and false otherwise.
func tryParseFlat(data []byte) (map[string]snmpModule, bool) {
	var modules map[string]snmpModule
	if err := yaml.Unmarshal(data, &modules); err != nil {
		return nil, false
	}

	if !validateModules(modules) {
		return nil, false
	}

	return modules, true
}

// validateModules checks that the parsed modules map is non-empty and
// at least one module contains a metric with a non-empty name.
func validateModules(modules map[string]snmpModule) bool {
	if len(modules) == 0 {
		return false
	}

	for _, mod := range modules {
		if hasNamedMetric(mod) {
			return true
		}
	}

	return false
}

// hasNamedMetric reports whether the module contains at least one metric with a non-empty name.
func hasNamedMetric(mod snmpModule) bool {
	for _, m := range mod.Metrics {
		if m.Name != "" {
			return true
		}
	}

	return false
}
