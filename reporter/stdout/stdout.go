package stdout

import (
	"fmt"

	"github.com/factorysh/jaeger-traefik/conf"
	"github.com/factorysh/jaeger-traefik/reporter"
	_reporter "github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
)

type Reporter struct {
	tags *conf.TagsConfig
}

func init() {
	reporter.Reporters["stdout"] = New
}

func New(tags *conf.TagsConfig, config map[string]interface{}) (_reporter.Reporter, error) {
	return &Reporter{}, nil
}

func (r *Reporter) EmitZipkinBatch(spans []*zipkincore.Span) (err error) {
	return nil
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
