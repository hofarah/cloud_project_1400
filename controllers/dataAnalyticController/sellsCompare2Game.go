package dataAnalyticController

import (
	"cloudProject/controllers/mainController"
	"cloudProject/pkg/draw"
	"github.com/gofiber/fiber/v2"
	"github.com/wcharczuk/go-chart/v2"
	"go.uber.org/zap"
)

func SellsCompare2Game(ctx *fiber.Ctx) error {
	_, traceID := mainController.InitAPI(ctx, "50000", "genreSells")
	zap.L().Debug("sellsCompare2Game_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)

	game1 := ctx.Query("game1")
	game2 := ctx.Query("game2")
	if game1 == "" || game2 == "" {
		zap.L().Error("sellsCompare2Game_params_invalid_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", "producer not found", 404)
	}
	//todo fetch data from main service
	data := map[string][]chart.Value{}
	bytes, err := draw.SellsCompare2Game(data)
	if err != nil {
		zap.L().Error("sellsCompare2Game_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", err.Error(), 500)
	}
	ctx.Response().Header.SetContentType("image/png")
	return ctx.Send(bytes)
}
