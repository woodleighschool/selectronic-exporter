package selectronic

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCollectorHealthyScrape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveFixture(t, w, r.URL.Path)
	}))
	defer server.Close()

	collector := &Collector{
		Client:   newTestClient(t, server.URL),
		DeviceID: testDeviceID,
		Timeout:  time.Second,
	}

	expected := `
# HELP selectronic_battery_in_watt_hours_total Lifetime battery charge energy total reported by the device.
# TYPE selectronic_battery_in_watt_hours_total counter
selectronic_battery_in_watt_hours_total 1000.25
# HELP selectronic_battery_state_of_charge_percent Battery state of charge reported by the device.
# TYPE selectronic_battery_state_of_charge_percent gauge
selectronic_battery_state_of_charge_percent 88.5
# HELP selectronic_battery_watts Instantaneous battery power reported by the device.
# TYPE selectronic_battery_watts gauge
selectronic_battery_watts 125.5
# HELP selectronic_device_info Selectronic device information.
# TYPE selectronic_device_info gauge
selectronic_device_info{board_firmware="4.2.1",device_id="ANONDEVICEID000000000000000000",device_type="SP-PRO",firmware="V15.47 Grid 4.00 ACC 5.11",manufacturer="Selectronic",serial_number="ANONSERIAL",short_id="ANONDEV1"} 1
# HELP selectronic_power_rating_watts Configured inverter power rating reported by the device.
# TYPE selectronic_power_rating_watts gauge
selectronic_power_rating_watts 5000
# HELP selectronic_scrape_success Whether the Selectronic scrape completed successfully.
# TYPE selectronic_scrape_success gauge
selectronic_scrape_success 1
`
	if err := testutil.CollectAndCompare(
		collector,
		strings.NewReader(expected),
		"selectronic_battery_in_watt_hours_total",
		"selectronic_battery_state_of_charge_percent",
		"selectronic_battery_watts",
		"selectronic_device_info",
		"selectronic_power_rating_watts",
		"selectronic_scrape_success",
	); err != nil {
		t.Fatal(err)
	}
}

func TestCollectorFailedScrape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "offline", http.StatusBadGateway)
	}))
	defer server.Close()

	collector := &Collector{
		Client:   newTestClient(t, server.URL),
		DeviceID: testDeviceID,
		Timeout:  time.Second,
	}

	expected := `
# HELP selectronic_scrape_success Whether the Selectronic scrape completed successfully.
# TYPE selectronic_scrape_success gauge
selectronic_scrape_success 0
`
	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected), "selectronic_scrape_success"); err != nil {
		t.Fatal(err)
	}
}
