package common

import (
	"fmt"
	"gin_web_demo/server/common/userconfig"
	"gin_web_demo/server/common/ylog"
	"github.com/spf13/pflag"
	"os"
)

func init() {
	// 结尾的Var表示支持将参数的值，绑定到变量
	pflag.StringVar(&StaticDir, "static", "./website/static/", "指定静态资源路径(可用于web界面)")
	pflag.StringVar(&UploadDir, "upload", "./website/upload/", "指定上传文件的路径，右斜线结尾")
	pflag.StringVar(&DownloadDir, "download", "./website/download/", "指定下载文件的路径，右斜线结尾")

	// 结尾的P表示支持短选项
	pflag.IntVarP(&HttpPort, "port", "P", 8888, "host port.")
	// 结尾的P表示支持短选项
	confPath := pflag.StringP("config", "C", "./conf/svr.yml", "ConfigPath")

	pflag.Parse()
	ConfPath = *confPath

	initConfig()
}

func initConfig() {
	var err error
	if UserConfig, err = userconfig.NewUserConfig(userconfig.WithPath(ConfPath)); err != nil {
		fmt.Printf("####LOAD_CONFIG_ERROR: %v", err)
		os.Exit(-1)
	}
	initLog()
	initDefault()

}

func initDefault() {
	// 打印初始化
	fmt.Println("####初始化:", "initDefault")

	SSLKeyFile = UserConfig.GetString("server.ssl.keyfile")
	SSLCertFile = UserConfig.GetString("server.ssl.certfile")
	SSLRawDataKeyFile = UserConfig.GetString("server.ssl.rawdata_keyfile")
	SSLRawDataCertFile = UserConfig.GetString("server.ssl.rawdata_certfile")
	SSLCaFile = UserConfig.GetString("server.ssl.cafile")
	if SSLKeyFile == "" || SSLCertFile == "" || SSLCaFile == "" || SSLRawDataKeyFile == "" || SSLRawDataCertFile == "" {
		ylog.Fatalf("init", "ssl file empty SSLKeyFile:%s SSLCertFile:%s SSLCaFile:%s SSLRawDataKeyFile:%s SSLRawDataCertFile:%s", SSLKeyFile, SSLCertFile, SSLCaFile, SSLRawDataKeyFile, SSLRawDataCertFile)
	}

	HttpPort = UserConfig.GetInt("server.http.port")
	HttpSSLEnable = UserConfig.GetBool("server.http.ssl.enable")
	HttpAuthEnable = UserConfig.GetBool("server.http.auth.enable")
	HttpAkSkMap = UserConfig.GetStringMapString("server.http.auth.aksk")

	// for 业务
	FileEnable = UserConfig.GetBool("server.http.file.enable")

}

func initLog() {
	logLevel := UserConfig.GetInt("server.log.applog.loglevel")
	logPath := UserConfig.GetString("server.log.applog.path")
	logger := ylog.NewYLog(
		ylog.WithLogFile(logPath),
		ylog.WithMaxAge(3),
		ylog.WithMaxSize(10),
		ylog.WithMaxBackups(3),
		ylog.WithLevel(logLevel),
	)
	ylog.InitLogger(logger)
}
