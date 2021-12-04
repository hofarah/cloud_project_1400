package dataAnalyticController

import (
	"cloudProject/controllers/mainController"
	"cloudProject/pkg/draw"
	"github.com/gofiber/fiber/v2"
	"github.com/wcharczuk/go-chart/v2"
	"go.uber.org/zap"
)

func SumSellAnnually(ctx *fiber.Ctx) error {
	_, traceID := mainController.InitAPI(ctx, "50000", "genreSells")
	zap.L().Debug("sumSellAnnually_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)

	from := ctx.Query("from")
	to := ctx.Query("to")
	if from == "" || to == "" {
		zap.L().Error("sumSellAnnually_params_invalid_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", "producer not found", 404)
	}
	//todo fetch data from main service
	data := []chart.Value{}
	bytes, err := draw.SumSellAnnually(data)
	if err != nil {
		zap.L().Error("sumSellAnnually_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", err.Error(), 500)
	}
	ctx.Response().Header.SetContentType("image/png")
	return ctx.Send(bytes)
}
