package utils

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func HttpRequest(method, url string, request, response interface{}, headers map[string]string, encryption ...bool) error {

	dataBytes, err := jsoniter.Marshal(request)
	if err != nil {
		return err
	}
	var buff = &bytes.Buffer{}
	buff.Write(dataBytes)
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(resBytes, response)
	if err != nil {
		return err
	}
	return nil
}

func HttpFormRequest(url, fileName string, data []byte, response interface{}) error {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part1, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}
	_, err = part1.Write(data)
	if err != nil {
		return err
	}
	writer.Close()
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(resBytes, &response)
	if err != nil {
		return err
	}
	return nil
}
