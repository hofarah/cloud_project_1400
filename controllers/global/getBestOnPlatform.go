package global

import (
	"cloudProject/controllers/mainController"
	gameRepo "cloudProject/models/game"
	"cloudProject/pkg/cast"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetBestOnPlatform(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "2002", "getBestOnPlatform")
	defer mainController.FinishAPISpan(ctx)

	platform := ctx.Query("platform")
	if platform == "" {
		return mainController.Error(ctx, "01", "01", 404, "platformNotFound")
	}
	NBest, _ := cast.ToInt(ctx.Query("N"))
	if NBest == 0 {
		return mainController.Error(ctx, "02", "01", 404, "N_NotFound")
	}
	games, errStr, err := gameRepo.Repo.GetBestOnPlatform(spanCtx, platform, NBest)
	if err != nil {
		zap.L().Error("get best on platform err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "03", errStr, 500)
	}
	return mainController.Response(ctx, games)
}
