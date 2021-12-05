package global

import (
	"cloudProject/controllers/mainController"
	gameRepo "cloudProject/models/game"
	"cloudProject/pkg/cast"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetBest5SellsByYearAndPlatform(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "2005", "getBest5SellsByYearAndPlatform")
	defer mainController.FinishAPISpan(ctx)

	year, _ := cast.ToInt(ctx.Query("year"))
	if year == 0 {
		return mainController.Error(ctx, "01", "01", 404)
	}
	platform := ctx.Query("platform")
	if platform == "" {
		return mainController.Error(ctx, "02", "01", 404)
	}
	games, errStr, err := gameRepo.Repo.GetBestOnYearAndPlatform(spanCtx, platform, year, 5)
	if err != nil {
		zap.L().Error("get best on year err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "03", errStr, 500)
	}
	return mainController.Response(ctx, games)
}
