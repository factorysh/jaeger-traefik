package reporter

import (
	"time"

	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	log "github.com/sirupsen/logrus"
)

type Traefik struct {
	Backend    string        // backend.name
	Frontend   string        // frontend.name
	URL        string        // http.url
	Host       string        // http.host
	Method     string        // http.method
	StatusCode int64         // http.status_code
	Duration   time.Duration // duration
}

func TraefikSpan(span *jaeger.Span) *Traefik {
	t := &Traefik{}
	for _, tag := range span.GetTags() {
		log.WithField(tag.GetKey(), tag).Debug("TraefikSpan")
		switch tag.GetKey() {
		case "span.kind":
			if tag.GetVStr() != "server" {
				continue
			}
		case "backend.name":
			t.Backend = tag.GetVStr()
		case "frontend.name":
			t.Frontend = tag.GetVStr()
		case "http.host":
			t.Host = tag.GetVStr()
		case "http.url":
			t.URL = tag.GetVStr()
		case "http.method":
			t.Method = tag.GetVStr()
		case "component":
			if tag.GetVStr() != "traefik" {
				return nil
			}
		case "http.status_code":
			t.StatusCode = tag.GetVLong()
		}
	}
	t.Duration = time.Duration(span.GetDuration()) * time.Microsecond
	return t
}
