package core

import (
	"encoding/json"
	"fmt"
	"MultiThreadDownloader/component/base"
)

type FileListItem struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir int    `json:"isdir"`
	Fid   int    `json:"fs_id"`
	Size  int    `json:"size"`
}
type FileInfo struct {
	FsId     int    `json:"fs_id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	IsDir    int    `json:"isdir"`
	Dlink    string `json:"dlink"`
}

func get_file_list(dirPath string, headers map[string]string, platform string) ([]FileListItem, error) { // Platform: "BaiduND" ...
	FileList := []FileListItem{}
	FileListResp, err := base.NetRequest("http://127.0.0.1:1242/"+platform+"/filelist", "POST", headers, map[string]interface{}{"ND_Path": dirPath}, 20)
	if err != nil {
		return nil, fmt.Errorf("get file list error: %v\n", err)
	}
	if err := json.Unmarshal(FileListResp.RespBody, &FileList); err != nil {
		return nil, fmt.Errorf("unmarshal json data failed: %v\n", err)
	}
	return FileList, nil
}

func get_file_info_api(fileList []FileListItem, headers map[string]string, platform string) ([]FileInfo, error) {
	FileInfo := []FileInfo{}
	FileInfoResp, err := base.NetRequest("http://127.0.0.1:1242/"+platform+"/directlink", "POST", headers, fileList, 20)
	if err != nil {
		return nil, fmt.Errorf("get direct link error: %v\n", err)
	}
	if err := json.Unmarshal(FileInfoResp.RespBody, &FileInfo); err != nil {
		return nil, fmt.Errorf("unmarshal json data failed: %v\n", err)
	}
	return FileInfo, nil
}