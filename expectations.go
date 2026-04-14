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
func loadExpectations(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snmp.yml: %w", err)
	}

	var config snmpConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing snmp.yml: %w", err)
	}

	result := make(map[string][]string)
	for name, mod := range config.Modules {
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

type snmpConfig struct {
	Modules map[string]snmpModule `yaml:"modules"`
}

type snmpModule struct {
	Metrics []snmpMetric `yaml:"metrics"`
}

type snmpMetric struct {
	Name string `yaml:"name"`
}
