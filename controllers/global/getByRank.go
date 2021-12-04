package global

import (
	"cloudProject/controllers/mainController"
	gameRepo "cloudProject/models/game"
	"cloudProject/pkg/cast"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetByRank(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "2000", "gerByRank")
	defer mainController.FinishAPISpan(ctx)
	rank, _ := cast.ToInt(ctx.Query("rank"))
	game, _, errStr, err := gameRepo.Repo.GetByRank(spanCtx, rank)
	if err != nil {
		zap.L().Error("get by rank err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "01", errStr, 200, err.Error())
	}
	return mainController.Response(ctx, game)
}
