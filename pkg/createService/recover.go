package createService

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"runtime"
)

func Recover() fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 8096)
				stack = stack[:runtime.Stack(stack, false)]
				zap.L().Debug("panic_err", zap.String("stack", string(stack)))
				ctx.Status(500)
			}
		}()
		return ctx.Next()
	}
}
