package common

import (
	"io"
	"net/http"
	"os"

	"k8s.io/klog"
)

// DownloadBingWallpaper 下载文件存储到指定的目录
//	directory 文件存储目录
// 	filename 文件名，需要带文件类型后缀，如：test.jpg
// 	uri 文件的远端地址
func Download(directory, filename, uri string) {
	if uri == "" || filename == "" {
		klog.Warningln("Warn: uri or filename is empty")
		return
	}
	res, err := http.Get(uri)
	if err != nil {
		klog.Errorln("Error: get uri error: ", uri, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		klog.Errorln("Error: get uri error code: ", uri, res.StatusCode)
		return
	}
	_, statErr := os.Stat(directory + filename)
	if statErr == nil {
		klog.Errorln("Error: file is already exist: ", directory+filename)
		return
	}
	file, fileErr := os.Create(directory + filename)
	if fileErr != nil {
		klog.Errorln("Error: create file error: ", fileErr)
		return
	}
	defer file.Close()

	_, copyErr := io.Copy(file, res.Body)
	if copyErr != nil {
		klog.Errorln("Error: copy res.body to file error: ", err)
		return
	}
}
