package global

import (
	"cloudProject/controllers/mainController"
	gameRepo "cloudProject/models/game"
	"cloudProject/pkg/cast"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetBestOnGenre(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "2004", "getBestOnGenre")
	defer mainController.FinishAPISpan(ctx)

	genre := ctx.Query("genre")
	if genre == "" {
		return mainController.Error(ctx, "01", "01", 404, "GenreNotFound")
	}
	NBest, _ := cast.ToInt(ctx.Query("N"))
	if NBest == 0 {
		return mainController.Error(ctx, "02", "01", 404, "N_NotFound")
	}
	games, errStr, err := gameRepo.Repo.GetBestOnGenre(spanCtx, genre, NBest)
	if err != nil {
		zap.L().Error("get best on genre err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "03", errStr, 500)
	}
	return mainController.Response(ctx, games)
}
