package authentication

import (
	"cloudProject/controllers/mainController"
	"cloudProject/models/user"
	"cloudProject/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Check(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "1001", "check")
	defer mainController.FinishAPISpan(ctx)

	payload, err := jwt.Verify(ctx.Get("token"))
	if err != nil {
		zap.L().Error("verify token err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "01", "01", 403)
	}
	u, exist, errStr, err := user.Repo.GetUserByUserName(spanCtx, payload.Username)
	if err != nil || !exist {
		zap.L().Error("un authorized user", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "02", errStr, 403)
	}
	if u.Secret != ctx.Get("secret") {
		zap.L().Error("invalid secret", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "03", "01", 403)
	}
	return mainController.Response(ctx, nil)
}
