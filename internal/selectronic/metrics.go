package selectronic

import "github.com/prometheus/client_golang/prometheus"

var (
	deviceInfoDesc = prometheus.NewDesc(
		"selectronic_device_info",
		"Selectronic device information.",
		[]string{"device_id", "short_id", "serial_number", "manufacturer", "device_type", "firmware", "board_firmware"},
		nil,
	)
	powerRatingDesc = prometheus.NewDesc(
		"selectronic_power_rating_watts",
		"Configured inverter power rating reported by the device.",
		nil,
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		"selectronic_scrape_success",
		"Whether the Selectronic scrape completed successfully.",
		nil,
		nil,
	)
	scrapeDurationDesc = prometheus.NewDesc(
		"selectronic_scrape_duration_seconds",
		"Duration of the Selectronic scrape.",
		nil,
		nil,
	)
	pointItemsDesc = prometheus.NewDesc(
		"selectronic_point_items",
		"Number of items reported in the point response.",
		nil,
		nil,
	)
	batteryWattsDesc = prometheus.NewDesc(
		"selectronic_battery_watts",
		"Instantaneous battery power reported by the device.",
		nil,
		nil,
	)
	batterySOCDesc = prometheus.NewDesc(
		"selectronic_battery_state_of_charge_percent",
		"Battery state of charge reported by the device.",
		nil,
		nil,
	)
	batteryInTodayDesc = prometheus.NewDesc(
		"selectronic_battery_in_watt_hours_today",
		"Battery charge energy reported for the current device day.",
		nil,
		nil,
	)
	batteryInTotalDesc = prometheus.NewDesc(
		"selectronic_battery_in_watt_hours_total",
		"Lifetime battery charge energy total reported by the device.",
		nil,
		nil,
	)
	batteryOutTodayDesc = prometheus.NewDesc(
		"selectronic_battery_out_watt_hours_today",
		"Battery discharge energy reported for the current device day.",
		nil,
		nil,
	)
	batteryOutTotalDesc = prometheus.NewDesc(
		"selectronic_battery_out_watt_hours_total",
		"Lifetime battery discharge energy total reported by the device.",
		nil,
		nil,
	)
	gridWattsDesc = prometheus.NewDesc(
		"selectronic_grid_watts",
		"Instantaneous grid power reported by the device.",
		nil,
		nil,
	)
	gridInTodayDesc = prometheus.NewDesc(
		"selectronic_grid_in_watt_hours_today",
		"Grid import energy reported for the current device day.",
		nil,
		nil,
	)
	gridInTotalDesc = prometheus.NewDesc(
		"selectronic_grid_in_watt_hours_total",
		"Lifetime grid import energy total reported by the device.",
		nil,
		nil,
	)
	gridOutTodayDesc = prometheus.NewDesc(
		"selectronic_grid_out_watt_hours_today",
		"Grid export energy reported for the current device day.",
		nil,
		nil,
	)
	gridOutTotalDesc = prometheus.NewDesc(
		"selectronic_grid_out_watt_hours_total",
		"Lifetime grid export energy total reported by the device.",
		nil,
		nil,
	)
	loadWattsDesc = prometheus.NewDesc(
		"selectronic_load_watts",
		"Instantaneous load power reported by the device.",
		nil,
		nil,
	)
	loadTodayDesc = prometheus.NewDesc(
		"selectronic_load_watt_hours_today",
		"Load energy reported for the current device day.",
		nil,
		nil,
	)
	loadTotalDesc = prometheus.NewDesc(
		"selectronic_load_watt_hours_total",
		"Lifetime load energy total reported by the device.",
		nil,
		nil,
	)
	solarInverterWattsDesc = prometheus.NewDesc(
		"selectronic_solar_inverter_watts",
		"Instantaneous solar inverter power reported by the device.",
		nil,
		nil,
	)
	solarTodayDesc = prometheus.NewDesc(
		"selectronic_solar_watt_hours_today",
		"Solar energy reported for the current device day.",
		nil,
		nil,
	)
	solarTotalDesc = prometheus.NewDesc(
		"selectronic_solar_watt_hours_total",
		"Lifetime solar energy total reported by the device.",
		nil,
		nil,
	)
	shuntWattsDesc = prometheus.NewDesc(
		"selectronic_shunt_watts",
		"Instantaneous shunt power reported by the device.",
		nil,
		nil,
	)
	generatorStatusDesc = prometheus.NewDesc(
		"selectronic_generator_status",
		"Numeric generator status reported by the device.",
		nil,
		nil,
	)
	faultCodeDesc = prometheus.NewDesc(
		"selectronic_fault_code",
		"Numeric fault code reported by the device.",
		nil,
		nil,
	)
	faultTimestampDesc = prometheus.NewDesc(
		"selectronic_fault_timestamp_seconds",
		"Unix timestamp of the current fault, or 0 when no fault is active.",
		nil,
		nil,
	)
	sampleTimestampDesc = prometheus.NewDesc(
		"selectronic_sample_timestamp_seconds",
		"Unix timestamp of the device sample.",
		nil,
		nil,
	)
	deviceNowTimestampDesc = prometheus.NewDesc(
		"selectronic_device_now_timestamp_seconds",
		"Unix timestamp returned by the device API.",
		nil,
		nil,
	)
)

var allDescs = []*prometheus.Desc{
	deviceInfoDesc,
	powerRatingDesc,
	scrapeSuccessDesc,
	scrapeDurationDesc,
	pointItemsDesc,
	batteryWattsDesc,
	batterySOCDesc,
	batteryInTodayDesc,
	batteryInTotalDesc,
	batteryOutTodayDesc,
	batteryOutTotalDesc,
	gridWattsDesc,
	gridInTodayDesc,
	gridInTotalDesc,
	gridOutTodayDesc,
	gridOutTotalDesc,
	loadWattsDesc,
	loadTodayDesc,
	loadTotalDesc,
	solarInverterWattsDesc,
	solarTodayDesc,
	solarTotalDesc,
	shuntWattsDesc,
	generatorStatusDesc,
	faultCodeDesc,
	faultTimestampDesc,
	sampleTimestampDesc,
	deviceNowTimestampDesc,
}
