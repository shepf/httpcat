#!/bin/bash

if ! type go >/dev/null 2>&1; then
    echo 'go not installed';
    exit 1
fi

# Clean output directory
rm -rf release
mkdir -p release
go mod tidy

# 检查前端文件是否存在
if [ ! -f "dist.zip" ]; then
  echo "dist.zip 文件不存在，请先执行 npm run build 命令,打包前端,再执行 build.sh"
  exit 1
fi
cp dist.zip release/ -rf


# 构建 Linux x86 版本
echo "Building for Linux x86"

HTTPCAT_VERSION=v0.1.5
HTTPCAT_BUILD=$(date "+%Y%m%d%H%M")
COMMIT_ID=$(git rev-parse HEAD)
GOOS=linux GOARCH=amd64 go build \
 -ldflags "-s -w" \
 -ldflags "-X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
 -o ./release/httpcat-linux-x86 ./cmd/httpcat.go

# 构建 Linux ARM 版本
echo "Building for Linux ARM"

# 需要添加 CGO_ENABLED=1
# 在编译 Linux x86 架构的版本时，你没有遇到 CGO 的问题，因为在 x86 架构上，默认情况下，CGO 是启用的。
#不同的架构对 CGO 的支持情况是不同的。在 x86 架构上，通常默认启用 CGO，而在 ARM 架构上，默认情况下是禁用 CGO。这就是为什么在编译 ARM 架构的版本时，需要显式启用 CGO。
# go-sqlite3 包是一个 Cgo 包，它依赖于 CGO 来与 SQLite 库进行交互。因此，当你禁用 CGO 后，go-sqlite3 将无法正常工作

# 另外：如果你在 x86 架构的计算机上编译 ARM 架构的程序，并且使用了 CGO，你需要安装适用于 ARM 架构的交叉编译工具链。
#交叉编译工具链是用于在一个平台上构建另一个平台的工具集。在这种情况下，你需要安装适用于 ARM 架构的交叉编译工具链，例如 gcc-aarch64-linux-gnu。
#
sudo apt-get install gcc-aarch64-linux-gnu -y
CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build \
 -ldflags "-s -w" \
 -ldflags "-X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
 -o ./release/httpcat-linux-aarch64 ./cmd/httpcat.go

## 构建 Windows x86 版本
#echo "Building for Windows x86"
#echo "Building for Windows"
#GOOS=windows GOARCH=amd64 go build \
# -ldflags "-s -w" \
# -ldflags "-X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
# -o ./release/httpcat.exe ./cmd/httpcat.go


# 修改源码 README.md、translations/README-cn.md 文件中的 httpcat_version="v0.x.x" 为当前版本号
sed -i "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" README.md
sed -i "s/httpcat_version=\".*\"/httpcat_version=\"$HTTPCAT_VERSION\"/g" translations/README-cn.md

# Package configuration file and static files
cp -r server/conf release/
cp -r static release/
cp -r httpcat.service release/
cp -r README.md release/
mkdir -p release/translations
cp -rf translations/* release/translations/
# copy install.sh
cp -r install.sh release/
cp -r uninstall.sh release/
chmod +x release/install.sh
chmod +x release/uninstall.sh

# Create release archive for Linux
tar zcvf httpcat_$HTTPCAT_VERSION.tar.gz release/*
mv httpcat_$HTTPCAT_VERSION.tar.gz release/

# Return to the root directory
cd ..

echo "Build complete"


config_systemd(){

      if [ -f /usr/bin/systemctl ]; then
        #For CentOS7
        #check and add httpcat service
        rm -rf /usr/lib/systemd/system/httpcat.service
        cp httpcat.service  /usr/lib/systemd/system/ -f

        systemctl is-enabled httpcat
        if [ $? -ne 0 ]; then
            systemctl enable httpcat
        fi


      fi


}

