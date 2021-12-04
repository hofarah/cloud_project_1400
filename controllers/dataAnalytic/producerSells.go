package dataAnalytic

import (
	"cloudProject/controllers/mainController"
	"cloudProject/pkg/draw"
	"github.com/gofiber/fiber/v2"
	"github.com/wcharczuk/go-chart/v2"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func ProducerSells(ctx *fiber.Ctx) error {
	_, traceID := mainController.InitAPI(ctx, "50000", "genreSells")
	zap.L().Debug("ProducerSells_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)
	startYear := ctx.Query("from", "1970")
	endYear := ctx.Query("to", strconv.Itoa(time.Now().Year()))
	producer1 := ctx.Query("p1")
	producer2 := ctx.Query("p2")
	if producer1 == "" || producer2 == "" {
		zap.L().Error("ProducerSells_params_invalid_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", "producer not found", 404)
	}
	//todo fetch data from main service
	data := map[string][]chart.Value{}
	bytes, err := draw.ProducerSells(startYear, endYear, data)
	if err != nil {
		zap.L().Error("genreSells_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", err.Error(), 500)
	}
	ctx.Response().Header.SetContentType("image/png")
	return ctx.Send(bytes)
}
