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

	StaticDir   string
	Port        int
	UploadDir   string
	DownloadDir string

	PProfEnable bool
	PProfPort   int //pprof
)
