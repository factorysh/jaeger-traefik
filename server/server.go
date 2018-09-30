package server

import (
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers"
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers/thriftudp"
	"github.com/uber/jaeger-lib/metrics"
)

func New() (servers.Server, error) {

	metricsFactory := metrics.NewLocalFactory(0)

	transport, err := thriftudp.NewTUDPServerTransport("127.0.0.1:6831")
	if err != nil {
		return nil, err
	}
	maxPacketSize := 65000
	queueSize := 100
	server, err := servers.NewTBufferedServer(transport, queueSize, maxPacketSize, metricsFactory)

	return server, err
}
