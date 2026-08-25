package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"encoding/json"
	"bytes"
	"strconv"
	"regexp"
)

func requestLink(url string, method string, headers map[string]string, post_data interface{}, timeout_second int) ([]byte, int, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout_second) * time.Second,
	}

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

	resp, StatusCode, err := requestLink(url, "HEAD", headers, nil, timeout_second)
	if err != nil {
		fmt.Printf("download error: %v\n", err)
		return
	}

	fmt.Printf("status code: %d\n", StatusCode)
	fmt.Println(string(resp))

	var header map[string][]string
	err = json.Unmarshal(resp, &header)
	if err != nil {
		fmt.Printf("unmarshal json data failed: %v\n", err)
		return
	}
	ContentLength := header["Content-Length"][0]
	ContentDisposition := header["Content-Disposition"][0]

	totalSize, err := strconv.Atoi(ContentLength)
	if err != nil {
		fmt.Printf("convert string to int failed: %v\n", err)
		return
	}
	fmt.Printf("total size: %d bytes\n", totalSize) // file size(Bytes)

	re := regexp.MustCompile(`filename=\"(.+?)\"`)
	filename := re.FindStringSubmatch(ContentDisposition)[1]
	fmt.Printf("filename: %s\n", filename)
}