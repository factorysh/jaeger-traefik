package apdex

import (
	"time"

	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	Satisfied   = "satisified"
	Tolerating  = "toleration"
	Unsatisfied = "unsatisfied"
)

var apdexCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "apdex",
		Help: "APDEX raw data, clustered by satisfaction",
	},
	[]string{"satisfaction", "domain", "backend"})

func init() {
	prometheus.MustRegister(apdexCounter)
}

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
		return nil
	}

	var backend, host string
	for _, tag := range p.GetTags() {
		switch tag.GetKey() {
		case "backend.name":
			backend = tag.GetVStr()
		case "http.host":
			host = tag.GetVStr()
		}
	}
	for _, span := range batch.GetSpans() {
		var satisfaction string
		for _, tag := range span.GetTags() {
			switch tag.GetKey() {
			case "span.kind":
				if tag.GetVStr() != "server" {
					return nil
				}
			case "component":
				if tag.GetVStr() != "traefik" {
					return nil
				}
			case "http.status_code":
				status := tag.GetVLong()
				if status < 200 { // 1xx
					return nil
				}
				if status >= 300 && status < 500 { // 3xx, 4xx
					return nil
				}
				if status >= 500 {
					satisfaction = Unsatisfied
				}
			}
		}
		if satisfaction == "" {
			d := time.Duration(span.GetDuration()) * time.Millisecond
			if d <= a.SatisfiedTarget {
				satisfaction = Satisfied
			} else if d <= a.ToleratingTarget {
				satisfaction = Tolerating
			} else {
				satisfaction = Unsatisfied
			}
		}
		apdexCounter.With(prometheus.Labels{
			"satisfaction": satisfaction,
			"backend":      backend,
			"host":         host,
		}).Inc()
		log.Debug("paf", satisfaction)
	}

	return nil
}

func New(toleratingTarget, satisfiedTarget time.Duration) *ApdexReporter {
	log.Info("New Apdex reporter")
	return &ApdexReporter{
		ToleratingTarget: toleratingTarget,
		SatisfiedTarget:  satisfiedTarget,
	}
}
