package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "targets":
		err = runTargets(os.Args[2:])
	case "expectations":
		err = runExpectations(os.Args[2:])
	case "check":
		err = runCheck(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: prom-snmp-verify <command> [flags]

Commands:
  targets        Fetch scrape targets from Prometheus API
  expectations   Extract expected metric names from snmp.yml
  check          Verify metrics against snmp_exporter
`)
}
