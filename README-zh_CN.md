[English](README.md) | 简体中文
## HttpCat 概述
一个基于HTTP的文件传输瑞士军刀。
HttpCat专注于利用HTTP/HTTPS协议进行简单、高效、稳定的文件上传和下载。

HttpCat是一个可靠、高效、易用的HTTP文件传输瑞士军刀,它将大大提高你的文件传输控制力和体验。无论是临时分享还是批量传输文件,HttpCat都将是你的优秀助手。

## 功能特点
- 简单
- 无依赖，可移植性好

## 编译
chmod +x build.sh
./build.sh

生成可执行文件： output/httpcat

cp output/httpcat /usr/local/bin/
httpcat -h

配置文件：
mkdir -p /etc/httpdcat
cp server/conf/svr.yml /etc/httpdcat/svr.yml

可以利用tmux方式后台运行:
cd /root
Create a new tmux session using a socket file named tmux_httpcat
$ tmux -S tmux_httpcat


Move process to background by detaching
Ctrl+b d OR ⌘+b d (Mac)

To re-attach
$ tmux -S tmux_httpcat attach

Alternatively, you can use the following single command to both create (if not exists already) and attach to a session:
$ tmux new-session -A -D -s tmux_httpcat

To delete farming session
$ tmux kill-session -t tmux_httpcat

### 配置开机自启动



### 无鉴权直接访问上传文件
当我想直接访问访问上传的文件，也不需要鉴权场景，我们可以在启动参数中指定静态资源目录为上传目录，这样就可以直接访问上传的文件了。
例如：
go run cmd/httpcat.go --static=/home/web/website/upload/  --c server/conf/svr.yml


## 使用
### 使用curl上传文件
注意： f1 为服务端代码定义的，修改为其他，如file，会报错上传失败。
```bash
# curl -vF "f1=@/root/test.lz4" http://localhost:8888/api/upload
*   Trying 127.0.0.1:80...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 80 (#0)
> POST /upload HTTP/1.1
> Host: localhost
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Length: 734
> Content-Type: multipart/form-data; boundary=------------------------1538dd9d9ac92293
>
* We are completely uploaded and fine
* Mark bundle as not supporting multiuse
  < HTTP/1.1 201 Created
  < Content-Type: text/plain; charset=utf-8
  < Date: Tue, 07 Nov 2023 07:46:18 GMT
  < Content-Length: 19
  <
  upload successful
```

###  下载文件

## 
指定静态资源目录为上传目录，这样就可以直接访问上传的文件了。
go run cmd/httpcat.go --static=/home/web/website/upload/  --c server/conf/svr.yml

## 提交代码，检查git 用户名和邮箱
使用以下命令来查看全局配置
git config --global user.name
git config --global user.email

查看当前仓库配置
git config user.name
git config user.email