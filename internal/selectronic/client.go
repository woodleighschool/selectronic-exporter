package selectronic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"

	promconfig "github.com/prometheus/common/config"
)

const userAgent = "selectronic_exporter"

type Client struct {
	baseURL    *url.URL
	pathPrefix string
	httpClient *http.Client
	logger     *slog.Logger
}

func NewClient(target, pathPrefix string, httpClientConfig promconfig.HTTPClientConfig, logger *slog.Logger) (*Client, error) {
	baseURL, err := url.Parse(strings.TrimRight(target, "/"))
	if err != nil {
		return nil, err
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("target must include scheme and host")
	}

	httpClient, err := promconfig.NewClientFromConfig(httpClientConfig, userAgent, promconfig.WithUserAgent(userAgent))
	if err != nil {
		return nil, err
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &Client{
		baseURL:    baseURL,
		pathPrefix: "/" + strings.Trim(strings.TrimSpace(pathPrefix), "/"),
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

func (c *Client) Scrape(ctx context.Context, deviceID string) (Snapshot, error) {
	board, err := getJSON[Board](ctx, c, "board/")
	if err != nil {
		return Snapshot{}, err
	}

	device, err := getJSON[Device](ctx, c, "devices", deviceID)
	if err != nil {
		return Snapshot{}, err
	}

	// /usage is advertised by the device metadata, but the real controller timed
	// out without bytes during bootstrap probing. /point is the v1 source of truth.
	point, err := getJSON[Point](ctx, c, "devices", deviceID, "point")
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		Board:  board,
		Device: device,
		Point:  point,
	}, nil
}

func getJSON[T any](ctx context.Context, c *Client, parts ...string) (T, error) {
	var result T

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint(parts...).String(), nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode/100 != 2 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return result, fmt.Errorf("GET %s returned %s", req.URL.String(), resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("decode %s: %w", req.URL.String(), err)
	}

	return result, nil
}

func (c *Client) endpoint(parts ...string) *url.URL {
	u := *c.baseURL
	allParts := append([]string{strings.Trim(c.pathPrefix, "/")}, parts...)
	u.Path = path.Join(allParts...)
	if !strings.HasPrefix(u.Path, "/") {
		u.Path = "/" + u.Path
	}
	return &u
}
