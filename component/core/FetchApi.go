package core

import (
	"encoding/json"
	"fmt"
	"MultiThreadDownloader/component/base"
)

type File_FolderList struct{
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
func GetFileList(dirPath string, platform string) (File_FolderList, error) {
	var respList File_FolderList
	FileListResp, err := base.NetRequest(baseurl+platform+"/filelist", "POST", map[string]string{"User-Agent": "Mozilla/5.0", "Content-Type": "application/json"}, map[string]interface{}{"ND_Path": dirPath}, 20)
	if err != nil {
		return respList, fmt.Errorf("get file list error: %v\n", err)
	}
	if err := json.Unmarshal(FileListResp.RespBody, &respList); err != nil {
		return respList, fmt.Errorf("unmarshal json data (baidu file and folder list) failed: %v\n", err)
	}
	return respList, nil
}

func GetFileInfo(fileList []FileListItem, platform string) ([]FileInfo, error) {
	var file_info []FileInfo
	fileListInOne := map[string]any{ "FileList": fileList }
	FileInfoResp, err := base.NetRequest(baseurl+platform+"/directlink", "POST", map[string]string{"User-Agent": "Mozilla/5.0", "Content-Type": "application/json"}, fileListInOne, 20)
	if err != nil {
		return nil, fmt.Errorf("get direct link error: %v\n", err)
	}
	if err := json.Unmarshal(FileInfoResp.RespBody, &file_info); err != nil {
		return nil, fmt.Errorf("unmarshal json data failed (baidu file info with direct link): %v\n", err)
	}
	return file_info, nil
}

func example() {
	resp, err := GetFileList("/", "BaiduND")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("File List: %+v\n", resp)
}