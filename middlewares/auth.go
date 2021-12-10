package middlewares

import (
	"cloudProject/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
)

var authServiceUrl = os.Getenv("AUTH_SERVICE_URL")

func Auth() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		statusCode, err := utils.HttpRequest("GET", authServiceUrl, nil, nil, map[string]string{
			"token":      ctx.Get("token"),
			"secret":     ctx.Get("secret"),
			"serviceKey": os.Getenv("SERVICE_KEY"),
		})
		if statusCode != 200 || err != nil {
			zap.L().Error("http request err", zap.Error(err))
			return ctx.SendStatus(403)
		}
		return ctx.Next()
	}
}
