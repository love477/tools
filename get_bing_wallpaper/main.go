package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/love477/tools/common"
	"github.com/robfig/cron"
	"k8s.io/klog"
)

const (
	bingHost     = "https://cn.bing.com"
	homePageUrl  = "/?ensearch=1&FORM=BEHPTB"
	wallpaperUrl = "/th?id="
	suffix       = "UHD.jpg"
)

var directory = ""

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		klog.Errorln("get workDir error: ", err)
		return
	}
	klog.Info(workDir)
	directory = workDir + directory
	path := flag.String("path", directory, "存储下载的图片")
	flag.Parse()

	if path != nil && *path != "" {
		directory = *path
	}

	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
}

func GetBingWallpaper() (uri string, filename string) {
	uri, filename = "", ""
	url := fmt.Sprintf("%s%s", bingHost, homePageUrl)
	resp, err := http.Get(url)
	if err != nil {
		klog.Errorln("Get error: ", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Errorln("ReadAll error: ", err)
		return
	}
	// "Wallpaper":"/th?id=OHR.Migliarino_ROW2706888241_1920x1200.jpg\u0026rf=LaDigue_1920x1200.jpg"
	reg := regexp.MustCompile(`Wallpaper\":\"\/(th\?id=(.*.jpg)).*.jpg`)
	if reg == nil {
		klog.Errorln("regexp.MustCompile result is nil")
		return
	}

	res := reg.FindStringSubmatch(string(content))
	if len(res) < 3 {
		klog.Errorln("Error: reg error, res: ", res)
		return
	}
	file := res[2]
	l := strings.Split(strings.TrimSuffix(file, ".jpg"), "_")
	filename = fmt.Sprintf("%s_%s_%s", l[0], l[1], suffix)
	uri = fmt.Sprintf("%s%s%s", bingHost, wallpaperUrl, filename)
	return
}

func DownloadBingWallpaper() {
	uri, filename := GetBingWallpaper()
	common.Download(directory, filename, uri)
}

func main() {
	//DownloadBingWallpaper()
	c := cron.New()
	c.AddFunc("0 30 10 * * *", DownloadBingWallpaper)
	c.Start()

	select {}
}

// nohup ./getBingWallpaper -path=/Users/dingyuan.wu/Documents/wallpaper > /Users/dingyuan.wu/log/getBingWallpaper.log &
