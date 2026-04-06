package v1

import (
	"fmt"
	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"httpcat/internal/midware"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SysConfigResponse 系统配置响应结构
type SysConfigResponse struct {
	// 文件根目录（只读，只能通过配置文件修改）
	FileBaseDir string `json:"fileBaseDir"`

	// 存储子目录（相对于 FileBaseDir）
	UploadDir   string `json:"uploadDir"`
	DownloadDir string `json:"downloadDir"`

	// 完整路径（只读，展示用）
	FullUploadDir   string `json:"fullUploadDir"`
	FullDownloadDir string `json:"fullDownloadDir"`

	// HTTP 服务
	HttpPort int `json:"httpPort"`

	// 文件上传
	FileUploadEnable       bool   `json:"fileUploadEnable"`
	EnableUploadToken      bool   `json:"enableUploadToken"`
	UploadPolicyDeadline   int64  `json:"uploadPolicyDeadline"`
	UploadPolicyFSizeMin   int64  `json:"uploadPolicyFSizeMin"`
	UploadPolicyFSizeLimit int64  `json:"uploadPolicyFSizeLimit"`

	// 企业微信 Bot
	PersistentNotifyURL string `json:"persistentNotifyUrl"`
	NotifyEnable        bool   `json:"notifyEnable"`

	// 缩略图
	ThumbWidth  int `json:"thumbWidth"`
	ThumbHeight int `json:"thumbHeight"`

	// 日志
	LogLevel int `json:"logLevel"`
}

// SysConfigUpdateRequest 系统配置更新请求结构
type SysConfigUpdateRequest struct {
	// 存储子目录（相对于 FileBaseDir，如 upload/、download/）
	UploadDir   *string `json:"uploadDir"`
	DownloadDir *string `json:"downloadDir"`

	// HTTP 服务
	HttpPort *int `json:"httpPort"`

	// 文件上传
	FileUploadEnable       *bool  `json:"fileUploadEnable"`
	EnableUploadToken      *bool  `json:"enableUploadToken"`
	UploadPolicyDeadline   *int64 `json:"uploadPolicyDeadline"`
	UploadPolicyFSizeMin   *int64 `json:"uploadPolicyFSizeMin"`
	UploadPolicyFSizeLimit *int64 `json:"uploadPolicyFSizeLimit"`

	// 企业微信 Bot
	PersistentNotifyURL *string `json:"persistentNotifyUrl"`
	NotifyEnable        *bool   `json:"notifyEnable"`

	// 缩略图
	ThumbWidth  *int `json:"thumbWidth"`
	ThumbHeight *int `json:"thumbHeight"`

	// 日志
	LogLevel *int `json:"logLevel"`
}

// GetSysConfig 获取系统配置
func GetSysConfig(c *gin.Context) {
	// 文件根目录：如果是默认值 "./"，展示绝对路径
	fileBaseDir := common.FileBaseDir
	absFileBaseDir, _ := filepath.Abs(fileBaseDir)
	if fileBaseDir == "" || fileBaseDir == "./" || fileBaseDir == "." {
		fileBaseDir = absFileBaseDir
	}

	// 完整路径也用绝对路径
	absUploadDir, _ := filepath.Abs(common.GetUploadDir())
	absDownloadDir, _ := filepath.Abs(common.GetDownloadDir())

	config := SysConfigResponse{
		FileBaseDir:            fileBaseDir,
		UploadDir:              common.UploadDir,
		DownloadDir:            common.DownloadDir,
		FullUploadDir:          absUploadDir + "/",
		FullDownloadDir:        absDownloadDir + "/",
		HttpPort:               common.RunningHttpPort,
		FileUploadEnable:       common.FileUploadEnable,
		EnableUploadToken:      common.EnableUploadToken,
		UploadPolicyDeadline:   common.UploadPolicyDeadline,
		UploadPolicyFSizeMin:   common.UploadPolicyFSizeMin,
		UploadPolicyFSizeLimit: common.UploadPolicyFSizeLimit,
		PersistentNotifyURL:    common.PersistentNotifyURL,
		NotifyEnable:           common.PersistentNotifyURL != "",
		ThumbWidth:             common.ThumbWidth,
		ThumbHeight:            common.ThumbHeight,
		LogLevel:               common.LogLevel,
	}

	common.CreateResponse(c, common.SuccessCode, config)
}

// validateSubDirPath 验证子目录路径安全性
// 子目录路径基于 FileBaseDir，如 upload/、download/、images/upload/
func validateSubDirPath(path string) (string, error) {
	return common.NormalizeSafeSubDirPath(path)
}

// UpdateSysConfig 更新系统配置
func UpdateSysConfig(c *gin.Context) {
	var req SysConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "请求参数格式错误: "+err.Error())
		return
	}

	// 需要重启才生效的配置变更
	needRestart := false
	// 记录变更项
	changes := make([]string, 0)

	// === 存储子目录（需重启） ===
	if req.UploadDir != nil {
		validPath, err := validateSubDirPath(*req.UploadDir)
		if err != nil {
			common.BadRequest(c, "上传子目录无效: "+err.Error())
			return
		}
		if validPath != common.UploadDir {
			common.UploadDir = validPath
			common.UserConfig.Set("server.http.file.upload_dir", validPath)
			needRestart = true
			changes = append(changes, "uploadDir")
			// 确保完整目录存在
			fullPath := common.GetUploadDir()
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				os.MkdirAll(fullPath, 0755)
			}
		}
	}

	if req.DownloadDir != nil {
		validPath, err := validateSubDirPath(*req.DownloadDir)
		if err != nil {
			common.BadRequest(c, "下载子目录无效: "+err.Error())
			return
		}
		if validPath != common.DownloadDir {
			common.DownloadDir = validPath
			common.UserConfig.Set("server.http.file.download_dir", validPath)
			needRestart = true
			changes = append(changes, "downloadDir")
			// 确保完整目录存在
			fullPath := common.GetDownloadDir()
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				os.MkdirAll(fullPath, 0755)
			}
		}
	}

	// === HTTP 端口（需重启） ===
	if req.HttpPort != nil {
		port := *req.HttpPort
		if port < 1 || port > 65535 {
			common.BadRequest(c, "HTTP 端口必须在 1-65535 范围内")
			return
		}
		if port != common.RunningHttpPort {
			// 只写配置文件，不修改内存值，重启后才生效
			common.UserConfig.Set("server.http.port", port)
			needRestart = true
			changes = append(changes, "httpPort")
		}
	}

	// === 文件上传（热更新） ===
	if req.FileUploadEnable != nil {
		common.FileUploadEnable = *req.FileUploadEnable
		common.UserConfig.Set("server.http.file.upload_enable", *req.FileUploadEnable)
		changes = append(changes, "fileUploadEnable")
	}

	if req.EnableUploadToken != nil {
		common.EnableUploadToken = *req.EnableUploadToken
		common.UserConfig.Set("server.http.file.enable_upload_token", *req.EnableUploadToken)
		changes = append(changes, "enableUploadToken")
	}

	if req.UploadPolicyDeadline != nil {
		if *req.UploadPolicyDeadline < 0 {
			common.BadRequest(c, "上传策略有效期不能为负数")
			return
		}
		common.UploadPolicyDeadline = *req.UploadPolicyDeadline
		common.UserConfig.Set("server.http.file.upload_policy.deadline", *req.UploadPolicyDeadline)
		changes = append(changes, "uploadPolicyDeadline")
	}

	if req.UploadPolicyFSizeMin != nil {
		if *req.UploadPolicyFSizeMin < 0 {
			common.BadRequest(c, "文件最小值不能为负数")
			return
		}
		common.UploadPolicyFSizeMin = *req.UploadPolicyFSizeMin
		common.UserConfig.Set("server.http.file.upload_policy.fsizemin", *req.UploadPolicyFSizeMin)
		changes = append(changes, "uploadPolicyFSizeMin")
	}

	if req.UploadPolicyFSizeLimit != nil {
		if *req.UploadPolicyFSizeLimit < 0 {
			common.BadRequest(c, "文件最大值不能为负数")
			return
		}
		common.UploadPolicyFSizeLimit = *req.UploadPolicyFSizeLimit
		common.UserConfig.Set("server.http.file.upload_policy.fsizeLimit", *req.UploadPolicyFSizeLimit)
		changes = append(changes, "uploadPolicyFSizeLimit")
	}

	// === 企业微信 Bot（热更新） ===
	if req.NotifyEnable != nil && !*req.NotifyEnable {
		// 关闭通知时清空 URL
		common.PersistentNotifyURL = ""
		common.UserConfig.Set("server.http.file.upload_policy.persistent_notify_url", "")
		changes = append(changes, "notifyEnable")
	} else if req.PersistentNotifyURL != nil {
		url := strings.TrimSpace(*req.PersistentNotifyURL)
		// 基本 URL 校验
		if url != "" && !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
			common.BadRequest(c, "Webhook URL 必须以 http:// 或 https:// 开头")
			return
		}
		common.PersistentNotifyURL = url
		common.UserConfig.Set("server.http.file.upload_policy.persistent_notify_url", url)
		changes = append(changes, "persistentNotifyUrl")
	}

	// === 缩略图尺寸（热更新） ===
	if req.ThumbWidth != nil {
		if *req.ThumbWidth < 50 || *req.ThumbWidth > 2000 {
			common.BadRequest(c, "缩略图宽度必须在 50-2000 范围内")
			return
		}
		common.ThumbWidth = *req.ThumbWidth
		common.UserConfig.Set("server.http.file.thumb_width", *req.ThumbWidth)
		changes = append(changes, "thumbWidth")
	}

	if req.ThumbHeight != nil {
		if *req.ThumbHeight < 50 || *req.ThumbHeight > 2000 {
			common.BadRequest(c, "缩略图高度必须在 50-2000 范围内")
			return
		}
		common.ThumbHeight = *req.ThumbHeight
		common.UserConfig.Set("server.http.file.thumb_height", *req.ThumbHeight)
		changes = append(changes, "thumbHeight")
	}

	// === 日志级别（热更新） ===
	if req.LogLevel != nil {
		if *req.LogLevel < -1 || *req.LogLevel > 5 {
			common.BadRequest(c, "日志级别必须在 -1(Debug) 到 5(Fatal) 范围内")
			return
		}
		common.LogLevel = *req.LogLevel
		common.UserConfig.Set("server.log.applog.loglevel", *req.LogLevel)
		// 动态更新日志级别
		ylog.SetLevel(*req.LogLevel)
		changes = append(changes, "logLevel")
	}

	// 持久化到配置文件
	if len(changes) > 0 {
		if err := saveConfig(); err != nil {
			ylog.Errorf("UpdateSysConfig", "保存配置文件失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorCode": common.ErrorCode,
				"msg":       "保存配置文件失败: " + err.Error(),
			})
			return
		}
		ylog.Infof("UpdateSysConfig", "配置已更新: %v", changes)
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"changes":     changes,
		"needRestart": needRestart,
		"message":     formatResultMessage(changes, needRestart),
	})
}

// saveConfig 将当前 viper 配置保存到 YAML 文件
func saveConfig() error {
	confPath := common.ConfPath

	// 创建新的 viper 实例读取原始文件
	v := viper.New()
	v.SetConfigFile(confPath)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 将内存中的变更合并到文件配置
	for _, key := range common.UserConfig.AllKeys() {
		v.Set(key, common.UserConfig.Get(key))
	}

	// 写回文件
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

func formatResultMessage(changes []string, needRestart bool) string {
	if len(changes) == 0 {
		return "没有配置变更"
	}
	msg := fmt.Sprintf("已更新 %d 项配置", len(changes))
	if needRestart {
		msg += "，部分配置(存储路径/HTTP端口)需要重启服务后生效"
	}
	return msg
}

// RestartRequest 重启请求结构
type RestartRequest struct {
	Password string `json:"password" binding:"required"`
}

// RestartServer 重启服务（需验证管理员密码）
// 流程：验证密码 → 返回成功响应 → 延迟 1.5 秒 → 优雅关闭 → systemd/Docker 自动拉起
func RestartServer(c *gin.Context) {
	var req RestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "请输入管理员密码")
		return
	}

	// 验证管理员密码
	_, err := midware.CheckUser("admin", req.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"errorCode": common.AuthFailedErrorCode,
			"msg":       "管理员密码错误",
		})
		return
	}

	ylog.Infof("RestartServer", "管理员已确认重启服务，1.5 秒后执行优雅关闭")

	// 先返回响应给前端
	common.CreateResponse(c, common.SuccessCode, gin.H{
		"message": "服务即将重启，请稍候...",
	})

	// 延迟发送重启信号，确保 HTTP 响应已经发送给客户端
	go func() {
		time.Sleep(1500 * time.Millisecond)
		select {
		case common.RestartChan <- struct{}{}:
			ylog.Infof("RestartServer", "重启信号已发送")
		default:
			ylog.Errorf("RestartServer", "重启信号通道已满，可能已有重启请求在处理")
		}
	}()
}
