# selectronic-exporter

Prometheus exporter for Selectronic / Select.live Solarmon devices. It uses a multi-target probe shape: `/metrics` is for exporter self-metrics, and `/probe` scrapes one configured controller/device pair when Prometheus asks.

## Run locally

```sh
go run ./cmd/selectronic_exporter --web.listen-address=:9788
```

Example probe:

```sh
curl 'http://127.0.0.1:9788/probe?target=http://192.0.2.10&module=default&device_id=ANONDEVICEID000000000000000000'
```

The built-in `default` module uses `/cgi-bin/solarmonweb` and a 5 second timeout. A config file is only needed if that path or timeout needs to change.

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

## Metrics

| Metric                                        | Type    | Labels                                                                                                | Description                                                          |
| --------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| `selectronic_device_info`                     | gauge   | `device_id`, `short_id`, `serial_number`, `manufacturer`, `device_type`, `firmware`, `board_firmware` | Device metadata.                                                     |
| `selectronic_power_rating_watts`              | gauge   |                                                                                                       | Configured inverter power rating.                                    |
| `selectronic_scrape_success`                  | gauge   |                                                                                                       | `1` when the device scrape succeeds, otherwise `0`.                  |
| `selectronic_scrape_duration_seconds`         | gauge   |                                                                                                       | Device scrape duration.                                              |
| `selectronic_point_items`                     | gauge   |                                                                                                       | Number of items in the point response.                               |
| `selectronic_battery_watts`                   | gauge   |                                                                                                       | Instantaneous battery power.                                         |
| `selectronic_battery_state_of_charge_percent` | gauge   |                                                                                                       | Battery state of charge.                                             |
| `selectronic_battery_in_watt_hours_today`     | gauge   |                                                                                                       | Battery charge energy for the current device day.                    |
| `selectronic_battery_in_watt_hours_total`     | counter |                                                                                                       | Lifetime battery charge energy.                                      |
| `selectronic_battery_out_watt_hours_today`    | gauge   |                                                                                                       | Battery discharge energy for the current device day.                 |
| `selectronic_battery_out_watt_hours_total`    | counter |                                                                                                       | Lifetime battery discharge energy.                                   |
| `selectronic_grid_watts`                      | gauge   |                                                                                                       | Instantaneous grid power.                                            |
| `selectronic_grid_in_watt_hours_today`        | gauge   |                                                                                                       | Grid import energy for the current device day.                       |
| `selectronic_grid_in_watt_hours_total`        | counter |                                                                                                       | Lifetime grid import energy.                                         |
| `selectronic_grid_out_watt_hours_today`       | gauge   |                                                                                                       | Grid export energy for the current device day.                       |
| `selectronic_grid_out_watt_hours_total`       | counter |                                                                                                       | Lifetime grid export energy.                                         |
| `selectronic_load_watts`                      | gauge   |                                                                                                       | Instantaneous load power.                                            |
| `selectronic_load_watt_hours_today`           | gauge   |                                                                                                       | Load energy for the current device day.                              |
| `selectronic_load_watt_hours_total`           | counter |                                                                                                       | Lifetime load energy.                                                |
| `selectronic_solar_inverter_watts`            | gauge   |                                                                                                       | Instantaneous solar inverter power.                                  |
| `selectronic_solar_watt_hours_today`          | gauge   |                                                                                                       | Solar energy for the current device day.                             |
| `selectronic_solar_watt_hours_total`          | counter |                                                                                                       | Lifetime solar energy.                                               |
| `selectronic_shunt_watts`                     | gauge   |                                                                                                       | Instantaneous shunt power.                                           |
| `selectronic_generator_status`                | gauge   |                                                                                                       | Numeric generator status.                                            |
| `selectronic_fault_code`                      | gauge   |                                                                                                       | Numeric fault code.                                                  |
| `selectronic_fault_timestamp_seconds`         | gauge   |                                                                                                       | Unix timestamp of the current fault, or `0` when no fault is active. |
| `selectronic_sample_timestamp_seconds`        | gauge   |                                                                                                       | Unix timestamp of the device sample.                                 |
| `selectronic_device_now_timestamp_seconds`    | gauge   |                                                                                                       | Unix timestamp returned by the device API.                           |
