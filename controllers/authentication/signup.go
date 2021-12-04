package authentication

import (
	"cloudProject/apiSchema/auth"
	"cloudProject/controllers/mainController"
	"cloudProject/models/user"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func SignUp(ctx *fiber.Ctx) error {
	spanCtx, traceID := mainController.InitAPI(ctx, "1000", "singUp")
	defer mainController.FinishAPISpan(ctx)

	request := auth.SignUpRequest{}
	errStr, code, err := mainController.ParseBody(ctx, &request)
	if err != nil {
		return mainController.Error(ctx, "01", errStr, code, err.Error())
	}
	signUPUser, token, errStr, err := user.Repo.SignUPUser(spanCtx, request.UserName)
	if err != nil {
		zap.L().Error("SignUPUser err", zap.String("traceID", traceID), zap.Error(err))
		return mainController.Error(ctx, "02", errStr, 200, err.Error())
	}
	res := auth.SignUpResponse{
		Token:  token,
		Secret: signUPUser.Secret,
	}
	return mainController.Response(ctx, res)
}
