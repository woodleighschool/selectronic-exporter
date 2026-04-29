# selectronic-exporter

Prometheus exporter for Selectronic / Select.live Solarmon devices. It uses a multi-target probe shape: `/metrics` is for exporter self-metrics, and `/probe` scrapes one configured controller/device pair when Prometheus asks.

## Run locally

```sh
go run ./cmd/selectronic_exporter --config.file=config.example.yml --web.listen-address=:9788
```

Example probe:

```sh
curl 'http://127.0.0.1:9788/probe?target=http://192.0.2.10&module=default&device_id=ANONDEVICEID000000000000000000'
```

## Prometheus

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
          site: "example"
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
