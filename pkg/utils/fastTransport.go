package utils

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type FastTransport struct{}

func (t *FastTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	freq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(freq)
	fres := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(fres)

	if req.Method == "GET" && req.Body != nil {
		req.Method = "POST"
	}
	freq.SetHost(req.Host)
	freq.SetRequestURI(req.URL.String())
	freq.Header.SetRequestURI(req.URL.String())
	freq.Header.SetMethod(req.Method)
	for k, vv := range req.Header {
		for _, v := range vv {
			freq.Header.Set(k, v)
		}
	}
	reqID, _ := NewUUID()
	zap.L().Debug("request started", zap.String("reqID", reqID))
	if req.Body != nil {
		freq.SetBodyStream(req.Body, -1)
	}
	err := fasthttp.Do(freq, fres)
	if err != nil {
		zap.L().Debug("request ended", zap.String("reqID", reqID), zap.Error(err), zap.String("body", string(freq.Body())))
		return nil, err
	}

	zap.L().Debug("request ended", zap.String("reqID", reqID), zap.String("body", string(freq.Body())))
	res := &http.Response{Header: make(http.Header)}

	res.StatusCode = fres.StatusCode()
	fres.Header.VisitAll(func(k, v []byte) {
		res.Header.Set(string(k), string(v))
	})
	res.Body = ioutil.NopCloser(strings.NewReader(string(fres.Body())))
	return res, nil
}
