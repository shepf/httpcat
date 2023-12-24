package main

import (
	"fmt"
	"gin_web_demo/server"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	"net/http"
	"os/signal"
	"syscall"
)

func init() {
	signal.Notify(common.Sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

func main() {
	ylog.Infof("[MAIN]", "START_SERVER")

	fmt.Printf("Version: %s\n", common.Version)
	fmt.Printf("Commit: %s\n", common.Commit)
	fmt.Printf("Build: %s\n", common.Build)
	fmt.Printf("CI: %s\n", common.CI)

	go server.RunAPIServer(common.HttpPort, common.HttpSSLEnable, common.HttpAuthEnable, common.SSLCertFile, common.SSLKeyFile)

	<-common.Sig
}

func debug() {
	//start pprof for debug
	if common.PProfEnable {
		err := http.ListenAndServe(fmt.Sprintf(":%d", common.PProfPort), nil)
		if err != nil {
			ylog.Errorf("[MAIN]", "pprof ListenAndServe Error %s", err.Error())
		}
	}
}
