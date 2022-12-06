package wallhaven

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/love477/tools/common"
	"k8s.io/klog"
)

const defaultMaxPage = 100
const directory = "/Users/dingyuan.wu/Downloads/wallhaven2/"

// const wallUri = "https://wallhaven.cc/search?categories=110&purity=110&sorting=hot&order=desc&page="
const wallUri = "https://wallhaven.cc/search?categories=110&purity=110&sorting=toplist&order=desc&page="

// const topUri = "https://wallhaven.cc/toplist?page="

// getToplistWallpaperId 获取wallhaven中toplist相关的图片ID
func getToplistWallpaperId(maxPage int) ([]string, error) {
	idList := []string{}
	for i := 1; i <= maxPage; i++ {
		if i%2 == 0 {
			fmt.Println("getToplistWallpaperId sleep 3 * time.Second")
			time.Sleep(3 * time.Second)
		}
		res, err := http.Get(fmt.Sprintf("%s%d", wallUri, i))
		if err != nil {
			klog.Error("get wallhaven toplist error: ", err)
			return idList, err
		}
		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			klog.Errorln("ReadAll error: ", err)
			return idList, err
		}
		// data-wallpaper-id=\"kx397d\" st.... href=\"https://wallhaven.cc/w/zygqvv\"
		reg := regexp.MustCompile(`data-wallpaper-id="(\S*)"`)
		if reg == nil {
			klog.Errorln("regexp.MustCompile result is nil")
			return idList, errors.New("regexp.MustCompile result is nil")
		}
		matchList := reg.FindAllSubmatch(content, 100)
		for i := range matchList {
			idList = append(idList, string(matchList[i][1]))
		}
	}
	return idList, nil
}

// getWallpaperUri 获取wallhaven图片的下载地址
func getWallpaperUri(wallpaperId string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://wallhaven.cc/w/%s", wallpaperId))
	if err != nil {
		klog.Error("get wallhaven toplist error: ", err)
		return "", err
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.Errorln("getWallpaperUri ReadAll error: ", err)
		return "", err
	}

	// img id=\"wallpaper\" src=\"https://w.wallhaven.cc/full/x6/wallhaven-x6567z.jpg\"
	reg := regexp.MustCompile(`<img id="wallpaper" src="(\S*)"`)
	if reg == nil {
		klog.Errorln("regexp.MustCompile result is nil")
		return "", errors.New("regexp.MustCompile result is nil")
	}
	matchList := reg.FindSubmatch(content)
	if len(matchList) <= 0 {
		klog.Errorln("not found uri, id: ", wallpaperId)
		return "", fmt.Errorf("not found uri, id: %s", wallpaperId)
	}
	return string(matchList[1]), nil
}

func download(ch chan string, done chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case uri := <-ch:
			// get filename by uri
			l := strings.Split(uri, "/")
			filename := l[len(l)-1]
			common.Download(directory, filename, uri)
			wg.Done()
		case <-done:
			fmt.Println("done")
			return
		}
	}
}

// DownloadWallhavenToplist 下载wallhaven中toplist的图片
func DownloadWallhavenToplist(maxPage int) {
	if maxPage <= 0 {
		maxPage = defaultMaxPage
	}

	wallpaperIdList, err := getToplistWallpaperId(maxPage)
	if err != nil {
		klog.Error("getToplistWallpaperId failed, error: ", err)
		return
	}

	wg := sync.WaitGroup{}
	ch := make(chan string)
	done := make(chan struct{})

	go func(wg *sync.WaitGroup) {
		download(ch, done, wg)
	}(&wg)

	for i, v := range wallpaperIdList {
		if i%2 == 0 {
			fmt.Println("getWallpaperUri sleep 3 * time.Second")
			time.Sleep(3 * time.Second)
		}
		wg.Add(1)
		uri, err := getWallpaperUri(v)
		if err != nil {
			klog.Error("getWallpaperUri error: ", err)
			wg.Done()
			continue
		}
		ch <- uri
	}

	wg.Wait()
	done <- struct{}{}
	time.Sleep(3 * time.Second)
	fmt.Println("DownloadWallhavenToplist successful")
}
