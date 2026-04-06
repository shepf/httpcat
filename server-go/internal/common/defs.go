package common

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

var (
	Sig = make(chan os.Signal, 1)

	StartTime time.Time // 程序启动时间

	UserConfig *viper.Viper
	ConfPath   string

	HttpPort           int
	RunningHttpPort    int // 实际运行的 HTTP 端口（启动后不变）
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
	FileBaseDir         string // 文件根目录（默认 "./" 即项目工作目录，生产环境建议改为绝对路径，只能通过配置文件修改）
	UploadDir           string // 上传子目录（相对于 FileBaseDir）
	DownloadDir         string // 下载子目录（相对于 FileBaseDir）
	FileUploadEnable    bool //
	EnableUploadToken   bool //是否开启文件上传token校验
	AppKey              string
	AppSecret           string
	PersistentNotifyURL string
	EnableSqlite        bool
	SqliteDBPath        string

	McpEnable    bool   // 是否启用 MCP Server
	McpAuthToken string // MCP 认证 Token（可选，为空则不验证）

	ShareEnable          bool // 是否启用分享功能
	ShareAnonymousAccess bool // 是否允许匿名访问分享链接（false 时需要登录）

	OpenAPIEnable bool // 是否启用 Open API（AK/SK 签名认证）

	PProfEnable bool
	PProfPort   int //pprof

	// 缩略图配置
	ThumbWidth  int // 缩略图宽度
	ThumbHeight int // 缩略图高度

	// 上传策略
	UploadPolicyDeadline   int64 // 上传策略有效期(秒)
	UploadPolicyFSizeMin   int64 // 上传文件最小值(字节)
	UploadPolicyFSizeLimit int64 // 上传文件最大值(字节)

	// 日志级别
	LogLevel int // 日志级别

	// 重启信号通道（handler 写入，main 消费）
	RestartChan = make(chan struct{}, 1)
)

// GetUploadDir 获取完整的上传目录路径（FileBaseDir + UploadDir）
func GetUploadDir() string {
	return filepath.Join(FileBaseDir, UploadDir) + "/"
}

// GetDownloadDir 获取完整的下载目录路径（FileBaseDir + DownloadDir）
func GetDownloadDir() string {
	return filepath.Join(FileBaseDir, DownloadDir) + "/"
}
