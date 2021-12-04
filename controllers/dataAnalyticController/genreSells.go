package dataAnalyticController

import (
	"cloudProject/controllers/mainController"
	"cloudProject/pkg/draw"
	"github.com/gofiber/fiber/v2"
	"github.com/wcharczuk/go-chart/v2"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GenreSells(ctx *fiber.Ctx) error {
	_, traceID := mainController.InitAPI(ctx, "50000", "genreSells")
	zap.L().Debug("genreSells_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)
	startYear := ctx.Query("from", "1970")
	endYear := ctx.Query("to", strconv.Itoa(time.Now().Year()))
	//todo fetch data from main service
	data := []chart.Value{}
	bytes, err := draw.GenreSells(startYear, endYear, data)
	if err != nil {
		zap.L().Error("genreSells_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", err.Error(), 500)
	}
	ctx.Response().Header.SetContentType("image/png")
	return ctx.Send(bytes)
}
