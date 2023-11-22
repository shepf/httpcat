package common

import (
	"github.com/spf13/viper"
	"os"
)

var (
	Sig = make(chan os.Signal, 1)

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
	HttpAuthEnable     bool
	HttpAkSkMap        map[string]string //access key and secret key list, which used to identify whether the http request comes from a known subject
	SvrAK              string            // access key, which use for http sign
	SvrSK              string            // secret key, which use for http sign
	P2pEnable          bool
	P2pListenIP        string
	P2pListenPort      int
	EnableMdns         bool
	RendezvousString   string

	StaticDir   string
	UploadDir   string
	DownloadDir string
	FileEnable  bool // 决定是否注册file路由，false：就只能做http服务使用 true：文件上传下载等功能

	PProfEnable bool
	PProfPort   int //pprof
)
