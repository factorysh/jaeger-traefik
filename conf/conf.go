package conf

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAdmin  string                            `yaml:"listen_admin"`
	ListenJaeger string                            `yaml:"listen_jaeger"`
	Reporters    map[string]map[string]interface{} `yaml:"reporters"`
}

func New(raw []byte) (*Config, error) {
	cfg := &Config{
		Reporters: make(map[string]map[string]interface{}),
	}
	err := yaml.Unmarshal(raw, cfg)
	if err != nil {
		return nil, err
	}
	if cfg.ListenAdmin == "" {
		cfg.ListenAdmin = "127.0.0.1:8080"
	}
	if cfg.ListenJaeger == "" {
		cfg.ListenJaeger = "127.0.0.1:6831"
	}
	return cfg, nil
}

func Read(cfgPath string) (*Config, error) {
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return New(raw)
}
