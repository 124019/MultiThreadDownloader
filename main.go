package main

import (
	"os"
	"fmt"
	"net/http"
	"io"
)

func getLocalText(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("Failed to read file : %s: %w", path, err)
    }
    return string(data), nil
}

func download() ([]byte, error) {
	client := &http.Client{}

	url, err := getLocalText("./BaiduNDApi/url.txt")
	if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}
	if url == "" {
		return nil , fmt.Errorf("url is empty.")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request error: %w", err)
	}

	req.Header.Set("User-Agent", "pan.baidu.com")

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
	resp, err := download()
	if err != nil {
		fmt.Printf("download error: %v\n", err)
		return
	}

	fmt.Println(string(resp))
}