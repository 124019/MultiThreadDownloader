package main

import (
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"regexp"
	"time"
)
import (
	"MultiThreadDownloader/utils"
)

type Chunk struct {
	// no int
	st int
	ed int
}


func makeChunks(total int, step int) []Chunk {
	var result []Chunk
	for st := 1; st <= total; st += step {
		ed := st + step - 1
		if ed > total {
			ed = total
		}
		result = append(result, Chunk{
			// no: (st - 1) / step,
			st: st,
			ed: ed,
		})
	}
	return result
	// fmt.Println(result)
}

func get_file_info(headers map[string]string, url string) (int, string, error) {
	timeout_second := 20

	resp, StatusCode, Latency, err := utils.NetRequest(url, "HEAD", headers, nil, timeout_second)
	if err != nil {
		return 0, "", fmt.Errorf("download error: %v\n", err)
	}

	fmt.Printf("Latency: %d ms\n", Latency)
	fmt.Printf("status code: %d\n", StatusCode)
	fmt.Println(string(resp))

	var header map[string][]string
	err = json.Unmarshal(resp, &header)
	if err != nil {
		return 0, "", fmt.Errorf("unmarshal json data failed: %v\n", err)
	}
	ContentLength := header["Content-Length"][0]
	ContentDisposition := header["Content-Disposition"][0]

	totalSize, err := strconv.Atoi(ContentLength)
	if err != nil {
		return 0, "", fmt.Errorf("convert string to int failed: %v\n", err)
	}
	fmt.Printf("total size: %d bytes\n", totalSize) // file size(Bytes)
	re := regexp.MustCompile(`filename=\"(.+?)\"`)
	filename := re.FindStringSubmatch(ContentDisposition)[1]
	fmt.Printf("filename: %s\n", filename)
	return totalSize, filename, nil
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

	data, err = os.ReadFile("./BaiduNDApi/RunUrlHeader.json")
	if err != nil {
		fmt.Printf("read file error: %v\n", err)
	}
	headers := map[string]string{}
	err = json.Unmarshal(data, &headers)
	if err != nil {
		fmt.Printf("unmarshal json data failed:while reading headers, %v\n", err)
		return
	}
	
	totalSize, filename, err := get_file_info(headers, url)
	if err != nil {
		fmt.Printf("get file info error: %v\n", err)
		return
	}
	fmt.Printf("total size: %d bytes\n", totalSize)
	fmt.Printf("filename: %s\n", filename)

	chunk := makeChunks(totalSize, 2048)
	// fmt.Println(chunk)

	chunk0 := chunk[0]
	fmt.Printf("chunk0: %d - %d\n", chunk0.st, chunk0.ed)
	headers["Range"] = fmt.Sprintf("bytes=%d-%d", chunk0.st, chunk0.ed-1)
	start := time.Now()
	resp, StatusCode, _, err := utils.NetRequest(url, "GET", headers, nil, 30)
	if err != nil {
		fmt.Printf("download error: %v\n", err)
		return
	}
	Elapsed := time.Since(start)
	fmt.Printf("download time: %s\n", Elapsed)
	fmt.Printf("status code: %d\n", StatusCode)
	fmt.Println(string(resp))

}