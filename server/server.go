package server

import (
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/factorysh/jaeger-lite/reporter/apdex"
	"github.com/jaegertracing/jaeger/cmd/agent/app/processors"
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers/thriftudp"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

type Server interface {
	Serve()
}

func NewServer(listen string, rep reporter.Reporter) (Server, error) {
	//metricsFactory := metrics.NewLocalFactory(0)
	f := &Factory{}

	transport, err := thriftudp.NewTUDPServerTransport(listen)
	if err != nil {
		return nil, err
	}
	maxPacketSize := 65000
	queueSize := 100
	server, err := servers.NewTBufferedServer(transport, queueSize, maxPacketSize, f)
	if err != nil {
		return nil, err
	}
	compactFactory := thrift.NewTCompactProtocolFactory()
	l := zap.NewExample()
	handler := jaegerThrift.NewAgentProcessor(rep)
	return processors.NewThriftProcessor(server, 1, f, compactFactory, handler, l)
}

func New() (Server, error) {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "127.0.0.1:6831"
	}
	return NewServer(listen, apdex.New(250*time.Millisecond, time.Second))
}
