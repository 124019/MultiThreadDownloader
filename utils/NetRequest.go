package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
)

func NetRequest(url string, method string, headers map[string]string, post_data interface{}, timeout_second int) ([]byte, int, time.Duration, error) { // Return : data , status code , elapsed time/ns , error
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

	start := time.Now()
	resp, err := client.Do(req)
	var StatusCode int
	if resp != nil { 
		StatusCode = resp.StatusCode
	} else { 
		return nil, 0, time.Since(start), err
	}
	if err != nil {
		return nil, StatusCode, 0, fmt.Errorf("get response error: %w", err)
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
		return nil, StatusCode, elapsed, fmt.Errorf("read body error: %w", err)
	}

	return body, StatusCode, elapsed, nil
}
