package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"httpcat/internal/models"

	"github.com/gin-gonic/gin"
)

// GetOperationLogs 获取操作日志列表（分页 + 筛选）
func GetOperationLogs(c *gin.Context) {
	var params struct {
		Current  int    `form:"current" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		Action   string `form:"action"`
		Username string `form:"username"`
		IP       string `form:"ip"`
		Path     string `form:"path"`
		Detail   string `form:"detail"`
		DateFrom string `form:"dateFrom"` // 开始日期 YYYY-MM-DD
		DateTo   string `form:"dateTo"`   // 结束日期 YYYY-MM-DD
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.CreateResponse(c, common.ParamInvalidErrorCode, err.Error())
		return
	}

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Current - 1) * params.PageSize

	// 构建查询条件
	query := db.Model(&models.OperationLogModel{})
	countQuery := db.Model(&models.OperationLogModel{})

	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
		countQuery = countQuery.Where("action = ?", params.Action)
	}
	if params.Username != "" {
		query = query.Where("username LIKE ?", "%"+params.Username+"%")
		countQuery = countQuery.Where("username LIKE ?", "%"+params.Username+"%")
	}
	if params.IP != "" {
		query = query.Where("ip LIKE ?", "%"+params.IP+"%")
		countQuery = countQuery.Where("ip LIKE ?", "%"+params.IP+"%")
	}
	if params.Path != "" {
		query = query.Where("path LIKE ?", "%"+params.Path+"%")
		countQuery = countQuery.Where("path LIKE ?", "%"+params.Path+"%")
	}
	if params.Detail != "" {
		query = query.Where("detail LIKE ?", "%"+params.Detail+"%")
		countQuery = countQuery.Where("detail LIKE ?", "%"+params.Detail+"%")
	}
	if params.DateFrom != "" {
		query = query.Where("created_at >= ?", params.DateFrom+" 00:00:00")
		countQuery = countQuery.Where("created_at >= ?", params.DateFrom+" 00:00:00")
	}
	if params.DateTo != "" {
		query = query.Where("created_at <= ?", params.DateTo+" 23:59:59")
		countQuery = countQuery.Where("created_at <= ?", params.DateTo+" 23:59:59")
	}

	// 查总数
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 分页查询
	var logs []models.OperationLogModel
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"list":     logs,
		"current":  params.Current,
		"pageSize": params.PageSize,
		"total":    total,
	})
}

// GetOperationStats 获取操作日志统计
func GetOperationStats(c *gin.Context) {
	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 总操作数
	var totalCount int64
	db.Model(&models.OperationLogModel{}).Count(&totalCount)

	// 今日操作数
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	db.Model(&models.OperationLogModel{}).Where("created_at >= ?", today+" 00:00:00").Count(&todayCount)

	// 各类型操作统计
	type ActionCount struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	var actionCounts []ActionCount
	db.Model(&models.OperationLogModel{}).
		Select("action, count(*) as count").
		Group("action").
		Order("count DESC").
		Find(&actionCounts)

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"totalCount":   totalCount,
		"todayCount":   todayCount,
		"actionCounts": actionCounts,
	})
}

// ========== 操作日志记录中间件 ==========

// 路径到操作类型的映射
var pathActionMap = map[string]string{
	"/api/v1/file/upload":      "upload",
	"/api/v1/file/download":    "download",
	"/api/v1/file/delete":      "delete",
	"/api/v1/file/rename":      "rename",
	"/api/v1/file/mkdir":       "mkdir",
	"/api/v1/file/preview":     "preview",
	"/api/v1/file/previewInfo": "preview",
	"/api/v1/file/downloadZip": "download_zip",
	"/api/v1/imageManage/upload":  "image_upload",
	"/api/v1/imageManage/delete":  "image_delete",
	"/api/v1/imageManage/rename":  "image_rename",
	"/api/v1/imageManage/clear":   "image_clear",
	"/api/v1/share":               "share_create",
	"/api/v1/user/login/account":  "login",
	"/api/v1/user/changePasswd":   "change_password",
	"/api/v1/conf/sysConfig":      "config_update",
	"/api/v1/conf/restart":        "restart",
}

// OperationLogger Gin 中间件：自动记录操作日志
func OperationLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过 GET 请求中不需要记录的（列表查询类）
		path := c.Request.URL.Path
		method := c.Request.Method

		// 只记录写操作 + 关键读操作（download/preview）
		action := resolveAction(path, method)
		if action == "" {
			c.Next()
			return
		}

		// 特殊处理：DELETE /api/v1/share/:code
		if strings.HasPrefix(path, "/api/v1/share/") && method == "DELETE" {
			action = "share_delete"
		}
		// PUT /api/v1/conf/sysConfig
		if path == "/api/v1/conf/sysConfig" && method == "PUT" {
			action = "config_update"
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()

		// 获取用户名
		username := ""
		if user, exists := c.Get("user"); exists {
			username = fmt.Sprintf("%v", user)
		}

		// 获取请求详情
		detail := buildDetail(c, action)

		// 异步写入日志
		logEntry := models.OperationLogModel{
			Username:  username,
			IP:        c.ClientIP(),
			Method:    method,
			Path:      path,
			Action:    action,
			Detail:    detail,
			Status:    c.Writer.Status(),
			Latency:   latency,
			UserAgent: truncateString(c.Request.UserAgent(), 512),
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		go saveOperationLog(logEntry)
	}
}

func resolveAction(path, method string) string {
	// 精确匹配
	if action, ok := pathActionMap[path]; ok {
		// 区分 GET/POST/PUT/DELETE
		switch {
		case path == "/api/v1/conf/sysConfig" && method == "GET":
			return "" // 查看配置不记录
		case path == "/api/v1/file/download" && method == "GET":
			return action
		case path == "/api/v1/file/preview" && method == "GET":
			return action
		case path == "/api/v1/file/previewInfo" && method == "GET":
			return "" // previewInfo 是辅助接口，不单独记录
		case method == "GET":
			return "" // 大部分 GET 不记录
		default:
			return action
		}
	}

	// 模糊匹配：DELETE /api/v1/share/:code
	if strings.HasPrefix(path, "/api/v1/share/") && method == "DELETE" {
		return "share_delete"
	}

	return ""
}

func buildDetail(c *gin.Context, action string) string {
	switch action {
	case "upload":
		dir := c.PostForm("dir")
		if dir != "" {
			return fmt.Sprintf("上传到目录: %s", dir)
		}
		return "上传文件"
	case "download":
		filename := c.Query("filename")
		return fmt.Sprintf("下载文件: %s", filename)
	case "delete":
		return "删除文件"
	case "rename":
		return "重命名文件"
	case "mkdir":
		return "创建文件夹"
	case "preview":
		filename := c.Query("filename")
		return fmt.Sprintf("预览文件: %s", filename)
	case "download_zip":
		return "打包下载"
	case "share_create":
		return "创建分享"
	case "share_delete":
		code := strings.TrimPrefix(c.Request.URL.Path, "/api/v1/share/")
		return fmt.Sprintf("删除分享: %s", code)
	case "login":
		return "用户登录"
	case "change_password":
		return "修改密码"
	case "config_update":
		return "更新系统配置"
	case "restart":
		return "重启服务"
	case "image_upload":
		return "上传图片"
	case "image_delete":
		filename := c.Query("filename")
		return fmt.Sprintf("删除图片: %s", filename)
	case "image_rename":
		return "重命名图片"
	case "image_clear":
		return "清空所有图片"
	case "chunk_upload_init":
		return "初始化分片上传"
	case "chunk_upload_complete":
		return "完成分片上传"
	case "chunk_upload_abort":
		return "中止分片上传"
	default:
		return ""
	}
}

func saveOperationLog(log models.OperationLogModel) {
	db, err := common.GetDB()
	if err != nil {
		ylog.Errorf("saveOperationLog", "获取数据库连接失败: %v", err)
		return
	}
	if err := db.Create(&log).Error; err != nil {
		ylog.Errorf("saveOperationLog", "记录操作日志失败: %v", err)
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
