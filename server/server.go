package server

import (
	"go.uber.org/zap"

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

func New() (Server, error) {
	//metricsFactory := metrics.NewLocalFactory(0)
	f := &Factory{}

	transport, err := thriftudp.NewTUDPServerTransport("127.0.0.1:6831")
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
	var rep reporter.Reporter
	rep = &Reporter{}
	handler := jaegerThrift.NewAgentProcessor(rep)
	p, err := processors.NewThriftProcessor(server, 1, f, compactFactory, handler, l)

	return p, err
}
