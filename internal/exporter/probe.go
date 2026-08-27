package exporter

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/woodleighschool/selectronic-exporter/internal/config"
	"github.com/woodleighschool/selectronic-exporter/internal/selectronic"
)

// NewProbeHandler returns a handler that collects metrics from one configured controller target.
func NewProbeHandler(cfg config.Config, logger *slog.Logger) http.Handler {
	return &probeHandler{config: cfg, logger: logger}
}

type probeHandler struct {
	config config.Config
	logger *slog.Logger
}

func (h *probeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	deviceID := r.URL.Query().Get("device_id")
	moduleName := r.URL.Query().Get("module")
	if moduleName == "" {
		moduleName = "default"
	}

	module, ok := h.config.Modules[moduleName]
	if !ok {
		http.Error(w, fmt.Sprintf("unknown module %q", moduleName), http.StatusBadRequest)
		return
	}
	if target == "" || deviceID == "" {
		http.Error(w, "target and device_id are required", http.StatusBadRequest)
		return
	}

	client, err := selectronic.NewClient(target, module.PathPrefix, module.HTTPClientConfig, h.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(&selectronic.Collector{
		Client:   client,
		DeviceID: deviceID,
		Timeout:  module.Timeout,
		Logger:   h.logger,
	})

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
