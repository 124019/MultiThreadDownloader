package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
)

func NetRequest(url string, method string, headers map[string]string, post_data interface{}, timeout_second int) ([]byte, int, error) {
	var data io.Reader
	switch method {
		case "POST", "PUT", "PATCH":
			jsonData, err := json.Marshal(post_data)
			if err != nil {
				return nil, 0, fmt.Errorf("marshal json data failed : while marshal post_data,  %w", err)
			}
			data = bytes.NewReader(jsonData)
		case "GET", "DELETE", "OPTIONS", "HEAD":
			data = nil
		default:
			return nil, 0, fmt.Errorf("unsupported method: %s", method)
	}
	client := &http.Client{
		Timeout: time.Duration(timeout_second) * time.Second,
	}

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, 0, fmt.Errorf("new request error: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("get response error: %w", err)
	}
	defer resp.Body.Close()

	if method == "HEAD" {
		headers ,err := json.Marshal(resp.Header)
		if err != nil {
			return nil, resp.StatusCode, fmt.Errorf("marshal json data failed: while marshal resp.Header, %w", err)
		}
		return headers, resp.StatusCode, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read body error: %w", err)
	}

	return body, resp.StatusCode, nil
}
