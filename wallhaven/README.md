# download wallhaven wallpaper
下载wallhaven中的toplist的壁纸
## 使用方法
### 直接使用提供的可执行文件
下载仓库中的可执行文件：tools/wallhaven download_wallhaven_wallpaper，本地直接运行
ps：本仓库中提供的可执行文件是mac Intel Core 64位系统的
```sh
# terminal进入对应文件的目录后执行下面的指令，查看帮助
./download_wallhaven_wallpaper --help

# 开始下载
./download_wallhaven_wallpaper --directory="/Users/xxx/Downloads" --maxpage=11
```

### 自己编译源文件
编译产生可执行文件并执行：
1. clone仓库到本地
2. 进入tools目录，执行go build -o wallhaven/download_wallhaven_wallpaper wallhaven/main.go
3. 执行./wallhaven/download_wallhaven_wallpaper --help 查看帮助
4. 执行./wallhaven/download_wallhaven_wallpaper --directory="/Users/xxx/Downloads" --maxpage=11 开始下载文件

安装docker_rmi指令到本地：
1. clone仓库到本地
2. 进入wallhaven目录，执行go install
3. 启动terminal，执行wallhaven --help 查看帮助
4. 执行wallhaven --directory="/Users/xxx/Downloads"  --maxpage=11开始下载文件
