package main

import (
	"MultiThreadDownloader/component/core"
	"MultiThreadDownloader/component/base"
	"encoding/json"
	"fmt"
	"os"
)

type Interval struct {
	st int
	ed int
}


func ChunksIndex(total_size int, chunk_size int) []Interval {
	var result []Interval
	for st_byte := 0; st_byte <= total_size; st_byte += chunk_size {
		ed_byte := st_byte + chunk_size - 1
		if ed_byte > total_size {
			ed_byte = total_size - 1 // first is 0, last is total_size - 1
		}
		result = append(result, Interval{
			// no: (st - 1) / step,
			st: st_byte,
			ed: ed_byte,
		})
	}
	// fmt.Println(result)
	return result
}

func bd_str_chunk_index_for(intervals []Interval) []string {
	result := make([]string, 0, len(intervals))
	for i := range intervals {
		Range := fmt.Sprintf("bytes=%d-%d", intervals[i].st, intervals[i].ed)
		result = append(result, Range)
	}
	// fmt.Println(result)
	return result
}

func main() {
	maxThread := 15
	chunkSize := 512 * 1024
	ND_Path := "/我的资源"
	platform := "BaiduND"
	// For TEST

	File_FolderList, err := core.GetFileList(ND_Path, platform)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}
	fmt.Println(File_FolderList) // For TEST

	fileList := File_FolderList.FileList // For TEST : Get all infomation in the list

	file_info_with_url, err := core.GetFileInfo(fileList, platform)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(file_info_with_url) // For TEST
	info0 := file_info_with_url[0] // For TEST : Get one file randomly
	url := info0.Dlink // For TEST : Get one file randomly
	totalSize := info0.Size // For TEST : Get one file randomly
	filename := info0.Filename // For TEST : Get one file randomly
	fmt.Printf("total size: %d bytes\n", totalSize)
	fmt.Printf("filename: %s\n", filename)

	baidu_header_data, err := os.ReadFile("./NDApi/BaiduNDApi/RunUrlHeader.json")
	if err != nil {
		fmt.Printf("read file error: %v\n", err)
	}
	headers := map[string]string{}
	err = json.Unmarshal(baidu_header_data, &headers)
	if err != nil {
		fmt.Printf("unmarshal json data failed:while reading headers, %v\n", err)
		return
	}

	// Get Latency
	response, err := base.NetRequest("https://d.pcs.baidu.com/", "GET", headers, nil, 30)
	if err != nil {
		fmt.Printf("get Latency error: %v\n", err)
		return
	}
	fmt.Printf("Latency: %d ms\n", response.ReqElapsed / 1000000)
	// Get Latency End

	str_range := bd_str_chunk_index_for(ChunksIndex(totalSize, chunkSize))
	// fmt.Println(chunk)
	length := len(str_range)
	fmt.Printf("total chunk: %d\n", length)

	err = core.MultiTGet(url, headers, maxThread, 5, str_range)
	if err != nil {
		fmt.Printf("multi-threaded download error: %v\n", err)
		return
	}
}