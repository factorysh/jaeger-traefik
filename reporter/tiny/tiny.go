package tiny

import (
	"fmt"
	"strings"

	"github.com/factorysh/jaeger-traefik/reporter"
	_reporter "github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	labelDomain     = "domain"
	labelBackend    = "backend"
	labelProject    = "project"
	labelStatus     = "status"
	labelPhatStatus = "status_xx"
)

var tinyHisogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "tiny",
		Help: "Tiny histogram",
	}, []string{labelDomain, labelBackend, labelProject, labelStatus, labelPhatStatus})

var tinyCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tiny_counter",
		Help: "Tiny counter",
	}, []string{labelDomain, labelBackend, labelProject, labelStatus, labelPhatStatus})

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
		if len(b) == 0 || traefik.StatusCode == 0 {
			continue
		}
		project := b[len(b)-1]
		// b is something like [backend web front demo]
		// service := strings.Join(b[1:len(b)], "-")
		phat := fmt.Sprintf("%vxx", traefik.StatusCode%100)
		tinyCounter.With(prometheus.Labels{
			labelProject:    project,
			labelBackend:    traefik.Backend,
			labelDomain:     traefik.Host,
			labelPhatStatus: phat,
			labelStatus:     string(traefik.StatusCode),
		}).Inc()
		if traefik.StatusCode >= 200 && traefik.StatusCode < 300 {
			tinyHisogram.With(prometheus.Labels{
				labelProject:    project,
				labelBackend:    traefik.Backend,
				labelDomain:     traefik.Host,
				labelPhatStatus: phat,
				labelStatus:     string(traefik.StatusCode),
			}).Observe(float64(traefik.Duration))
		}
		log.WithFields(log.Fields{
			labelProject:    project,
			labelBackend:    traefik.Backend,
			labelDomain:     traefik.Host,
			labelPhatStatus: phat,
			labelStatus:     traefik.StatusCode,
		}).Info("Prometheus")
	}
	return nil
}
