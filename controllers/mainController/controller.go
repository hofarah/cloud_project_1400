package mainController

import (
	"cloudProject/apiSchema/common"
	"cloudProject/pkg/cast"
	"cloudProject/pkg/logger"
	"cloudProject/pkg/prometheus"
	"cloudProject/pkg/utils"
	translate "cloudProject/statics/translate/message"
	"context"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/opentracing/opentracing-go"
	goPrometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
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
	ctx.Set("content-type", fiber.MIMEApplicationJSONCharsetUTF8)
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
	stats := getAPIStats(ctx.OriginalURL())
	if stats != nil {
		stats.AddSuccess()
	}
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
	stats := getAPIStats(ctx.OriginalURL())
	if stats != nil {
		stats.AddError()
	}
	return prepareBody(ctx, res)
}
func SetAPIBaseErrCode(ctx *fiber.Ctx, baseErrCode string) {
	ctx.Locals("baseErrCode", baseErrCode)
}
func GetAPIBaseErrCode(ctx *fiber.Ctx) string {
	base, _ := cast.ToString(ctx.Locals("baseErrCode"))
	return base
}
func GetAcceptLanguage(ctx *fiber.Ctx) string {
	return ctx.Get(fiber.HeaderAcceptLanguage, "fa")
}
func GetParentSpan(ctx *fiber.Ctx) opentracing.Span {
	var parentSpan opentracing.Span
	span := ctx.Locals("parentSpan")
	if span != nil {
		ok := false
		parentSpan, ok = span.(opentracing.Span)
		if ok {
			return parentSpan
		}
	}
	return nil
}
func GetTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}
func StartAPISpan(ctx *fiber.Ctx, apiName string) (apiSpan opentracing.Span, traceID string) {
	parentSpan := GetParentSpan(ctx)
	tracer := GetTracer()
	if parentSpan != nil && tracer != nil {
		apiSpan = tracer.StartSpan(apiName, opentracing.ChildOf(parentSpan.Context()))
		traceID = logger.GetTraceIDFromSpan(apiSpan)
		ctx.Locals("apiSpan", apiSpan)
	}
	return apiSpan, traceID
}
func FinishAPISpan(ctx *fiber.Ctx) {
	apiSpan := GetAPISpan(ctx)
	logger.FinishSpan(apiSpan)
	parentSpan := GetParentSpan(ctx)
	logger.FinishSpan(parentSpan)
}
func GetAPISpan(ctx *fiber.Ctx) opentracing.Span {
	span, ok := ctx.Locals("apiSpan").(opentracing.Span)
	if !ok {
		return nil
	}
	return span
}
func GetContextWithSpan(ctx *fiber.Ctx) context.Context {
	cntx := opentracing.ContextWithSpan(ctx.Context(), GetAPISpan(ctx))
	return cntx
}
func InitAPI(ctx *fiber.Ctx, baseErrCode, apiName string) (context.Context, string) {
	SetAPIBaseErrCode(ctx, baseErrCode)
	_, traceID := StartAPISpan(ctx, apiName)
	return GetContextWithSpan(ctx), traceID
}
func StartPrometheus(routes map[string]string) {
	metrics := make([]prometheus.Metric, len(routes))
	statistics := make([]*prometheus.Stats, len(routes))
	serviceStats = make(map[string]*prometheus.Stats)
	i := 0
	for name, url := range routes {
		metrics[i] = prometheus.Metric{Name: name}
		stats := &prometheus.Stats{}
		statistics[i] = stats
		serviceStats[strings.ToLower(url)] = stats
		i++
	}
	collector := prometheus.NewCollector()
	collector.BulkAddAPIMetric(metrics, statistics)
	go func() {
		goPrometheus.MustRegister(collector)
		http.Handle("/metrics", promhttp.Handler())
		zap.L().Info("metrics to serve on port :8080")
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			panic(err)
		}
	}()
}
func getAPIStats(url string) *prometheus.Stats {
	url = strings.ToLower(url)
	if strings.Contains(url, "?") {
		url = url[:strings.Index(url, "?")]
	}
	if _, found := serviceStats[url]; found {
		return serviceStats[url]
	}
	return nil
}
