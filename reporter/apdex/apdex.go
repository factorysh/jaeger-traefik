package apdex

import (
	"time"

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
}

// ApdexReporter is a jæger reporter, see github.com/jaegertracing/jaeger/cmd/agent/app/reporter
type ApdexReporter struct {
	SatisfiedTarget  time.Duration
	ToleratingTarget time.Duration
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

	var satisfaction, backend, host string

	batches := batch.GetSpans()
	log.WithField("length", len(batches)).Debug("spans")
	for s, span := range batches {
		tags := span.GetTags()
		for _, tag := range tags {
			log.WithField("key", tag.GetKey()).WithField("n", s).Debug("span key")
			switch tag.GetKey() {
			case "span.kind":
				log.WithField("span.kind", tag.GetVStr()).Debug("span kind")
				if tag.GetVStr() != "server" {
					log.WithField("span.kind", tag.GetVStr()).Error("Not a server")
					continue
				}
			case "backend.name":
				log.WithField(tag.GetKey(), tag.GetVStr()).Debug("process tag")
				backend = tag.GetVStr()
			case "http.host":
				log.WithField(tag.GetKey(), tag.GetVStr()).Debug("process tag")
				host = tag.GetVStr()
			case "component":
				if tag.GetVStr() != "traefik" {
					log.WithField("component", tag.GetVStr()).Error("Not traefik")
					return nil
				}
			case "http.status_code":
				status := tag.GetVLong()
				log.WithField("http.status_code", tag.GetVLong()).Debug("status")
				if status < 200 { // 1xx
					return nil
				}
				if status >= 300 && status < 500 { // 3xx, 4xx
					return nil
				}
				if status >= 500 {
					satisfaction = ApdexUnsatisfied
				}
			}
		}
		if satisfaction == "" {
			log.WithField("duration", span.GetDuration()).Debug("Duration")
			d := time.Duration(span.GetDuration()) * time.Microsecond
			if d <= a.SatisfiedTarget {
				satisfaction = ApdexSatisfied
			} else if d <= a.ToleratingTarget {
				satisfaction = ApdexTolerating
			} else {
				satisfaction = ApdexUnsatisfied
			}
		}
	}
	apdexCounter.With(prometheus.Labels{
		LabelSatsifaction: satisfaction,
		LabelBackend:      backend,
		LabelDomain:       host,
	}).Inc()
	log.WithFields(log.Fields{
		LabelSatsifaction: satisfaction,
		LabelBackend:      backend,
		LabelDomain:       host,
	}).Debug("Apdex inc")

	return nil
}

func New(toleratingTarget, satisfiedTarget time.Duration) *ApdexReporter {
	log.Info("New Apdex reporter")
	return &ApdexReporter{
		ToleratingTarget: toleratingTarget,
		SatisfiedTarget:  satisfiedTarget,
	}
}
