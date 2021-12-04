package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func CheckRequiredHeaders(ctx *fiber.Ctx) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		apiKey := ctx.Get("secret")
		if apiKey == "" || len(apiKey) != 32 {
			zap.L().Error("secret not found", zap.Any("url", ctx.OriginalURL()))
			return ctx.SendStatus(403)
		}
		if ctx.Get("token") == "" {
			zap.L().Error("token not found", zap.Any("url", ctx.OriginalURL()))
			return ctx.SendStatus(403)
		}
		if ctx.Get("serviceKey") == "" && ctx.OriginalURL() != "/authentication/signup" {
			zap.L().Error("serviceKey not found", zap.Any("url", ctx.OriginalURL()))
			return ctx.SendStatus(403)
		}
		return ctx.Next()
	}
}
