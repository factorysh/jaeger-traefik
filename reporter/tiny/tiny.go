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

func init() {
	reporter.Reporters["tiny"] = New
}

type Tiny struct {
	tags *reporter.TagsConfig
}

var (
	once         bool
	tinyHisogram *prometheus.HistogramVec
	tinyCounter  *prometheus.CounterVec
)

func New(tags *reporter.TagsConfig, config map[string]interface{}) (_reporter.Reporter, error) {
	if once {
		panic("You can register prometheus stuff just one time.")
	}
	once = true
	labels := []string{labelDomain, labelBackend, labelStatus, labelPhatStatus}
	for _, l := range tags.Labels {
		labels = append(labels, l)
	}
	tinyHisogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "tiny",
			Help: "Tiny histogram",
		}, labels)

	tinyCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tiny_counter",
			Help: "Tiny counter",
		}, labels)
	prometheus.MustRegister(tinyHisogram)
	prometheus.MustRegister(tinyCounter)
	return &Tiny{tags}, nil
}

func (t *Tiny) EmitZipkinBatch(spans []*zipkincore.Span) (err error) {
	return nil
}

func split(txt string, sep rune) []string {
	slugs := make([]string, 0)
	i := 0
	for {
		poz := strings.IndexRune(txt[i:], sep)
		fmt.Println(i, poz)
		var slug string
		if poz == -1 {
			slug = txt[i:len(txt)]
		} else {
			slug = txt[i : i+poz]
		}
		slugs = append(slugs, slug)
		if poz == -1 {
			break
		}
		i += poz + 1
	}
	return slugs
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
		b := split(traefik.Backend, t.tags.GetSeparator())
		if len(b) == 0 || traefik.StatusCode == 0 {
			continue
		}
		phat := fmt.Sprintf("%vxx", traefik.StatusCode/100)
		status := fmt.Sprintf("%d", traefik.StatusCode)
		labels := prometheus.Labels{
			labelBackend:    traefik.Backend,
			labelDomain:     traefik.Host,
			labelPhatStatus: phat,
			labelStatus:     status,
		}
		for i, label := range t.tags.Labels {
			if i < len(b) {
				labels[label] = b[i]
			} else {
				labels[label] = ""
			}
		}
		tinyCounter.With(labels).Inc()
		if traefik.StatusCode >= 200 && traefik.StatusCode < 300 {
			tinyHisogram.With(labels).Observe(float64(traefik.Duration))
		}
		fields := make(log.Fields)
		for k, v := range labels {
			fields[k] = v
		}
		log.WithFields(fields).Info("Prometheus")
	}
	return nil
}
