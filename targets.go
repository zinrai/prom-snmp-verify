package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"encoding/json"
)

func runTargets(args []string) error {
	fs := flag.NewFlagSet("targets", flag.ExitOnError)
	prometheusURL := fs.String("prometheus-url", "", "Prometheus URL (required)")
	metricsPath := fs.String("metrics-path", "/snmp", "Filter by metrics_path")
	output := fs.String("output", "targets.json", "Output file path")
	fs.Parse(args)

	if *prometheusURL == "" {
		return fmt.Errorf("--prometheus-url is required")
	}

	targets, err := fetchTargets(*prometheusURL, *metricsPath)
	if err != nil {
		return err
	}

	return writeJSON(*output, targets)
}

func fetchTargets(prometheusURL, metricsPath string) ([]Target, error) {
	u, err := url.JoinPath(prometheusURL, "/api/v1/targets")
	if err != nil {
		return nil, fmt.Errorf("building URL: %w", err)
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("requesting Prometheus API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Prometheus API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var apiResp prometheusTargetsResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	var targets []Target
	for _, at := range apiResp.Data.ActiveTargets {
		if at.DiscoveredLabels["__metrics_path__"] != metricsPath {
			continue
		}

		if at.ScrapeURL == "" {
			continue
		}

		targets = append(targets, Target{
			ScrapePool: at.ScrapePool,
			ScrapeURL:  at.ScrapeURL,
		})
	}

	return targets, nil
}

type prometheusTargetsResponse struct {
	Data struct {
		ActiveTargets []prometheusActiveTarget `json:"activeTargets"`
	} `json:"data"`
}

type prometheusActiveTarget struct {
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
	ScrapePool       string            `json:"scrapePool"`
	ScrapeURL        string            `json:"scrapeUrl"`
}
