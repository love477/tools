package main

import (
	"flag"
	"io/fs"
	"log"
	"os"

	"github.com/love477/tools/wallhaven/pkg"
)

func main() {
	var directory string
	var maxPage int
	flag.StringVar(&directory, "directory", "", "保存文件的目录,默认是当前目录下新建wallhaven目录")
	flag.IntVar(&maxPage, "maxpage", 10, "需要下载的页数,默认下载前10页")

	flag.Parse()

	if directory == "" {
		err := os.Mkdir("wallhaven", fs.ModeDir)
		if err != nil && !os.IsExist(err) {
			log.Fatal("os.Mkdir err: ", err)
		}
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal("os.Getwd err: ", err)
		}
		directory = dir + "/wallhaven"
	} else {
		// check directory is exist
		_, err := os.Stat(directory)
		if err != nil {
			log.Fatal("os.Stat err: ", err)
		}
	}

	pkg.DownloadWallhavenToplist(maxPage, directory)
}
