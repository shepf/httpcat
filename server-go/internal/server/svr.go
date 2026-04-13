package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"

	"github.com/gin-gonic/gin"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var ps *pubsub.PubSub

// 使用一个全局变量 subscribedTopics，它是一个映射（map），用于保存已加入的主题及其对应的 *pubsub.Topic 实例
var subscribedTopics = make(map[string]*pubsub.Topic)

// httpSrv 全局 HTTP Server 实例，用于优雅关闭
var httpSrv *http.Server

// GracefulShutdown 优雅关闭 HTTP 服务器，等待现有请求处理完毕（最多 10 秒）
// 关闭后进程将退出，由 systemd/Docker 自动重启
func GracefulShutdown() {
	if httpSrv == nil {
		ylog.Errorf("GracefulShutdown", "httpSrv is nil, 直接退出进程")
		os.Exit(1)
		return
	}

	ylog.Infof("GracefulShutdown", "开始优雅关闭服务器...")

	// 给 10 秒时间等待现有请求处理完毕
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		ylog.Errorf("GracefulShutdown", "优雅关闭失败: %v，强制退出", err)
	} else {
		ylog.Infof("GracefulShutdown", "服务器已优雅关闭")
	}

	// 退出进程，由 systemd(Restart=always) / Docker(restart:unless-stopped) 自动拉起
	os.Exit(0)
}

func RunAPIServer(port int, enableSSL, enableAuth bool, certFile, keyFile string) {
	//生成一个 Engine，这是 gin 的核心，默认带有 Logger 和 Recovery 两个中间件
	router := gin.Default()
	RegisterRouter(router)

	// 创建http server
	httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
		//ReadTimeout：服务器在读取客户端请求时，等待的最大时间。
		ReadTimeout: time.Duration(common.HttpReadTimeout) * time.Second,
		// 服务器在写回应应答时，等待的最大时间。
		WriteTimeout: time.Duration(common.HttpWriteTimeout) * time.Second,
		// 一个连接在空闲状态下可以存在的最长时间。
		IdleTimeout: time.Duration(common.HttpIdleTimeout) * time.Second,
	}
	ctx := context.Background()

	enableP2P := common.P2pEnable
	if enableP2P {
		go runP2PServer(ctx, router)
	}

	var err error
	ylog.Infof("RunServer", "####HTTP_LISTEN_ON:%d", port)
	if enableSSL {
		err = httpSrv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = httpSrv.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		ylog.Errorf("RunServer", "####http run error: %v", err)
	}
}
