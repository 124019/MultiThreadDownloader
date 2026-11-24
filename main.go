package main

import (
	"MultiThreadDownloader/src/core"
	"MultiThreadDownloader/src/core/netdisk"
	"MultiThreadDownloader/src/lib"
	"MultiThreadDownloader/src/lib/fetchApi"
	"encoding/json"
	"fmt"
	"os"
	"time"
)


var maxDownloadThread = 15
var chunkSize = 256 * 1024
var ND_Path = "/我的资源/爱言叶5动态素材包" // root path : "/" 
var platform = "BaiduNetdisk"
var RetryTimeInterval = 3 * time.Second
var maxWritingThread = 15
// For TEST, it will be post by frontend, but now we just hardcode it here.


func main() {
	File_FolderList, err := fetchApi.GetFileList(ND_Path, platform)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}
	fmt.Println(File_FolderList) // For TEST

	fileList := File_FolderList.FileList // Get all infomation in the list

	// Post fileList to frontend
	var select_file_list []fetchApi.FileListItem
	select_file_list = append(select_file_list, fileList[2]) // Responsed from frontend
	fmt.Printf("select_file_list: %v", select_file_list)

	file_info_with_url, err := fetchApi.GetFileInfo(fileList, platform)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(file_info_with_url)
	info0 := file_info_with_url[0] // Should get all info, but now ignore it tempory to develop other method.
	url := info0.Dlink
	totalSize := info0.Size
	filename := info0.Filename
	
	fmt.Printf("total size: %d bytes\n", totalSize)
	fmt.Printf("filename: %s\n", filename)

	baidu_header_data, err := os.ReadFile("./NDApi/static/headers/Baidu_fetchURL_header.json")
	if err != nil {
		fmt.Printf("read file error: %v\n", err)
	}
	headers := map[string]string{}
	err = json.Unmarshal(baidu_header_data, &headers)
	if err != nil {
		fmt.Printf("unmarshal json data failed:while reading headers, %v\n", err)
		return
	}

	response, err := core.NetRequest("https://d.pcs.baidu.com/", "GET", headers, nil, 30) // Get Latency
	if err != nil {
		fmt.Printf("get Latency error, try to check your internet: %v\n", err)
		return
	}
	fmt.Printf("Latency: %d ms\n", response.ReqElapsed / 1000000)

	chunk_range := netdisk.BaiduNetdiskChunkRangeInit(core.InitChunksIndex(totalSize, chunkSize))
	// fmt.Println(chunk)
	chunk_range_len := len(chunk_range)
	fmt.Printf("total chunk number: %d\n", chunk_range_len)

	f, err := lib.CreatePlaceholderFile("./download/" + filename, totalSize) // create file to write
	if err != nil {
		fmt.Printf("create placeholder file error: %v\n", err)
		return
	}
	fmt.Printf("placeholder file created successfully: %s\n", f.F.Name())

	ok, errs := lib.MultiTGet(url, chunk_range, lib.MultiGetConfig{WriteTo: *f, RetryTimeInterval: RetryTimeInterval, MaxRetry: 2, MaxDownloadThread: maxDownloadThread, MaxWritingThread: maxWritingThread, Headers: headers, SingleChunkTimeout: 5})
	if len(errs) == 0 && ok != true {
		for _, e := range errs {
			fmt.Printf("ERRORS\n=====\n %v\n", e)
		}
		return
	}
}