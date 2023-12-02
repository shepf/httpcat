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
	pflag.StringVar(&StaticDir, "static", "./website/static/", "指定静态资源路径(web)")
	pflag.StringVar(&UploadDir, "upload", "./website/upload/", "指定上传文件的路径,右斜线结尾")
	pflag.StringVar(&DownloadDir, "download", "./website/download/", "指定下载文件的路径,右斜线结尾")

	// 结尾的P表示支持短选项
	pflag.IntVarP(&HttpPort, "port", "P", 0, "host port.")
	pflag.IntVar(&P2pListenPort, "p2pport", 0, "p2p host port.")
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

	// 优先使用命令行传入值
	if HttpPort == 0 {
		HttpPort = UserConfig.GetInt("server.http.port")
	}
	HttpReadTimeout = UserConfig.GetInt64("server.http.read_timeout")
	HttpWriteTimeout = UserConfig.GetInt64("server.http.write_timeout")
	HttpIdleTimeout = UserConfig.GetInt64("server.http.idle_timeout")

	HttpSSLEnable = UserConfig.GetBool("server.http.ssl.enable")
	HttpAuthEnable = UserConfig.GetBool("server.http.auth.enable")
	HttpAkSkMap = UserConfig.GetStringMapString("server.http.auth.aksk")

	// for 业务
	FileEnable = UserConfig.GetBool("server.http.file.enable")
	EnableUploadToken = UserConfig.GetBool("server.http.file.enable_upload_token")
	AppKey = UserConfig.GetString("server.http.file.app_key")
	AppSecret = UserConfig.GetString("server.http.file.app_secret")
	PersistentNotifyURL = UserConfig.GetString("server.http.file.upload_policy.persistent_notify_url")
	EnableSqlite = UserConfig.GetBool("server.http.file.enable_sqlite")
	SqliteDBPath = UserConfig.GetString("server.http.file.sqlite_db_path")

	// for p2p
	P2pEnable = UserConfig.GetBool("server.p2p.enable")
	P2pListenIP = UserConfig.GetString("server.p2p.listen.ip")
	// 优先使用命令行传入值
	if P2pListenPort == 0 {
		P2pListenPort = UserConfig.GetInt("server.p2p.listen.port")
	}
	// discovery node
	EnableMdns = UserConfig.GetBool("server.p2p.mdns.enable")
	RendezvousString = UserConfig.GetString("server.p2p.mdns.rendezvous")
	// pubsub
	EnablePubSub = UserConfig.GetBool("server.p2p.pubsub.enable")
	TopicName = UserConfig.GetString("server.p2p.pubsub.topic_name")

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
