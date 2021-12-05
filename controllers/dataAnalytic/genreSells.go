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

func GenreSells(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "50000", "genreSells")
	zap.L().Debug("genreSells_start", zap.String("traceID", traceID))
	defer mainController.FinishAPISpan(ctx)
	startYear, _ := cast.ToInt(ctx.Query("from", "1970"))
	endYear, _ := cast.ToInt(ctx.Query("to", strconv.Itoa(time.Now().Year())))
	data, errStr, err := gameRepo.Repo.GetGenreBetweenYears(spanCtx, startYear, endYear)
	if err != nil {
		zap.L().Error("get best on year err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "03", errStr, 500)
	}
	bytes, err := draw.GenreSells(startYear, endYear, data)
	if err != nil {
		zap.L().Error("genreSells_draw_err", zap.String("traceID", traceID))
		return mainController.Error(ctx, "01", err.Error(), 500)
	}
	ctx.Response().Header.SetContentType("image/png")
	return ctx.Send(bytes)
}
