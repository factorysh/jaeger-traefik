package reporter

import (
	"github.com/factorysh/jaeger-traefik/conf"
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
)

var Reporters = make(map[string]ReporterConfig)

type ReporterConfig func(*conf.TagsConfig, map[string]interface{}) (reporter.Reporter, error)
