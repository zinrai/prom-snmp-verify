# prom-snmp-verify

Verification tool for Prometheus [snmp_exporter](https://github.com/prometheus/snmp_exporter). Detects metric collection regressions when upgrading snmp_exporter or changing snmp.yml configuration.

## Usage

### targets

Fetch scrape targets from Prometheus API.

```
prom-snmp-verify targets --prometheus-url http://prometheus:9090
```

Flags:

- `--prometheus-url` (required) - Prometheus URL
- `--metrics-path` - Filter by metrics_path (default: `/snmp`)
- `--output` - Output file path (default: `targets.json`)

### expectations

Extract expected metric names from snmp.yml.

```
prom-snmp-verify expectations --snmp-yml /path/to/snmp.yml
```

Flags:

- `--snmp-yml` (required) - Path to snmp.yml
- `--output` - Output file path (default: `expectations.json`)

### check

Verify metrics against snmp_exporter by comparing actual scrape results with expected metrics defined in snmp.yml.

```
prom-snmp-verify check \
  --snmp-yml /path/to/snmp.yml \
  --exporter-url http://localhost:9116 \
  --targets targets.json
```

Flags:

- `--snmp-yml` (required) - Path to snmp.yml
- `--exporter-url` (required) - snmp_exporter URL
- `--targets` (required) - Path to targets JSON file (output of `targets` subcommand)
- `--output` - Output file path (default: `check.json`)

Exit codes:

- 0 - All targets responded successfully
- 1 - One or more targets failed

## License

This project is licensed under the [MIT License](./LICENSE).
