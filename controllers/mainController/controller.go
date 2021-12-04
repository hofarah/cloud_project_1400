package mainController

import (
	"cloudProject/apiSchema/common"
	"cloudProject/pkg/cast"
	"cloudProject/pkg/prometheus"
	"cloudProject/pkg/utils"
	translate "cloudProject/statics/translate/message"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

var serviceStats map[string]*prometheus.Stats

type Request interface {
	Validate(ctx *fiber.Ctx) (string, int, error)
}

func ParseBody(ctx *fiber.Ctx, request Request) (string, int, error) {
	err := jsoniter.Unmarshal(ctx.Body(), request)
	if err != nil {
		return "02", 400, err
	}
	return request.Validate(ctx)
}
func prepareBody(ctx *fiber.Ctx, body interface{}) error {
	data, err := jsoniter.Marshal(body)
	if err != nil {
		return err
	}
	_, _ = ctx.Write(data)
	return nil
}
func Response(ctx *fiber.Ctx, data interface{}) error {
	res := common.Response{}
	lang := GetAcceptLanguage(ctx)
	res.Status = "success"
	message := ""
	if utils.IsNonEmptyString(data) {
		message, _ = cast.ToString(data)
	}

	if message != "" {
		msg := translate.GetMessage(message, lang)
		if msg != "" {
			message = msg
		}
		data = map[string]interface{}{
			"message": message,
		}
		res.Data = data
	} else {
		res.Data = data
	}
	ctx.Status(200)
	getAPIStats(ctx.OriginalURL()).AddSuccess()
	return prepareBody(ctx, res)
}
func Error(ctx *fiber.Ctx, section, errStr string, code int, msg ...string) error {
	res := common.Response{}
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	errCode := GetAPIBaseErrCode(ctx) + section + errStr
	res.Status = "error"
	lang := GetAcceptLanguage(ctx)

	res.Error.Code = errCode
	errMessage := translate.GetMessage(errCode, lang)
	res.Error.Message = errMessage
	if res.Error.Message == "" {
		if message == "" {
			res.Error.Message = translate.GetMessage(strconv.Itoa(code), lang)
		} else {
			msg := translate.GetMessage(message, lang)
			if msg != "" {
				res.Error.Message = msg
			} else if code == 200 {
				res.Error.Message = message
			}
		}
		if res.Error.Message == "" {
			res.Error.Message = translate.GetMessage(strconv.Itoa(500), lang)
		}
	}
	ctx.Status(code)
	getAPIStats(ctx.OriginalURL()).AddError()
	return prepareBody(ctx, res)
}
