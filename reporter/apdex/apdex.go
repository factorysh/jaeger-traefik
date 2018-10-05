package apdex

import (
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"github.com/prometheus/client_golang/prometheus"
)

type ApdexReporter struct {
	gauge prometheus.Gauge
}

func (a *ApdexReporter) EmitZipkinBatch(spans []*zipkincore.Span) (err error) {
	return nil
}

func (a *ApdexReporter) EmitBatch(batch *jaegerThrift.Batch) (err error) {
	return nil
}

func New() *ApdexReporter {
	return &ApdexReporter{
		gauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "apdex",
			Help: "Apdex compute from http status and performance, per domain",
		}),
	}
}
