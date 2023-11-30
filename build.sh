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

HTTPCAT_VERSION=v0.0.4
HTTPCAT_BUILD=$(date "+%Y%m%d%H%M")
GOOS=linux GOARCH=amd64 go build \
 -ldflags "-s -w" \
 -ldflags "-X gin_web_demo/server.Version=$HTTPCAT_VERSION -X gin_web_demo/server.Build=$HTTPCAT_BUILD" \
 -o ./release/httpcat ./cmd/httpcat.go

# Build for Windows
echo "Building for Windows"
GOOS=windows GOARCH=amd64 go build \
 -ldflags "-s -w" \
 -ldflags "-X gin_web_demo.server.Version=$HTTPCAT_VERSION -X gin_web_demo.server.Build=$HTTPCAT_BUILD" \
 -o ./release/httpcat.exe ./cmd/httpcat.go


# Package configuration file and static files
cp -r server/conf release/
cp -r static release/
cp -r httpcat.service release/

# Create release archive for Linux
cd release
tar zcvf httpcat_$HTTPCAT_VERSION.tar.gz ./*

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

