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
)

var ErrorRequestTimeout = errors.New("request timeout")

func NetRequest(url string, method string, headers map[string]string, post_data interface{}, timeout_second int) ([]byte, int, time.Duration, error) { // Return : data , status code , elapsed time , error
	var data io.Reader
	switch method {
		case "POST", "PUT", "PATCH":
			jsonData, err := json.Marshal(post_data)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("marshal json data failed : while marshal post_data,  %w", err)
			}
			data = bytes.NewReader(jsonData)
		case "GET", "DELETE", "OPTIONS", "HEAD":
			data = nil
		default:
			return nil, 0, 0, fmt.Errorf("unsupported method: %s", method)
	}
	client := &http.Client{
		Timeout: time.Duration(timeout_second) * time.Second,
	}

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("new request error: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var StatusCode int
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, 0, time.Since(start), fmt.Errorf("%w: after %v (timeout=%d):%w", ErrorRequestTimeout, time.Since(start), timeout_second, err)
		}else {
			return nil, 0, time.Since(start), fmt.Errorf("Request Error: %w", err)
		}
	}
	if resp != nil { 
		StatusCode = resp.StatusCode
	} else { 
		return nil, 0, time.Since(start), fmt.Errorf("No valid Response, error?: %w", err)
	}
	defer resp.Body.Close()

	if method == "HEAD" {
		elapsed := time.Since(start)
		fmt.Printf("Request took %v\n", elapsed)
		headers ,err := json.Marshal(resp.Header)
		if err != nil {
			return nil, StatusCode, elapsed, fmt.Errorf("marshal json data failed: while marshal resp.Header, %w", err)
		}
		return headers, StatusCode, elapsed, nil
	}

	body, err := io.ReadAll(resp.Body)
	elapsed := time.Since(start)
	fmt.Printf("Request took %v\n", elapsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, 0, time.Since(start), fmt.Errorf("%w: after %v (timeout=%d):%w", ErrorRequestTimeout, time.Since(start), timeout_second, err)
		}else {
			return nil, 0, time.Since(start), fmt.Errorf("Request Error: %w", err)
		}
	}

	return body, StatusCode, elapsed, nil
}
