#!/bin/bash

if ! type go >/dev/null 2>&1; then
    echo 'go not installed';
    exit 1
fi

# Clean output directory
rm -rf release
mkdir -p release
go mod tidy

# Build for Linux
echo "Building for Linux"

HTTPCAT_VERSION=v0.1.3
HTTPCAT_BUILD=$(date "+%Y%m%d%H%M")
COMMIT_ID=$(git rev-parse HEAD)
GOOS=linux GOARCH=amd64 go build \
 -ldflags "-s -w" \
 -ldflags "-X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
 -o ./release/httpcat ./cmd/httpcat.go

# Build for Windows
echo "Building for Windows"
GOOS=windows GOARCH=amd64 go build \
 -ldflags "-s -w" \
 -ldflags "-X gin_web_demo/server/common.Version=$HTTPCAT_VERSION -X gin_web_demo/server/common.Build=$HTTPCAT_BUILD -X gin_web_demo/server/common.Commit=$COMMIT_ID" \
 -o ./release/httpcat.exe ./cmd/httpcat.go

# 修改源码 README.md、translations/README-cn.md 文件中的 httpcat_version="v0.1.3" 为当前版本号
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
chmod +x release/install.sh

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

