package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

/*
 * http请求工具类
 */

// 默认超时时间
const REQ_TIMEOUT = 5 * time.Second

type options struct {
	requestUrl     string            // 请求URL
	requestParams  interface{}       // 请求参数
	requestMethod  string            // 请求方法 POST or GET
	basicAuth      string            // basic auth认证
	headers        map[string]string // 请求头参数
	requestTimeout time.Duration     // 请求超时时间
}

// 请求参数
type RequestOptions options

// 请求URL
func (opt *RequestOptions) Url(url string) *RequestOptions {
	opt.requestUrl = url
	return opt
}

// 请求参数
func (opt *RequestOptions) Params(params interface{}) *RequestOptions {
	opt.requestParams = params
	return opt
}

// 请求方法
func (opt *RequestOptions) Method(method string) *RequestOptions {
	opt.requestMethod = method
	return opt
}

// basic auth认证
func (opt *RequestOptions) BasicAuth(basicAuth string) *RequestOptions {
	opt.basicAuth = basicAuth
	return opt
}

// 请求超时时间
func (opt *RequestOptions) Timeout(timeout time.Duration) *RequestOptions {
	opt.requestTimeout = timeout
	return opt
}

// 请求头参数
func (opt *RequestOptions) Headers(headers map[string]string) *RequestOptions {
	opt.headers = headers
	return opt
}

/* http请求返回结果 */
func DoHttpExecute(result interface{}, ops *RequestOptions) error {
	return doHttpPost(result, ops)
}

// http请求
func doHttpPost(result interface{}, ops *RequestOptions) error {
	var err error
	var request *http.Request
	var method = ops.requestMethod

	if len(method) <= 0 {
		method = "POST"
	}

	if ops.requestParams != nil {
		bytesData, err := json.Marshal(ops.requestParams)
		if err != nil {
			return err
		}
		request, err = http.NewRequest(method, ops.requestUrl, bytes.NewReader(bytesData))
	} else {
		request, err = http.NewRequest(method, ops.requestUrl, nil)
	}

	if err != nil {
		return err
	}

	// set header
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if len(ops.basicAuth) > 0 {
		request.Header.Set("Authorization", ops.basicAuth)
	}
	if ops.headers != nil && len(ops.headers) > 0 {
		for key, val := range ops.headers {
			request.Header.Set(key, val)
		}
	}

	timeout := ops.requestTimeout
	if timeout.Nanoseconds() <= 0 {
		timeout = REQ_TIMEOUT
	}
	client := http.Client{Timeout: timeout}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBytes, result)
	if err != nil {
		return err
	}
	return nil
}
