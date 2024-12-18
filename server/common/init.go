package common

import (
	"fmt"
	"gin_web_demo/server/common/userconfig"
	"gin_web_demo/server/common/ylog"
	"github.com/spf13/pflag"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

var (
	Version string
	Commit  string
	Build   string
	CI      string
)

func init() {
	// 结尾的Var表示支持将参数的值，绑定到变量
	pflag.StringVar(&StaticDir, "static", "./website/static/", "Specify the path for static resources (web), ending with a forward slash (/)")
	pflag.StringVar(&UploadDir, "upload", "./website/upload/", "Specify the path for uploading files, ending with a forward slash (/)")
	pflag.StringVar(&DownloadDir, "download", "./website/download/", "Specify the path for downloading files, ending with a forward slash (/)")

	// 结尾的P表示支持短选项
	pflag.IntVarP(&HttpPort, "port", "P", 0, "host port.")
	pflag.IntVar(&P2pListenPort, "p2pport", 0, "p2p host port.")
	// 结尾的P表示支持短选项
	confPath := pflag.StringP("config", "C", "./conf/svr.yml", "ConfigPath")
	showVersion := pflag.BoolP("version", "v", false, "Show the version number")

	pflag.Parse()
	ConfPath = *confPath

	if *showVersion {
		fmt.Println("Version: ", Version) // 替换为实际的版本号
		fmt.Println("Build time: ", Build)
		fmt.Println("Commit id: ", Commit)

		os.Exit(0)
	}

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
	//数据库相关初始化
	initDB()

}

func initDefault() {
	// 打印初始化
	fmt.Println("####初始化:", "initDefault")
	// 在应用程序启动时记录起始时间
	StartTime = time.Now()

	SSLKeyFile = UserConfig.GetString("server.ssl.keyfile")
	SSLCertFile = UserConfig.GetString("server.ssl.certfile")
	SSLRawDataKeyFile = UserConfig.GetString("server.ssl.rawdata_keyfile")
	SSLRawDataCertFile = UserConfig.GetString("server.ssl.rawdata_certfile")
	SSLCaFile = UserConfig.GetString("server.ssl.cafile")
	if SSLKeyFile == "" || SSLCertFile == "" || SSLCaFile == "" || SSLRawDataKeyFile == "" || SSLRawDataCertFile == "" {
		ylog.Fatalf("init", "ssl file empty SSLKeyFile:%s SSLCertFile:%s SSLCaFile:%s SSLRawDataKeyFile:%s SSLRawDataCertFile:%s", SSLKeyFile, SSLCertFile, SSLCaFile, SSLRawDataKeyFile, SSLRawDataCertFile)
	}

	JwtSecret = UserConfig.GetString("server.http.jwt_secret")

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
	FileUploadEnable = UserConfig.GetBool("server.http.file.upload_enable")
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

	// 程序启动时候，判断上传下载目录是否存在，如果不存在则创建
	// 判断目录是否存在，如果不存在则创建
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(UploadDir, 0755)
		if err != nil {
			ylog.Errorf("initDefault", "创建UploadDir目录失败", err)
			panic(err)
		}
	}
	if _, err := os.Stat(DownloadDir); os.IsNotExist(err) {
		err := os.MkdirAll(DownloadDir, 0755)
		if err != nil {
			ylog.Errorf("initDefault", "创建DownloadDir目录失败", err)
			panic(err)
		}
	}

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

func initDB() {
	// 打印初始化
	fmt.Println("####初始化:", "initDB")
	if EnableSqlite {
		ylog.Infof("initDB", "init start~")
		// 读取 SQLite 数据库文件路径配置项
		dbPath := SqliteDBPath

		dir := filepath.Dir(dbPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				ylog.Errorf("initDB", "failed to create directory: %v", err)
				return
			}
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			file, err := os.Create(dbPath)
			if err != nil {
				ylog.Errorf("initDB", "failed to create database file: %v", err)
				return
			}
			file.Close()
		}

		// 打开 SQLite 数据库连接
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			ylog.Errorf("initDB", "open db failed, err:%v", err)
			return
		}

		type Notification struct {
			ID         uint   `gorm:"column:id;primary_key"`
			IP         string `gorm:"column:ip"`
			Appkey     string `gorm:"column:appkey"`
			UploadTime string `gorm:"column:upload_time"`
			Filename   string `gorm:"column:filename"`
			FileSize   string `gorm:"column:file_size"`
			FileMD5    string `gorm:"column:file_md5"`
		}
		db.AutoMigrate(&Notification{})
		if err != nil {
			ylog.Errorf("initDB", "create notifications table failed, err:%v", err)
			return
		}
		// 创建 users 表（如果不存在）
		InitializeUserTable(db)
		// 加载用户到缓存
		InitUser()

		// 创建 upload_tokens 表（如果不存在）
		InitializeUploadTokenTable(db)
		InitUploadToken()

		// 创建 download_logs 表（如果不存在）
		InitializeDownloadLogTable(db)

		// 创建 upload_logs 表（如果不存在）
		InitializeUploadLogTable(db)

		// 创建 t_upload_image 表（如果不存在）
		InitializeUploadImageTable(db)

		ylog.Infof("initDB", "init end~")
	}

}
