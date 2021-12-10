package dataAnalytic

import (
	"cloudProject/controllers/mainController"
	"cloudProject/models/game"
	"cloudProject/pkg/cast"
	"cloudProject/pkg/draw"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func ProducerSells(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "50000", "producerSells")
	zap.L().Debug("ProducerSells_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)
	startYear, _ := cast.ToInt(ctx.Query("from", "1970"))
	endYear, _ := cast.ToInt(ctx.Query("to", strconv.Itoa(time.Now().Year())))
	producer1 := ctx.Query("p1")
	producer2 := ctx.Query("p2")
	if producer1 == "" || producer2 == "" {
		zap.L().Error("ProducerSells_params_invalid_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", "01", 404, "producerNotFound")
	}
	data, errStr, err := game.Repo.GetProducerOnYears(spanCtx, producer1, producer2, startYear, endYear)
	if err != nil {
		zap.L().Error("get producer on years failed", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "02", errStr, 500)
	}
	bytes, err := draw.ProducerSells(data)
	if err != nil {
		zap.L().Error("genreSells_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "03", "01", 500, err.Error())
	}
	ctx.Response().Header.SetContentType("image/png")
	mainController.GetAPIStats(ctx.OriginalURL()).AddSuccess()
	return ctx.Send(bytes)
}
