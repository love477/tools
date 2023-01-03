# docker rmi
删除本地所有没有使用的docker image。

## 使用方法
### 直接使用提供的可执行文件
下载仓库中的可执行文件：tools/docker_rmi docker_rmi，本地直接运行
ps：本仓库中提供的可执行文件是mac Intel Core 64位系统的
```sh
# terminal进入对应文件的目录后执行下面的指令
./docker_rmi
```
### 自己编译源文件
编译产生可执行文件并执行：
1. clone仓库到本地
2. 进入tools目录，执行go build -o docker_rmi/docker_rmi docker_rmi/main.go
3. 执行./docker_rmi/docker_rmi

安装docker_rmi指令到本地：
1. clone仓库到本地
2. 进入docker_rmi目录，执行go install
3. 启动terminal，执行docker_rmi

## 注意点
1. 需要确保本地的docker是运行状态
2. 本指令不能删除正在被使用的镜像