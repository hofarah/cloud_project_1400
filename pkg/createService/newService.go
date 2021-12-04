package createService

import (
	"cloudProject/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func New(serviceName string) *fiber.App {
	zap.L().Info("start" + serviceName + "service")
	processConfig()
	tracer := logger.InitJaeger(serviceName)
	opentracing.SetGlobalTracer(tracer)
	app := fiber.New(fiber.Config{ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		ctx.Status(fiber.StatusInternalServerError)
		_, _ = ctx.Write(nil)
		return nil
	}})
	app.Use(Recover())
	app.Use(openTracing)
	gracefulShutdown(app)
	return app
}
