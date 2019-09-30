package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

/*
 * http请求客户端
 */

// 默认超时时间
const REQ_TIMEOUT = 5 * time.Second

// 请求参数
type RequestOptions struct {
	RequestUrl     string        // 请求URL
	RequestParams  interface{}   // 请求参数
	RequestMethod  string        // 请求方法 get or post
	BasicAuth      string        // basic auth认证
	RequestTimeout time.Duration // 请求超时时间
	ResultData     interface{}   // 返回结果
}

/* http请求返回结果 */
func DoExecute(ops RequestOptions) error {
	return doHttpPost(ops.ResultData, ops.RequestParams, ops.BasicAuth, ops.RequestUrl, ops.RequestMethod, ops.RequestTimeout)
}

// http请求
func doHttpPost(result interface{}, params interface{}, basicAuth string, url string, method string, timeout time.Duration) error {
	var err error
	var request *http.Request
	if params != nil {
		bytesData, err := json.Marshal(params)
		if err != nil {
			return err
		}
		request, err = http.NewRequest(method, url, bytes.NewReader(bytesData))
	} else {
		request, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if len(basicAuth) > 0 {
		request.Header.Set("Authorization", basicAuth)
	}
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
