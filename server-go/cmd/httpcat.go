package main

import (
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"httpcat/internal/server"
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
	go debug()

	// 同时监听系统信号和 API 触发的重启信号
	select {
	case <-common.Sig:
		ylog.Infof("[MAIN]", "收到系统信号，直接退出")
	case <-common.RestartChan:
		ylog.Infof("[MAIN]", "收到 API 重启信号，执行优雅关闭")
		server.GracefulShutdown()
	}
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
