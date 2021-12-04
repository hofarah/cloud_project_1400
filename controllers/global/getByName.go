package global

import (
	"cloudProject/controllers/mainController"
	gameRepo "cloudProject/models/game"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetByName(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "2001", "gerByName")
	defer mainController.FinishAPISpan(ctx)

	games, errStr, err := gameRepo.Repo.GetByName(spanCtx, ctx.Query("name"))
	if err != nil {
		zap.L().Error("get by name err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "01", errStr, 200, err.Error())
	}
	return mainController.Response(ctx, games)
}
