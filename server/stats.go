package server

import (
	"fmt"
	"time"

	"github.com/uber/jaeger-lib/metrics"
)

type Factory struct {
}

func (f *Factory) Counter(options metrics.Options) metrics.Counter {
	return &Counter{options.Name, options.Tags}
}

func (f *Factory) Timer(options metrics.TimerOptions) metrics.Timer {
	return &Timer{options.Name, options.Tags}
}

func (f *Factory) Histogram(options metrics.HistogramOptions) metrics.Histogram {
	return nil
}

func (f *Factory) Gauge(options metrics.Options) metrics.Gauge {
	return &Gauge{options.Name, options.Tags}
}

// Namespace returns a nested metrics factory.
func (f *Factory) Namespace(options metrics.NSOptions) metrics.Factory {
	fmt.Println("Namespace", options)
	return nil
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
