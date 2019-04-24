package conf

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen    string
	Reporters map[string]map[string]interface{}
}

func New(raw []byte) (*Config, error) {
	cfg := &Config{
		Reporters: make(map[string]map[string]interface{}),
	}
	err := yaml.Unmarshal(raw, cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Listen == "" {
		cfg.Listen = "127.0.0.1:5000"
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
