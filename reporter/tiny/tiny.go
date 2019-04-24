package tiny

import (
	"fmt"
	"strings"

	"github.com/factorysh/jaeger-lite/reporter"
	_reporter "github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	labelDomain  = "domain"
	labelBackend = "backend"
	labelProject = "project"
	labelStatus  = "status"
)

var tinyHisogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "tiny",
		Help: "Tiny histogram",
	}, []string{labelDomain, labelBackend, labelProject, labelStatus})

var tinyCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tiny_counter",
		Help: "Tiny counter",
	}, []string{labelDomain, labelBackend, labelProject, labelStatus})

func init() {
	prometheus.MustRegister(tinyHisogram)
	prometheus.MustRegister(tinyCounter)
	reporter.Reporters["tiny"] = New
}

type Tiny struct {
}

func New(config map[string]interface{}) (_reporter.Reporter, error) {
	return &Tiny{}, nil
}

func (t *Tiny) EmitZipkinBatch(spans []*zipkincore.Span) (err error) {
	return nil
}

func (t *Tiny) EmitBatch(batch *jaegerThrift.Batch) (err error) {
	p := batch.GetProcess()
	if p.GetServiceName() != "traefik" {
		log.WithField("ServiceName", p.GetServiceName()).Error("Not træfik")
		return nil
	}
	for _, span := range batch.GetSpans() {
		traefik := reporter.TraefikSpan(span)
		log.WithField("traefik", traefik).Debug("spans")
		b := strings.Split(traefik.Backend, "-")
		tinyCounter.With(prometheus.Labels{
			labelProject: b[1],
			labelBackend: traefik.Backend,
			labelDomain:  traefik.Host,
			labelStatus:  fmt.Sprintf("%vxx", traefik.StatusCode%100),
		}).Inc()
		if traefik.StatusCode >= 200 && traefik.StatusCode < 300 {
			tinyHisogram.With(prometheus.Labels{
				labelProject: b[1],
				labelBackend: traefik.Backend,
				labelDomain:  traefik.Host,
			}).Observe(float64(traefik.Duration))
		}
	}
	return nil
}
