package createService

import (
	"cloudProject/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
)

func openTracing(ctx *fiber.Ctx) error {
	tracer := opentracing.GlobalTracer()
	if tracer != nil {
		parentSpan := tracer.StartSpan("start service")
		logger.SetSpanTag(parentSpan, "requestID", ctx.Context().ID())
		ctx.Locals("parentSpan", parentSpan)
	}
	return ctx.Next()
}
