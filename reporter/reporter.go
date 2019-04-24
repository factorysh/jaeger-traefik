package reporter

import (
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
)

var Reporters = make(map[string]ReporterConfig)

type ReporterConfig func(map[string]interface{}) (reporter.Reporter, error)
