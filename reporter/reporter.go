package reporter

import (
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
)

var Reporters = make(map[string]ReporterConfig)

type TagsConfig struct {
	Separator rune
	Labels    []string
}

func (t *TagsConfig) GetSeparator() rune {
	if t.Separator == 0 {
		return ':'
	}
	return t.Separator
}

type ReporterConfig func(*TagsConfig, map[string]interface{}) (reporter.Reporter, error)
