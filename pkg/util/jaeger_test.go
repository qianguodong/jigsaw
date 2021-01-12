package util_test

import (
	"fmt"
	"testing"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

func TestJaeger(t *testing.T) {
	metrics := prometheus.New()

	// Initialize the tracing configuration.
	conf, err := config.FromEnv()
	if err != nil {
		t.Error(err)
	}
	conf.ServiceName = "testing"
	conf.Sampler = &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
	conf.Reporter = &config.ReporterConfig{
		LogSpans:           true,
		LocalAgentHostPort: fmt.Sprintf("%s:%d", "127.0.0.1", 6831),
	}
	// Use the configuration to create a new tracer.
	zipkinB3HTTPHeaderPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, closer, err := conf.NewTracer(
		config.Metrics(metrics),
		config.Logger(jaeger.StdLogger),
		config.ZipkinSharedRPCSpan(true),
		config.Injector(opentracing.HTTPHeaders, zipkinB3HTTPHeaderPropagator),
		config.Extractor(opentracing.HTTPHeaders, zipkinB3HTTPHeaderPropagator),
	)
	if err != nil {
		t.Error(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// build zipkin trace
	request := opentracing.HTTPHeadersCarrier{}
	request.Set("x-b3-traceid", "1")
	request.Set("x-b3-spanid", "2")
	request.Set("x-b3-parentspanid", "1")
	request.Set("x-b3-sampled", "1")
	parentSpanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, request)
	if err != nil {
		t.Error(err)
	}
	serverSpan1 := opentracing.GlobalTracer().StartSpan(
		"sp1",
		ext.RPCServerOption(parentSpanCtx),
	)
	defer serverSpan1.Finish()
}
