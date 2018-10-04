package server

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/cmd/agent/app/processors"
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers/thriftudp"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"github.com/uber/jaeger-lib/metrics"
)

type Server interface {
	Serve()
}

type Factory struct {
}

func (f *Factory) Counter(name string, tags map[string]string) metrics.Counter {
	return &Counter{name, tags}
}

func (f *Factory) Timer(name string, tags map[string]string) metrics.Timer {
	return &Timer{name, tags}
}

func (f *Factory) Gauge(name string, tags map[string]string) metrics.Gauge {
	return &Gauge{name, tags}
}

// Namespace returns a nested metrics factory.
func (f *Factory) Namespace(name string, tags map[string]string) metrics.Factory {
	fmt.Println("Namespace", name, tags)
	return f
}

type Counter struct {
	Name string
	Tags map[string]string
}

func (c *Counter) Inc(i int64) {
	fmt.Println(c.Name, c.Tags, "Inc", i)
}

type Timer struct {
	Name string
	Tags map[string]string
}

func (t *Timer) Record(d time.Duration) {
	fmt.Println(t.Name, t.Tags, "Record", d)
}

type Gauge struct {
	Name string
	Tags map[string]string
}

func (g *Gauge) Update(i int64) {
	fmt.Println(g.Name, g.Tags, "Update", i)
}

type Processor struct {
}

func (p *Processor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	fmt.Println("iprot:", iprot)
	fmt.Println("oprot:", oprot)
	return true, nil
}

type Agent struct {
}

type Reporter struct {
}

func (r *Reporter) EmitZipkinBatch(spans []*zipkincore.Span) (err error) {
	return nil
}

func printTags(tags []*jaegerThrift.Tag) {
	for _, tag := range tags {
		switch tag.GetVType() {
		case jaegerThrift.TagType_STRING:
			fmt.Println("\t", tag.GetKey(), "=>", tag.GetVStr())
		case jaegerThrift.TagType_BOOL:
			fmt.Println("\t", tag.GetKey(), "=>", tag.GetVBool())
		case jaegerThrift.TagType_LONG:
			fmt.Println("\t", tag.GetKey(), "=>", tag.GetVLong())
		default:
			fmt.Println("\t", tag.GetKey(), tag.GetVType(), tag)
		}
	}
}

func (r *Reporter) EmitBatch(batch *jaegerThrift.Batch) (err error) {
	p := batch.GetProcess()
	fmt.Println("batch process:", p.GetServiceName())
	printTags(p.GetTags())

	for _, span := range batch.GetSpans() {
		printTags(span.GetTags())
		fmt.Println("span:", span)
	}
	return nil
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

func eventsLoop(s servers.Server) {
	var r *servers.ReadBuf
	for {
		r = <-s.DataChan()
		data := r.GetBytes()
		s.DataRecd(r)
		fmt.Println("data:", data)
	}
}
