package dataAnalytic

import (
	"cloudProject/controllers/mainController"
	"cloudProject/models/game"
	"cloudProject/pkg/draw"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func SellsCompare2Game(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "50000", "genreSells")
	zap.L().Debug("sellsCompare2Game_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)

	game1 := ctx.Query("g1")
	game2 := ctx.Query("g2")
	if game1 == "" || game2 == "" {
		zap.L().Error("sellsCompare2Game_params_invalid_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", "game not found", 404)
	}
	data, errStr, err := game.Repo.GetGamesSell(spanCtx, game1, game2)
	if err != nil {
		zap.L().Error("get game sales failed", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "02", errStr, 500)
	}
	bytes, err := draw.SellsCompare2Game(data)
	if err != nil {
		zap.L().Error("sellsCompare2Game_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", err.Error(), 500)
	}
	ctx.Response().Header.SetContentType("image/png")
	return ctx.Send(bytes)
}
