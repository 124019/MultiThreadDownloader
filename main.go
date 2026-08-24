package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"encoding/json"
	"bytes"
)

func getLocalText(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("Failed to read file : %s: %w", path, err)
    }
    return string(data), nil
}

func requestLink(url string, method string, headers map[string]string, post_data interface{}, timeout_second int) ([]byte, error) {
	var data io.Reader
	switch method {
		case "POST", "PUT", "PATCH":
			jsonData, err := json.Marshal(post_data)
			if err != nil {
				return nil, fmt.Errorf("marshal json data failed: %w", err)
			}
			data = bytes.NewReader(jsonData)
		case "GET", "DELETE", "HEAD", "OPTIONS":
			data = nil
		default:
			return nil, fmt.Errorf("unsupported method: %s", method)
	}

	client := &http.Client{
		Timeout: time.Duration(timeout_second) * time.Second,
	}
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get response error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
	}

	return body, nil
}

func main() {
	data, err := os.ReadFile("./BaiduNDApi/url.txt")
	if err != nil {
		fmt.Printf("read file error: %v\n", err)
	}
	url := string(data)
	if url == "" {
		fmt.Printf("url file is empty.")
	}

	headers := map[string]string{ "User-Agent": "pan.baidu.com" }

	timeout_second := 20

	resp, err := requestLink(url, "GET", headers, nil, timeout_second)
	if err != nil {
		fmt.Printf("download error: %v\n", err)
		return
	}

	fmt.Println(string(resp))
}