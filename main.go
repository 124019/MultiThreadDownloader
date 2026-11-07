package main

import (
	"MultiThreadDownloader/src/core"
	"MultiThreadDownloader/src/core/netdisk"
	"MultiThreadDownloader/src/lib"
	"MultiThreadDownloader/src/lib/fetchApi"
	"encoding/json"
	"fmt"
	"os"
)


var maxThread = 15
var chunkSize = 512 * 1024
var ND_Path = "/我的资源"
var platform = "BaiduNetdisk"
// For TEST, it will be post by frontend, but now we just hardcode it here.


func main() {
	File_FolderList, err := fetchApi.GetFileList(ND_Path, platform)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 
	}
	fmt.Println(File_FolderList) // For TEST

	fileList := File_FolderList.FileList // For TEST : Get all infomation in the list

	file_info_with_url, err := fetchApi.GetFileInfo(fileList, platform)
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

	// Get Latency
	response, err := core.NetRequest("https://d.pcs.baidu.com/", "GET", headers, nil, 30)
	if err != nil {
		fmt.Printf("get Latency error: %v\n", err)
		return
	}
	fmt.Printf("Latency: %d ms\n", response.ReqElapsed / 1000000)
	// Get Latency End

	str_range := netdisk.BaiduNetdiskChunkRangeInit(core.InitChunksIndex(totalSize, chunkSize))
	// fmt.Println(chunk)
	length := len(str_range)
	fmt.Printf("total chunk: %d\n", length)

	err = lib.MultiTGet(url, headers, maxThread, 5, str_range)
	if err != nil {
		fmt.Printf("multi-threaded download error: %v\n", err)
		return
	}
}