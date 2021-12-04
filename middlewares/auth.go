package middlewares

import (
	"cloudProject/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
)

var authServiceUrl string

func init() {
	if authServiceUrl = os.Getenv("AUTH_SERVICE_URL"); authServiceUrl == "" {
		panic("AUTH_SERVICE_URL not found")
	}
}

func Auth(ctx *fiber.Ctx) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		statusCode, err := utils.HttpRequest("GET", authServiceUrl, nil, nil, map[string]string{
			"token":      ctx.Get("token"),
			"secret":     ctx.Get("secret"),
			"serviceKey": os.Getenv("serviceKey"),
		})
		if statusCode != 200 {
			zap.L().Error("http request err", zap.Error(err))
			return ctx.SendStatus(403)
		}
		return ctx.Next()
	}
}
