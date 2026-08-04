package selectronic

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	promconfig "github.com/prometheus/common/config"
)

const testDeviceID = "ANONDEVICEID000000000000000000"

func TestClientScrape(t *testing.T) {
	seen := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen[r.URL.Path] = true
		if got := r.Header.Get("User-Agent"); got != userAgent {
			t.Fatalf("User-Agent = %q, want %q", got, userAgent)
		}
		serveFixture(t, w, r.URL.Path)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	snapshot, err := client.Scrape(context.Background(), testDeviceID)
	if err != nil {
		t.Fatalf("Scrape returned error: %v", err)
	}

	for _, path := range []string{
		"/cgi-bin/solarmonweb/board",
		"/cgi-bin/solarmonweb/devices/" + testDeviceID,
		"/cgi-bin/solarmonweb/devices/" + testDeviceID + "/point",
	} {
		if !seen[path] {
			t.Fatalf("expected request for %s", path)
		}
	}
	if snapshot.Device.ID != testDeviceID {
		t.Fatalf("Device.ID = %q, want test device", snapshot.Device.ID)
	}
	if snapshot.Point.Items.BatterySOC != 88.5 {
		t.Fatalf("BatterySOC = %f, want 88.5", snapshot.Point.Items.BatterySOC)
	}
}

func TestClientNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusBadGateway)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	if _, err := client.Scrape(context.Background(), testDeviceID); err == nil {
		t.Fatal("Scrape succeeded after non-2xx response")
	}
}

func TestClientInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	if _, err := client.Scrape(context.Background(), testDeviceID); err == nil {
		t.Fatal("Scrape succeeded after invalid JSON")
	}
}

func TestClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, err := client.Scrape(ctx, testDeviceID); err == nil {
		t.Fatal("Scrape succeeded after context timeout")
	}
}

func TestNewClientRequiresHTTPURL(t *testing.T) {
	if _, err := NewClient("controller.local", "/cgi-bin/solarmonweb", promconfig.DefaultHTTPClientConfig, nil); err == nil {
		t.Fatal("NewClient accepted target without scheme")
	}
}

func newTestClient(t *testing.T, target string) *Client {
	t.Helper()
	client, err := NewClient(target, "/cgi-bin/solarmonweb", promconfig.DefaultHTTPClientConfig, nil)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	return client
}

func serveFixture(t *testing.T, w http.ResponseWriter, reqPath string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	switch reqPath {
	case "/cgi-bin/solarmonweb/board":
		writeFixture(t, w, "board.json")
	case "/cgi-bin/solarmonweb/devices/" + testDeviceID:
		writeFixture(t, w, "device.json")
	case "/cgi-bin/solarmonweb/devices/" + testDeviceID + "/point":
		writeFixture(t, w, "point.json")
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func writeFixture(t *testing.T, w http.ResponseWriter, name string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write(data)
}
