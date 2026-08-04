package exporter

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	promconfig "github.com/prometheus/common/config"

	"github.com/woodleighschool/selectronic-exporter/internal/config"
)

const testDeviceID = "ANONDEVICEID000000000000000000"

func TestProbeRequiresTargetAndDevice(t *testing.T) {
	server := testServer(t)

	for _, path := range []string{
		"/probe",
		"/probe?target=http://example.test",
		"/probe?device_id=" + testDeviceID,
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		server.Probe(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s returned status %d, want 400", path, rec.Code)
		}
	}
}

func TestProbeUnknownModule(t *testing.T) {
	server := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/probe?target=http://example.test&device_id="+testDeviceID+"&module=missing", nil)
	rec := httptest.NewRecorder()
	server.Probe(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestProbeHealthy(t *testing.T) {
	device := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveProbeFixture(t, w, r.URL.Path)
	}))
	defer device.Close()

	server := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/probe?target="+url.QueryEscape(device.URL)+"&device_id="+testDeviceID, nil)
	rec := httptest.NewRecorder()
	server.Probe(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{
		"selectronic_scrape_success 1",
		"selectronic_device_info",
		"selectronic_battery_watts 125.5",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("probe body missing %q:\n%s", want, body)
		}
	}
}

func TestProbeUpstreamFailureIsMetricFailure(t *testing.T) {
	device := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "offline", http.StatusBadGateway)
	}))
	defer device.Close()

	server := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/probe?target="+url.QueryEscape(device.URL)+"&device_id="+testDeviceID, nil)
	rec := httptest.NewRecorder()
	server.Probe(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "selectronic_scrape_success 0") {
		t.Fatalf("probe body missing failure metric:\n%s", rec.Body.String())
	}
}

func testServer(t *testing.T) *Server {
	t.Helper()
	cfg := config.Config{
		Modules: map[string]config.Module{
			"default": {
				Timeout:          5 * time.Second,
				PathPrefix:       "/cgi-bin/solarmonweb",
				HTTPClientConfig: promconfig.DefaultHTTPClientConfig,
			},
		},
	}
	return &Server{Config: cfg, MetricsPath: "/metrics"}
}

func serveProbeFixture(t *testing.T, w http.ResponseWriter, reqPath string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	switch reqPath {
	case "/cgi-bin/solarmonweb/board":
		writeProbeFixture(t, w, "../selectronic/testdata/board.json")
	case "/cgi-bin/solarmonweb/devices/" + testDeviceID:
		writeProbeFixture(t, w, "../selectronic/testdata/device.json")
	case "/cgi-bin/solarmonweb/devices/" + testDeviceID + "/point":
		writeProbeFixture(t, w, "../selectronic/testdata/point.json")
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func writeProbeFixture(t *testing.T, w http.ResponseWriter, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write(data)
}
