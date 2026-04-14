package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	snmpYml := fs.String("snmp-yml", "", "Path to snmp.yml (required)")
	exporterURL := fs.String("exporter-url", "", "snmp_exporter URL (required)")
	targetsFile := fs.String("targets", "", "Path to targets JSON file (required)")
	output := fs.String("output", "check.json", "Output file path")
	fs.Parse(args)

	if *snmpYml == "" || *exporterURL == "" || *targetsFile == "" {
		return fmt.Errorf("--snmp-yml, --exporter-url, and --targets are all required")
	}

	targets, err := loadTargets(*targetsFile)
	if err != nil {
		return err
	}

	expectations, err := loadExpectations(*snmpYml)
	if err != nil {
		return err
	}

	results, hasError := checkAll(*exporterURL, targets, expectations)

	if err := writeJSON(*output, results); err != nil {
		return err
	}

	if hasError {
		os.Exit(1)
	}

	return nil
}

func loadTargets(path string) ([]Target, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading targets file: %w", err)
	}

	var targets []Target
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, fmt.Errorf("parsing targets file: %w", err)
	}

	return targets, nil
}

func checkAll(exporterURL string, targets []Target, expectations map[string][]string) ([]CheckResult, bool) {
	var results []CheckResult
	hasError := false

	for _, t := range targets {
		result := checkTarget(exporterURL, t, expectations)
		if result.Status == "error" {
			hasError = true
		}
		results = append(results, result)
	}

	return results, hasError
}

func checkTarget(exporterURL string, t Target, expectations map[string][]string) CheckResult {
	scraped, err := url.Parse(t.ScrapeURL)
	if err != nil {
		return CheckResult{
			ScrapePool: t.ScrapePool,
			Status:     "error",
			Error:      fmt.Sprintf("parsing scrapeUrl: %v", err),
		}
	}

	params := scraped.Query()
	target := params.Get("target")
	module := params.Get("module")

	expected, ok := expectations[module]
	if !ok {
		return CheckResult{
			ScrapePool: t.ScrapePool,
			Target:     target,
			Module:     module,
			Status:     "error",
			Error:      fmt.Sprintf("module %q not found in snmp.yml", module),
		}
	}

	actual, err := scrapeMetrics(exporterURL, scraped.Query())
	if err != nil {
		return CheckResult{
			ScrapePool: t.ScrapePool,
			Target:     target,
			Module:     module,
			Status:     "error",
			Error:      err.Error(),
		}
	}

	var okNames []string
	var missing []string
	for _, name := range expected {
		if actual[name] {
			okNames = append(okNames, name)
		} else {
			missing = append(missing, name)
		}
	}

	status := "pass"
	if len(missing) > 0 {
		status = "warn"
	}

	return CheckResult{
		ScrapePool: t.ScrapePool,
		Target:     target,
		Module:     module,
		OK:         okNames,
		Missing:    missing,
		Status:     status,
	}
}

func scrapeMetrics(exporterURL string, params url.Values) (map[string]bool, error) {
	base, err := url.Parse(exporterURL)
	if err != nil {
		return nil, fmt.Errorf("parsing exporter URL: %w", err)
	}
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return nil, fmt.Errorf("requesting snmp_exporter: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("snmp_exporter returned status %d", resp.StatusCode)
	}

	return parseMetricNames(resp.Body), nil
}
