package main

import (
	"fmt"
	"os"
	"encoding/json"
	"strconv"
	"regexp"
)
import (
	"MultiThreadDownloader/utils"
)

type Interval struct {
	st int
	ed int
}


func makeChunks(total int, step int) []Interval {
	var result []Interval
	for st := 0; st <= total; st += step {
		ed := st + step - 1
		if ed > total {
			ed = total - 1 // first is 0, last is tt - 1
		}
		result = append(result, Interval{
			// no: (st - 1) / step,
			st: st,
			ed: ed,
		})
	}
	return result
	// fmt.Println(result)
}

func str_chunk(intervals []Interval) []string {
	result := make([]string, 0, len(intervals))
	for i := range intervals {
		Range := fmt.Sprintf("bytes=%d-%d", intervals[i].st, intervals[i].ed)
		result = append(result, Range)
	}
	// fmt.Println(result)
	return result
}

func get_file_info(headers map[string]string, url string) (int, string, error) {
	timeout_second := 20

	resp, StatusCode, _, err := utils.NetRequest(url, "HEAD", headers, nil, timeout_second)
	if err != nil {
		return 0, "", fmt.Errorf("download error: %v\n", err)
	}

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

	// Get Latency
	_, _, Latency, err := utils.NetRequest("https://d.pcs.baidu.com/", "GET", headers, nil, 30)
	if err != nil {
		fmt.Printf("get Latency error: %v\n", err)
		return
	}
	fmt.Printf("Latency: %d ms\n", Latency / 1000000)
	// Get Latency End

	str_range := str_chunk(makeChunks(totalSize, 150*1024))
	// fmt.Println(chunk)
	length := len(str_range)
	fmt.Printf("total chunk: %d\n", length)

	str_range0 := str_range[0]
	fmt.Printf("range: %s\n", str_range0)
	headers["Range"] = str_range0
	_, StatusCode, time_cost, err := utils.NetRequest(url, "GET", headers, nil, 30)
	if err != nil {
		fmt.Printf("download error: %v\n", err)
		return
	}
	time_elapsed := (time_cost - Latency)
	fmt.Printf("download time: %s\n", time_elapsed)
	fmt.Printf("status code: %d\n", StatusCode)
	// fmt.Println(string(resp))

}