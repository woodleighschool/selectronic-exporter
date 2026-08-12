# selectronic-exporter

Prometheus multi-target exporter for Selectronic and Select.live Solarmon devices. `/metrics` serves exporter self-metrics; `/probe` reads one controller and device pair on demand.

## 🚀 Usage

Run the exporter, then probe a controller:

```bash
go run ./cmd/selectronic_exporter --web.listen-address=:9788
curl 'http://127.0.0.1:9788/probe?target=http://192.0.2.10&module=default&device_id=ANONDEVICEID000000000000000000'
```

A Prometheus scrape job looks like this:

```yaml
scrape_configs:
    - job_name: selectronic
      metrics_path: /probe
      params:
          module: [default]
      static_configs:
          - targets: ["http://192.0.2.10"]
            labels:
                device_id: "ANONDEVICEID000000000000000000"
      relabel_configs:
          - source_labels: [__address__]
            target_label: __param_target
          - source_labels: [device_id]
            target_label: __param_device_id
          - source_labels: [device_id]
            target_label: instance
          - target_label: __address__
            replacement: selectronic-exporter:9788
```

## ⚙️ Configuration

The built-in `default` module uses `/cgi-bin/solarmonweb`, collects faults, and times out after five seconds. Start from [`config.example.yml`](config.example.yml) only when those settings or HTTP behaviour need to change.

```bash
go run ./cmd/selectronic_exporter --config.file=config.yml --config.check
```

| Flag                   | Default         | Purpose                         |
| ---------------------- | --------------- | ------------------------------- |
| `--web.listen-address` | `:9788`         | HTTP listen address             |
| `--web.telemetry-path` | `/metrics`      | Self-metrics path               |
| `--config.file`        | Built-in module | YAML module configuration       |
| `--config.check`       | `false`         | Validate configuration and exit |

## 📈 Metrics

| Family                                        | Purpose                                                   |
| --------------------------------------------- | --------------------------------------------------------- |
| `selectronic_device_info`                     | Device, firmware, and hardware metadata                   |
| `selectronic_*_watts`                         | Instantaneous battery, grid, load, solar, and shunt power |
| `selectronic_*_watt_hours_today`              | Energy for the current device day                         |
| `selectronic_*_watt_hours_total`              | Lifetime energy counters                                  |
| `selectronic_battery_state_of_charge_percent` | Battery state of charge                                   |
| `selectronic_generator_status`                | Generator status                                          |
| `selectronic_fault_*`                         | Active fault code and timestamp                           |
| `selectronic_sample_timestamp_seconds`        | Device sample timestamp                                   |
| `selectronic_scrape_success`                  | Whether the target scrape succeeded                       |
| `selectronic_scrape_duration_seconds`         | Target scrape duration                                    |

## 🧑‍💻 Development

```bash
mise run deps
mise run build
mise run test
mise run lint
mise run fmt-check
```

Tests use captured responses and local servers; no controller is required.

## 📄 License

Licensed under the [Apache License 2.0](LICENSE).
