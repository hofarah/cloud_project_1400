package dataAnalytic

import (
	"cloudProject/controllers/mainController"
	gameRepo "cloudProject/models/game"
	"cloudProject/pkg/cast"
	"cloudProject/pkg/draw"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func SumSellAnnually(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "50000", "sumSell")
	zap.L().Debug("sumSellAnnually_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)

	startYear, _ := cast.ToInt(ctx.Query("from", "1970"))
	endYear, _ := cast.ToInt(ctx.Query("to", strconv.Itoa(time.Now().Year())))
	data, errStr, err := gameRepo.Repo.GetYearBetweenYears(spanCtx, startYear, endYear)
	if err != nil {
		zap.L().Error("get best on year err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "01", errStr, 500)
	}
	bytes, err := draw.SumSellAnnually(data)
	if err != nil {
		zap.L().Error("sumSellAnnually_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "02", "01", 500, err.Error())
	}
	ctx.Response().Header.SetContentType("image/png")
	mainController.GetAPIStats(ctx.OriginalURL()).AddSuccess()
	return ctx.Send(bytes)
}
