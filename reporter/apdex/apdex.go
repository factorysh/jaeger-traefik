package apdex

import (
	"time"

	"github.com/factorysh/jaeger-traefik/conf"
	"github.com/factorysh/jaeger-traefik/reporter"
	_reporter "github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	ApdexSatisfied    = "satisfied"
	ApdexTolerating   = "tolerating"
	ApdexUnsatisfied  = "unsatisfied"
	LabelSatsifaction = "satisfaction"
	LabelDomain       = "domain"
	LabelBackend      = "backend"
)

var apdexCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "apdex",
		Help: "APDEX raw data, clustered by satisfaction",
	},
	[]string{LabelSatsifaction, LabelDomain, LabelBackend})

func init() {
	prometheus.MustRegister(apdexCounter)
	reporter.Reporters["apdex"] = New
}

// ApdexReporter is a jæger reporter, see github.com/jaegertracing/jaeger/cmd/agent/app/reporter
type ApdexReporter struct {
	tags             *conf.TagsConfig
	SatisfiedTarget  time.Duration
	ToleratingTarget time.Duration
}

func New(tags *conf.TagsConfig, config map[string]interface{}) (_reporter.Reporter, error) {
	r := &ApdexReporter{tags: tags}
	s := config["satisfied"]
	if s == "" {
		r.SatisfiedTarget = 250 * time.Millisecond
	} else {
		duration, err := time.ParseDuration(s.(string))
		if err != nil {
			return nil, err
		}
		r.SatisfiedTarget = duration
	}
	t := config["tolerating"]
	if t == "" {
		r.ToleratingTarget = time.Second
	} else {
		duration, err := time.ParseDuration(t.(string))
		if err != nil {
			return nil, err
		}
		r.ToleratingTarget = duration
	}
	return r, nil
}

func (a *ApdexReporter) EmitZipkinBatch(spans []*zipkincore.Span) (err error) {
	return nil
}

func (a *ApdexReporter) EmitBatch(batch *jaegerThrift.Batch) (err error) {
	/*
		process:
		   batch process: traefik
		   	 jaeger.version => Go-2.9.0
		   	 hostname => c626598ba941
		   	 ip => 172.29.0.4
		   	 span.kind => client
		   	 frontend.name => frontend-Host-traefik-2
		   	 backend.name => backend-web-demo
		   	 http.method => GET
		   	 http.url => http://172.29.0.2:80/
		   	 http.host => traefik
		   	 http.status_code => 200

		span:
			sampler.type => const
			sampler.param => true
			span.kind => server
			component => traefik
			http.method => GET
			http.url => /
			http.host => traefik
			span.kind => server
			http.status_code => 200
	*/
	p := batch.GetProcess()
	if p.GetServiceName() != "traefik" {
		log.WithField("ServiceName", p.GetServiceName()).Error("Not træfik")
		return nil
	}

	var satisfaction string

	batches := batch.GetSpans()
	log.WithField("batches length", len(batches)).Debug("spans")
	for _, span := range batches {
		traefik := reporter.TraefikSpan(span)
		log.WithField("traefik", traefik).Debug("spans")
		if traefik.StatusCode < 200 { // 1xx
			return nil
		}
		if traefik.StatusCode >= 300 && traefik.StatusCode < 500 { // 3xx, 4xx
			return nil
		}
		if traefik.StatusCode >= 500 {
			satisfaction = ApdexUnsatisfied
		}
		if satisfaction == "" {
			if traefik.Duration <= a.SatisfiedTarget {
				satisfaction = ApdexSatisfied
			} else if traefik.Duration <= a.ToleratingTarget {
				satisfaction = ApdexTolerating
			} else {
				satisfaction = ApdexUnsatisfied
			}
		}
		apdexCounter.With(prometheus.Labels{
			LabelSatsifaction: satisfaction,
			LabelBackend:      traefik.Backend,
			LabelDomain:       traefik.Host,
		}).Inc()
		log.WithFields(log.Fields{
			LabelSatsifaction: satisfaction,
			LabelBackend:      traefik.Backend,
			LabelDomain:       traefik.Host,
		}).Debug("Apdex inc")
	}

	return nil
}
