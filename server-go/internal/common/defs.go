package common

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

var (
	Sig = make(chan os.Signal, 1)

	StartTime time.Time // 程序启动时间

	UserConfig *viper.Viper
	ConfPath   string

	HttpPort           int
	HttpReadTimeout    int64
	HttpWriteTimeout   int64
	HttpIdleTimeout    int64
	HttpSSLEnable      bool
	SSLKeyFile         string
	SSLCertFile        string
	SSLRawDataKeyFile  string
	SSLRawDataCertFile string
	SSLCaFile          string
	JwtSecret          string
	HttpAuthEnable     bool
	HttpAkSkMap        map[string]string //access key and secret key list, which used to identify whether the http request comes from a known subject
	SvrAK              string            // access key, which use for http sign
	SvrSK              string            // secret key, which use for http sign
	P2pEnable          bool
	P2pListenIP        string
	P2pListenPort      int
	EnableMdns         bool
	RendezvousString   string
	EnablePubSub       bool
	TopicName          string

	StaticDir           string
	UploadDir           string
	DownloadDir         string
	FileUploadEnable    bool //
	EnableUploadToken   bool //是否开启文件上传token校验
	AppKey              string
	AppSecret           string
	PersistentNotifyURL string
	EnableSqlite        bool
	SqliteDBPath        string

	McpEnable    bool   // 是否启用 MCP Server
	McpAuthToken string // MCP 认证 Token（可选，为空则不验证）

	PProfEnable bool
	PProfPort   int //pprof
)
