package reporter

import (
	"time"

	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

type Traefik struct {
	Backend    string        // backend.name
	Host       string        // http.host
	StatusCode int64         // http.status_code
	Duration   time.Duration // duration
}

func TraefikSpan(span *jaeger.Span) *Traefik {
	t := &Traefik{}
	for _, tag := range span.GetTags() {
		switch tag.GetKey() {
		case "span.kind":
			if tag.GetVStr() != "server" {
				continue
			}
		case "backend.name":
			t.Backend = tag.GetVStr()
		case "http.host":
			t.Host = tag.GetVStr()
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
