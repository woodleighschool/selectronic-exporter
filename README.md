# selectronic-exporter

[![Release](https://img.shields.io/github/v/release/woodleighschool/selectronic-exporter?display_name=tag&sort=semver)](https://github.com/woodleighschool/selectronic-exporter/releases/latest)
[![CI](https://github.com/woodleighschool/selectronic-exporter/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/woodleighschool/selectronic-exporter/actions/workflows/ci.yaml)
[![Go](https://img.shields.io/github/go-mod/go-version/woodleighschool/selectronic-exporter?logo=go)](https://github.com/woodleighschool/selectronic-exporter/blob/main/go.mod)
[![Container](https://img.shields.io/badge/container-ghcr.io-2496ED?logo=github&logoColor=white)](https://github.com/orgs/woodleighschool/packages/container/package/selectronic-exporter)
[![License](https://img.shields.io/github/license/woodleighschool/selectronic-exporter)](https://github.com/woodleighschool/selectronic-exporter/blob/main/LICENSE)

Prometheus multi-target exporter for Selectronic and Select.live Solarmon devices. `/metrics` serves exporter self-metrics; `/probe` reads one controller and device pair on demand.

## 🚀 Usage

A container is published with each [release](https://github.com/woodleighschool/selectronic-exporter/releases/latest):

```bash
docker run --rm \
  --publish 9788:9788 \
  ghcr.io/woodleighschool/selectronic-exporter:rolling
```

Probe a controller:

```bash
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
docker run --rm \
  --volume "$PWD/config.yml:/config.yml:ro" \
  ghcr.io/woodleighschool/selectronic-exporter:rolling \
  --config.file=/config.yml \
  --config.check
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

Run the current checkout directly:

```bash
go run ./cmd/selectronic_exporter --web.listen-address=:9788
```

Repository checks:

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
