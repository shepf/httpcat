#!/bin/bash

if ! type go >/dev/null 2>&1; then
    echo 'go not installed';
    exit 1
fi

rm -rf output
mkdir -p output/
go mod download
go build -ldflags "-s -w" -o ./output/httpcat ./cmd/httpcat.go

#cd output
#tar zcvf bin.tar.gz ./*
#rm -f httpcat


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

