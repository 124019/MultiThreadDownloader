package core

import (
	"encoding/json"
	"fmt"
	"MultiThreadDownloader/component/base"
)

type FileListResponse struct{
	FileList   []FileListItem `json:"FileList"`
	FolderList []FileListItem `json:"FolderList"`
}
type FileListItem struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir int    `json:"isdir"`
	Fid   int    `json:"fid"`
	Size  int    `json:"size"`
}
type FileInfo struct {
	FsId     int    `json:"fs_id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	IsDir    int    `json:"isdir"`
	Dlink    string `json:"dlink"`
}

var baseurl = "http://127.0.0.1:1242/mtd/"
func get_file_list(dirPath string, platform string) (FileListResponse, error) { // Platform: "BaiduND" ...
	var respList FileListResponse
	FileListResp, err := base.NetRequest(baseurl+platform+"/filelist", "POST", map[string]string{"User-Agent": "Mozilla/5.0", "Content-Type": "application/json"}, map[string]interface{}{"ND_Path": dirPath}, 20)
	if err != nil {
		return respList, fmt.Errorf("get file list error: %v\n", err)
	}
	if err := json.Unmarshal(FileListResp.RespBody, &respList); err != nil {
		return respList, fmt.Errorf("unmarshal json data failed: %v\n", err)
	}
	return respList, nil
}

func get_file_info_api(fileList []FileListItem, platform string) ([]FileInfo, error) {
	var file_info []FileInfo
	FileInfoResp, err := base.NetRequest(baseurl+platform+"/directlink", "POST", map[string]string{"User-Agent": "Mozilla/5.0"}, fileList, 20)
	if err != nil {
		return nil, fmt.Errorf("get direct link error: %v\n", err)
	}
	if err := json.Unmarshal(FileInfoResp.RespBody, &file_info); err != nil {
		return nil, fmt.Errorf("unmarshal json data failed: %v\n", err)
	}
	return file_info, nil
}

func example() {
	resp, err := get_file_list("/", "BaiduND")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("File List: %+v\n", resp)
}