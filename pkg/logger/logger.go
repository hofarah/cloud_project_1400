package logger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	//jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var tracers map[string]opentracing.Tracer

func init() {
	tracers = map[string]opentracing.Tracer{}
	debug := os.Getenv("LOGLEVEL")
	var level zapcore.Level
	switch debug {
	case "debug":
		level = zapcore.DebugLevel
	case "error":
		level = zapcore.ErrorLevel
	case "warn":
		level = zapcore.WarnLevel
	default:
		level = zapcore.InfoLevel
	}
	path := "/logs/app.log"
	if os.Getenv("DEV_MODE") == "1" {
		path = "stdout"
	}
	logger, _ := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{path},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()
	zap.ReplaceGlobals(logger)
}

// InitJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func InitJaeger(service string) opentracing.Tracer {
	if trace, ok := tracers[service]; ok {
		return trace
	}
	agentIP := os.Getenv("JAEGER_AGENT_HOST")
	if agentIP == "" {
		return nil
	}
	agentPort := os.Getenv("JAEGER_AGENT_PORT")
	zap.L().Info("sending trace to", zap.String("agentIP", agentIP), zap.String("agentPort", agentPort))
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: agentIP + ":" + agentPort,
		},
		ServiceName: service,
	}
	metricsFactory := metrics.NullFactory
	tracer, _, err := cfg.NewTracer(config.Logger(jaeger.StdLogger), config.Metrics(metricsFactory))
	if err != nil {
		panic(err)
	}
	tracers[service] = tracer
	return tracer
}
func JaegerErrorLog(span opentracing.Span, err error) {
	SetSpanTag(span, "error", true)
	SetSpanFields(span, log.String("event", "error"), log.String("value", err.Error()))
}
func SetSpanTag(span opentracing.Span, key string, value interface{}) {
	if span != nil && len(key) != 0 && value != nil {
		span.SetTag(key, value)
	}
}
func SetSpanFields(span opentracing.Span, fields ...log.Field) {
	if span != nil && fields != nil {
		span.LogFields(fields...)
	}
}
func StartSpan(ctx context.Context, tracer opentracing.Tracer, operationName string) (span opentracing.Span, traceID string) {
	if tracer == nil {
		return nil, ""
	}
	if s := GetSpanFromContext(ctx); span != nil {
		span = tracer.StartSpan(operationName, opentracing.ChildOf(s.Context()))
		traceID = GetTraceIDFromSpan(span)
	}
	return
}
func FinishSpan(span opentracing.Span) {
	if span == nil {
		return
	}
	span.Finish()
}
func GetSpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}
func GetTraceIDFromSpan(span opentracing.Span) string {
	if span == nil {
		return ""
	}
	spanCtx, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return ""
	}
	return spanCtx.TraceID().String()
}
func GetTraceIDFromContext(context context.Context) string {
	return GetTraceIDFromSpan(GetSpanFromContext(context))
}
