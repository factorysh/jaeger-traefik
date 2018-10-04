package server

import (
	"fmt"
	"time"

	"github.com/uber/jaeger-lib/metrics"
)

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
