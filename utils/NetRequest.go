package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"errors"
	"context"
	"os"
	"net"
)

type NetResp struct {
	RespBody []byte
	RespStatusCode int
	RespHeader http.Header
	ReqElapsed time.Duration
	ReqErr error
}

var ErrorRequestTimeout = errors.New("request timeout")
func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return os.IsTimeout(err)
}

func NetRequest(url string, method string, headers map[string]string, post_data interface{}, timeout_second int) NetResp { // Return : data , status code , elapsed time , error
	var data io.Reader
	switch method {
		case "POST", "PUT", "PATCH":
			jsonData, err := json.Marshal(post_data)
			if err != nil {
				return NetResp{ReqErr: fmt.Errorf("marshal json data failed : while marshal post_data,  %w", err)}
			}
			data = bytes.NewReader(jsonData)
		case "GET", "DELETE", "OPTIONS", "HEAD":
			data = nil
		default:
			return NetResp{ReqErr: fmt.Errorf("unsupported method: %s", method)}
	}
	client := &http.Client{
		Timeout: time.Duration(timeout_second) * time.Second,
	}

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return NetResp{ReqErr: fmt.Errorf("new request error: %w", err)}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var StatusCode int
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return NetResp{ReqElapsed: time.Since(start), ReqErr: fmt.Errorf("%w: after %v (timeout=%d):%w", ErrorRequestTimeout, time.Since(start), timeout_second, err)}
		}else {
			return NetResp{ReqElapsed: time.Since(start), ReqErr: fmt.Errorf("Request Error: %w", err)}
		}
	}
	if resp != nil { 
		StatusCode = resp.StatusCode
	} else { 
		return NetResp{ReqElapsed: time.Since(start), ReqErr: fmt.Errorf("No valid Response, error?: %w", err)}
	}
	defer resp.Body.Close()

	if method == "HEAD" {
		elapsed := time.Since(start)
		fmt.Printf("Request took %v\n", elapsed)
		return NetResp{RespStatusCode: StatusCode, RespHeader: resp.Header, ReqElapsed: elapsed, ReqErr: nil}
	}

	body, err := io.ReadAll(resp.Body)
	header := resp.Header
	elapsed := time.Since(start)
	fmt.Printf("Request took %v\n", elapsed)
	if err != nil {
		if isTimeout(err) {
			return NetResp{RespBody: body, RespStatusCode: StatusCode, RespHeader: header, ReqElapsed: elapsed, ReqErr: fmt.Errorf("%w: after %v (timeout=%d):%w", ErrorRequestTimeout, time.Since(start), timeout_second, err)}
		} else {
			return NetResp{RespHeader: header, ReqElapsed: time.Since(start), ReqErr: fmt.Errorf("Read Body Error: %w", err)}
		}
	}

	return NetResp{RespBody: body, RespStatusCode: StatusCode, RespHeader: header, ReqElapsed: elapsed, ReqErr: nil}
}
