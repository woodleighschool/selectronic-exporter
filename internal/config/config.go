package config

import (
	"errors"
	"os"
	"strings"
	"time"

	promconfig "github.com/prometheus/common/config"
	"go.yaml.in/yaml/v3"
)

const defaultPathPrefix = "/cgi-bin/solarmonweb"

type Config struct {
	Modules map[string]Module `yaml:"modules"`
}

type Module struct {
	Timeout          time.Duration               `yaml:"timeout"`
	PathPrefix       string                      `yaml:"path_prefix"`
	CollectFault     bool                        `yaml:"collect_fault"`
	HTTPClientConfig promconfig.HTTPClientConfig `yaml:"http_client_config"`
}

func (m *Module) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Module
	*m = Module{
		HTTPClientConfig: promconfig.DefaultHTTPClientConfig,
	}
	return unmarshal((*plain)(m))
}

func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return Load(data)
}

func Default() Config {
	cfg := Config{
		Modules: map[string]Module{
			"default": {
				Timeout:          5 * time.Second,
				PathPrefix:       defaultPathPrefix,
				CollectFault:     true,
				HTTPClientConfig: promconfig.DefaultHTTPClientConfig,
			},
		},
	}
	return cfg
}

func Load(data []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if len(cfg.Modules) == 0 {
		return Config{}, errors.New("config must define at least one module")
	}

	return normalize(cfg)
}

func normalize(cfg Config) (Config, error) {
	for name, module := range cfg.Modules {
		if module.Timeout == 0 {
			module.Timeout = 5 * time.Second
		}
		if module.PathPrefix == "" {
			module.PathPrefix = defaultPathPrefix
		}
		module.PathPrefix = "/" + strings.Trim(strings.TrimSpace(module.PathPrefix), "/")
		if err := module.HTTPClientConfig.Validate(); err != nil {
			return Config{}, err
		}
		cfg.Modules[name] = module
	}

	return cfg, nil
}
