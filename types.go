package main

// Target represents an snmp_exporter scrape target.
type Target struct {
	ScrapePool string `json:"scrapePool"`
	ScrapeURL  string `json:"scrapeUrl"`
}

// CheckResult represents the verification result for a single target.
type CheckResult struct {
	ScrapePool string   `json:"scrapePool"`
	Target     string   `json:"target"`
	Module     string   `json:"module"`
	OK         []string `json:"ok"`
	Missing    []string `json:"missing"`
	Status     string   `json:"status"`
	Error      string   `json:"error,omitempty"`
}
