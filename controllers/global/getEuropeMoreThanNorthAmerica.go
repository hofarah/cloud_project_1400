package global

import (
	"cloudProject/controllers/mainController"
	gameRepo "cloudProject/models/game"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetEuropeMoreThanNorthAmerica(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "2006", "getEuropeMoreThanNorthAmerica")
	defer mainController.FinishAPISpan(ctx)

	games, errStr, err := gameRepo.Repo.GetEuropeMoreThanNorthAmerica(spanCtx)
	if err != nil {
		zap.L().Error("get best on year err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "03", errStr, 500)
	}
	return mainController.Response(ctx, games)
}
